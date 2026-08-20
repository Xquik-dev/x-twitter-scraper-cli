// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"net/http/httputil"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/Xquik-dev/x-twitter-scraper-cli/internal/jsonview"
	"github.com/Xquik-dev/x-twitter-scraper-go/option"

	"github.com/charmbracelet/x/term"
	"github.com/itchyny/json2yaml"
	"github.com/muesli/reflow/wrap"
	"github.com/tidwall/gjson"
	"github.com/tidwall/pretty"
	"github.com/urfave/cli/v3"
)

var OutputFormats = []string{"auto", "explore", "json", "jsonl", "pretty", "raw", "yaml"}

// ValidateBaseURL requires an HTTP or HTTPS scheme.
func ValidateBaseURL(value, source string) error {
	if value != "" && !strings.HasPrefix(value, "http://") && !strings.HasPrefix(value, "https://") {
		return fmt.Errorf("%s %q is missing a scheme (expected http:// or https://)", source, value)
	}
	return nil
}

func getDefaultRequestOptions(cmd *cli.Command) []option.RequestOption {
	opts := []option.RequestOption{
		option.WithHeader("User-Agent", fmt.Sprintf("XTwitterScraper/CLI %s", Version)),
		option.WithHeader("X-Stainless-Lang", "cli"),
		option.WithHeader("X-Stainless-Package-Version", Version),
		option.WithHeader("X-Stainless-Runtime", "cli"),
		option.WithHeader("X-Stainless-CLI-Command", cmd.FullName()),
	}
	if cmd.IsSet("api-key") {
		opts = append(opts, option.WithAPIKey(cmd.String("api-key")))
	}
	if cmd.IsSet("bearer-token") {
		opts = append(opts, option.WithBearerToken(cmd.String("bearer-token")))
	}

	if baseURL := cmd.String("base-url"); baseURL != "" {
		opts = append(opts, option.WithBaseURL(baseURL))
	}

	return opts
}

var debugMiddlewareOption = option.WithMiddleware(
	func(r *http.Request, mn option.MiddlewareNext) (*http.Response, error) {
		logger := log.Default()

		if reqBytes, err := httputil.DumpRequest(r, true); err == nil {
			logger.Printf("HTTP Request:\n%s\n", reqBytes)
		}

		resp, err := mn(r)
		if err != nil {
			return resp, err
		}

		if respBytes, err := httputil.DumpResponse(resp, true); err == nil {
			logger.Printf("HTTP Response:\n%s\n", respBytes)
		}

		return resp, err
	},
)

// isInputPiped reports whether stdin has readable data.
func isInputPiped() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	mode := stat.Mode()

	// Empty regular files include Windows NUL and contain no input.
	if mode.IsRegular() && stat.Size() > 0 {
		return true
	}

	// Some terminals connect an empty pipe, so poll before reading.
	if mode&(os.ModeNamedPipe|os.ModeSocket) != 0 {
		return isPipedDataAvailableOSSpecific()
	}

	return false
}

func isTerminal(w io.Writer) bool {
	switch v := w.(type) {
	case *os.File:
		return term.IsTerminal(v.Fd())
	default:
		return false
	}
}

func streamOutput(label string, generateOutput func(w *os.File) error) error {
	if !isTerminal(os.Stdout) {
		return streamToStdout(generateOutput)
	}

	// Unix sockets reduce pagination; Windows and failures use pipes.
	return streamOutputOSSpecific(label, generateOutput)
}

func streamToPagerWithPipe(label string, generateOutput func(w *os.File) error) error {
	r, w, err := os.Pipe()
	if err != nil {
		return err
	}
	defer r.Close()
	defer w.Close()

	pagerProgram := os.Getenv("PAGER")
	if pagerProgram == "" {
		pagerProgram = "less"
	}

	if _, err := exec.LookPath(pagerProgram); err != nil {
		return err
	}

	cmd := exec.Command(pagerProgram)
	cmd.Stdin = r
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(),
		"LESS=-X -r -P "+label,
		"MORE=-r -P "+label,
	)

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := r.Close(); err != nil {
		return err
	}

	// Preserve terminal colors inside the pager.
	if isTerminal(os.Stdout) && os.Getenv("FORCE_COLOR") == "" {
		os.Setenv("FORCE_COLOR", "1")
	}

	if err := generateOutput(w); err != nil && !strings.Contains(err.Error(), "broken pipe") {
		return err
	}

	w.Close()
	return cmd.Wait()
}

func streamToStdout(generateOutput func(w *os.File) error) error {
	signal.Ignore(syscall.SIGPIPE)
	err := generateOutput(os.Stdout)
	if err != nil && strings.Contains(err.Error(), "broken pipe") {
		return nil
	}
	return err
}

// writeBinaryResponse writes bytes to stdout or a file.
func writeBinaryResponse(response *http.Response, stdout io.Writer, outfile string) (string, error) {
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	switch outfile {
	case "-", "/dev/stdout":
		_, err := stdout.Write(body)
		return "", err
	case "":
		// Print text and piped output directly.
		if !isTerminal(os.Stdout) || isUTF8TextFile(body) {
			_, err := stdout.Write(body)
			return "", err
		}

		// Otherwise, use the response filename or a unique fallback.
		file, err := createDownloadFile(response, body)
		if err != nil {
			return "", err
		}
		defer file.Close()
		if _, err := file.Write(body); err != nil {
			return "", err
		}
		return fmt.Sprintf("Wrote output to %s.", file.Name()), nil
	default:
		if err := os.WriteFile(outfile, body, 0644); err != nil {
			return "", err
		}
		return fmt.Sprintf("Wrote output to %s.", outfile), nil
	}
}

// createDownloadFile uses response metadata or detected MIME type.
func createDownloadFile(response *http.Response, data []byte) (*os.File, error) {
	filename := "file"
	disp := response.Header.Get("Content-Disposition")
	_, params, err := mime.ParseMediaType(disp)
	if err == nil {
		if dispFilename, ok := params["filename"]; ok {
			// Strip directories to prevent traversal.
			filename = filepath.Base(dispFilename)
			file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
			if err == nil {
				return file, nil
			}
		}
	}

	ext := filepath.Ext(filename)
	if ext == "" {
		ext = guessExtension(data)
	}
	base := strings.TrimSuffix(filename, ext)
	return os.CreateTemp(".", base+"-*"+ext)
}

func guessExtension(data []byte) string {
	ct := http.DetectContentType(data)

	// Prefer common extensions.
	switch ct {
	case "application/gzip":
		return ".gz"
	case "application/pdf":
		return ".pdf"
	case "application/zip":
		return ".zip"
	case "audio/mpeg":
		return ".mp3"
	case "image/bmp":
		return ".bmp"
	case "image/gif":
		return ".gif"
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	case "video/mp4":
		return ".mp4"
	}

	exts, err := mime.ExtensionsByType(ct)
	if err == nil && len(exts) > 0 {
		return exts[0]
	} else if isUTF8TextFile(data) {
		return ".txt"
	} else {
		return ".bin"
	}
}

func shouldUseColors(w io.Writer) bool {
	force, ok := os.LookupEnv("FORCE_COLOR")
	if ok {
		if force == "1" {
			return true
		}
		if force == "0" {
			return false
		}
	}
	return isTerminal(w)
}

func formatJSON(res gjson.Result, opts ShowJSONOpts) ([]byte, error) {
	if opts.Transform != "" {
		transformed := res.Get(opts.Transform)
		if transformed.Exists() {
			res = transformed
		}
	}
	// Match jq -r by printing strings without JSON quotes.
	if opts.RawOutput && res.Type == gjson.String {
		return []byte(res.Str + "\n"), nil
	}
	switch strings.ToLower(opts.Format) {
	case "auto":
		autoOpts := opts
		autoOpts.Format = "json"
		autoOpts.Transform = ""
		return formatJSON(res, autoOpts)
	case "pretty":
		return []byte(jsonview.RenderJSON(opts.Title, res) + "\n"), nil
	case "json":
		prettyJSON := pretty.Pretty([]byte(res.Raw))
		if shouldUseColors(opts.Stdout) {
			return pretty.Color(prettyJSON, pretty.TerminalStyle), nil
		} else {
			return prettyJSON, nil
		}
	case "jsonl":
		// @ugly removes whitespace for JSON Lines.
		oneLineJSON := res.Get("@ugly").Raw
		if shouldUseColors(opts.Stdout) {
			bytes := append(pretty.Color([]byte(oneLineJSON), pretty.TerminalStyle), '\n')
			return bytes, nil
		} else {
			return []byte(oneLineJSON + "\n"), nil
		}
	case "raw":
		return []byte(res.Raw + "\n"), nil
	case "yaml":
		input := strings.NewReader(res.Raw)
		var yaml strings.Builder
		if err := json2yaml.Convert(&yaml, input); err != nil {
			return nil, err
		}
		_, err := opts.Stdout.Write([]byte(yaml.String()))
		return nil, err
	default:
		return nil, fmt.Errorf("format %q is invalid. Choose one of: %s", opts.Format, strings.Join(OutputFormats, ", "))
	}
}

const warningExploreNotSupported = "Explore format requires a terminal. Using JSON instead.\n"

// ShowJSONOpts configures how JSON output is displayed.
type ShowJSONOpts struct {
	ExplicitFormat bool      // true if the user explicitly passed --format
	Format         string    // output format (auto, explore, json, jsonl, pretty, raw, yaml)
	RawOutput      bool      // like jq -r: print strings without JSON quotes
	Stderr         io.Writer // stderr for warnings; injectable for testing; defaults to os.Stderr
	Stdout         *os.File  // stdout (or pager); injectable for testing; defaults to os.Stdout
	Title          string    // display title
	Transform      string    // GJSON path to extract before displaying
}

func (o *ShowJSONOpts) setDefaults() {
	if o.Stderr == nil {
		o.Stderr = os.Stderr
	}
	if o.Stdout == nil {
		o.Stdout = os.Stdout
	}
}

// ShowJSON displays a single JSON result to the user.
func ShowJSON(res gjson.Result, opts ShowJSONOpts) error {
	opts.setDefaults()

	switch strings.ToLower(opts.Format) {
	case "auto":
		autoOpts := opts
		autoOpts.Format = "json"
		return ShowJSON(res, autoOpts)
	case "explore":
		if !isTerminal(opts.Stdout) {
			if opts.ExplicitFormat {
				fmt.Fprint(opts.Stderr, warningExploreNotSupported)
			}
			jsonOpts := opts
			jsonOpts.Format = "json"
			return ShowJSON(res, jsonOpts)
		}
		if opts.Transform != "" {
			transformed := res.Get(opts.Transform)
			if transformed.Exists() {
				res = transformed
			}
		}
		return jsonview.ExploreJSON(opts.Title, res)
	default:
		bytes, err := formatJSON(res, opts)
		if err != nil {
			return err
		}

		_, err = opts.Stdout.Write(bytes)
		return err
	}
}

// countTerminalLines returns the rendered terminal height.
func countTerminalLines(data []byte, terminalWidth int) int {
	return bytes.Count([]byte(wrap.String(string(data), terminalWidth)), []byte("\n"))
}

type hasRawJSON interface {
	RawJSON() string
}

// ShowJSONIterator displays streamed values. Use -1 for no limit.
func ShowJSONIterator[T any](iter jsonview.Iterator[T], itemsToDisplay int64, opts ShowJSONOpts) error {
	opts.setDefaults()

	if opts.Format == "explore" {
		if isTerminal(opts.Stdout) {
			return jsonview.ExploreJSONStream(opts.Title, iter)
		}
		if opts.ExplicitFormat {
			fmt.Fprint(opts.Stderr, warningExploreNotSupported)
		}
		opts.Format = "json"
	}

	terminalWidth, terminalHeight, err := term.GetSize(os.Stdout.Fd())
	if err != nil {
		terminalWidth = 100
		terminalHeight = 40
	}

	// Page output that exceeds the terminal.
	usePager := false
	output := []byte{}
	numberOfNewlines := 0
	for itemsToDisplay != 0 && iter.Next() {
		item := iter.Current()
		var obj gjson.Result
		if hasRaw, ok := any(item).(hasRawJSON); ok {
			obj = gjson.Parse(hasRaw.RawJSON())
		} else {
			jsonData, err := json.Marshal(item)
			if err != nil {
				return err
			}
			obj = gjson.ParseBytes(jsonData)
		}
		json, err := formatJSON(obj, opts)
		if err != nil {
			return err
		}

		output = append(output, json...)
		itemsToDisplay -= 1
		numberOfNewlines += countTerminalLines(json, terminalWidth)

		if numberOfNewlines >= terminalHeight-3 {
			usePager = true
			break
		}
	}

	if !usePager {
		_, err := opts.Stdout.Write(output)
		if err != nil {
			return err
		}

		return iter.Err()
	}

	return streamOutput(opts.Title, func(pager *os.File) error {
		_, err := pager.Write(output)
		if err != nil {
			return err
		}

		pagerOpts := opts
		pagerOpts.Stdout = pager

		for iter.Next() {
			if itemsToDisplay == 0 {
				break
			}
			item := iter.Current()
			var obj gjson.Result
			if hasRaw, ok := any(item).(hasRawJSON); ok {
				obj = gjson.Parse(hasRaw.RawJSON())
			} else {
				jsonData, err := json.Marshal(item)
				if err != nil {
					return err
				}
				obj = gjson.ParseBytes(jsonData)
			}
			if err := ShowJSON(obj, pagerOpts); err != nil {
				return err
			}
			itemsToDisplay -= 1
		}
		return iter.Err()
	})
}
