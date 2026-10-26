// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Xquik-dev/x-twitter-scraper-cli/internal/jsonview"
	"github.com/Xquik-dev/x-twitter-scraper-go/option"

	"github.com/charmbracelet/x/term"
	"github.com/itchyny/json2yaml"
	"github.com/tidwall/gjson"
	"github.com/tidwall/pretty"
	"github.com/urfave/cli/v3"
)

var OutputFormats = []string{"auto", "explore", "json", "jsonl", "pretty", "raw", "yaml"}

func formatChoices() string {
	last := len(OutputFormats) - 1
	return strings.Join(OutputFormats[:last], ", ") + " or " + OutputFormats[last]
}

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

func isInputPiped() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	mode := stat.Mode()

	if mode.IsRegular() && stat.Size() > 0 {
		return true
	}

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
		if !isTerminal(os.Stdout) || isUTF8TextFile(body) {
			_, err := stdout.Write(body)
			return "", err
		}

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

func createDownloadFile(response *http.Response, data []byte) (*os.File, error) {
	filename := "file"
	disp := response.Header.Get("Content-Disposition")
	_, params, err := mime.ParseMediaType(disp)
	if err == nil {
		if dispFilename, ok := params["filename"]; ok {
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
		return nil, unknownFormatError(opts.Format)
	}
}

type unknownFormatError string

func (e unknownFormatError) Error() string {
	return fmt.Sprintf("Format %q is unknown. Choose %s.", string(e), formatChoices())
}

type unexpectedArgumentsError []string

func (e unexpectedArgumentsError) Error() string {
	return fmt.Sprintf("Unexpected extra arguments %v. Remove them, then run the command again.", []string(e))
}

const warningExploreNotSupported = "Explore format requires a terminal. Using JSON instead.\n"

type ShowJSONOpts struct {
	ExplicitFormat bool
	Format         string
	RawOutput      bool
	Stderr         io.Writer
	Stdout         *os.File
	Title          string
	Transform      string
}

func (o *ShowJSONOpts) setDefaults() {
	if o.Stderr == nil {
		o.Stderr = os.Stderr
	}
	if o.Stdout == nil {
		o.Stdout = os.Stdout
	}
}

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
