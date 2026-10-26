// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

// Package jsonview shows JSON as static text or in an interactive explorer built
// on https://github.com/charmbracelet/bubbletea.
package jsonview

import (
	"errors"
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
	"github.com/muesli/reflow/wordwrap"
	"github.com/tidwall/gjson"
)

const (
	borderPadding     = 2
	heightOffset      = 5
	tableMinHeight    = 2
	titlePaddingLeft  = 2
	titlePaddingTop   = 0
	footerPaddingLeft = 1

	defaultColumnWidth = 10
	keyColumnWidth     = 3
	valueColumnWidth   = 5

	maxStringLength  = 100
	maxPreviewLength = 24

	arrayColor  = lipgloss.Color("1")
	stringColor = lipgloss.Color("5")
	objectColor = lipgloss.Color("4")
)

type keyMap struct {
	Up         key.Binding
	Down       key.Binding
	Enter      key.Binding
	Back       key.Binding
	PrintValue key.Binding
	Raw        key.Binding
	Quit       key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit, k.Up, k.Down, k.Back, k.Enter, k.PrintValue, k.Raw}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	Back: key.NewBinding(
		key.WithKeys("left", "h", "backspace"),
		key.WithHelp("←/h", "back"),
	),
	Enter: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "expand"),
	),
	PrintValue: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "print & exit"),
	),
	Raw: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "raw JSON"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c", "enter"),
		key.WithHelp("q/enter", "quit"),
	),
}

var (
	titleStyle         = lipgloss.NewStyle().Bold(true).PaddingLeft(titlePaddingLeft).PaddingTop(titlePaddingTop)
	arrayStyle         = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(arrayColor)
	stringStyle        = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(stringColor)
	objectStyle        = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(objectColor)
	stringLiteralStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("2"))
)

type JSONView interface {
	GetPath() string
	GetData() gjson.Result
	Update(tea.KeyMsg, bool)
	View() string
	Resize(width, height int)
}

type TableView struct {
	width   int
	height  int
	path    string
	data    gjson.Result
	table   table.Model
	rowData []gjson.Result
	columns []table.Column
}

func (tv *TableView) GetPath() string       { return tv.path }
func (tv *TableView) GetData() gjson.Result { return tv.data }
func (tv *TableView) View() string          { return tv.table.View() }

func (tv *TableView) Update(msg tea.KeyMsg, raw bool) {
	tv.table, _ = tv.table.Update(msg)
}

func (tv *TableView) Resize(width, height int) {
	tv.width = width
	tv.height = height
	tv.updateColumnWidths(width)
	tv.table.SetHeight(min(height-heightOffset, tableMinHeight+len(tv.table.Rows())))
}

func (tv *TableView) updateColumnWidths(width int) {
	columns := tv.table.Columns()
	widths := make([]int, len(columns))

	for i, col := range columns {
		widths[i] = lipgloss.Width(col.Title)
	}

	for _, row := range tv.table.Rows() {
		for i, cell := range row {
			if i < len(widths) {
				widths[i] = max(widths[i], lipgloss.Width(cell))
			}
		}
	}

	totalWidth := sum(widths)
	available := width - borderPadding*len(columns)

	if totalWidth <= available {
		for i, w := range widths {
			columns[i].Width = w
		}
		return
	}

	fairShare := float64(available) / float64(len(columns))
	shrinkable := 0.0

	for _, w := range widths {
		if float64(w) > fairShare {
			shrinkable += float64(w) - fairShare
		}
	}

	if shrinkable > 0 {
		excess := float64(totalWidth - available)
		for i, w := range widths {
			if float64(w) > fairShare {
				reduction := (float64(w) - fairShare) * (excess / shrinkable)
				widths[i] = int(math.Round(float64(w) - reduction))
			}
		}
	}

	for i, w := range widths {
		columns[i].Width = w
	}

	tv.table.SetColumns(columns)
}

type TextView struct {
	path     string
	data     gjson.Result
	viewport viewport.Model
	ready    bool
}

func (tv *TextView) GetPath() string       { return tv.path }
func (tv *TextView) GetData() gjson.Result { return tv.data }
func (tv *TextView) View() string          { return tv.viewport.View() }

func (tv *TextView) Update(msg tea.KeyMsg, raw bool) {
	tv.viewport, _ = tv.viewport.Update(msg)
}

func (tv *TextView) Resize(width, height int) {
	h := height - heightOffset
	if !tv.ready {
		tv.viewport = viewport.New(width, h)
		tv.viewport.SetContent(wordwrap.String(tv.data.String(), width))
		tv.ready = true
		return
	}
	tv.viewport.Width = width
	tv.viewport.Height = h
}

type JSONViewer struct {
	stack   []JSONView
	root    string
	width   int
	height  int
	rawMode bool
	message string
	help    help.Model
}

func ExploreJSON(title string, json gjson.Result, options ...Option) error {
	view, err := newView("", json, false)
	if err != nil {
		return err
	}

	viewer := &JSONViewer{stack: []JSONView{view}, root: title, rawMode: false, help: help.New()}

	terminal := &screen{output: os.Stdout}
	for _, option := range options {
		option(terminal)
	}
	err = terminal.run(viewer)
	if viewer.message != "" {
		_, msgErr := fmt.Println("\n" + viewer.message)
		err = errors.Join(err, msgErr)
	}
	return err
}

func (v *JSONViewer) current() JSONView { return v.stack[len(v.stack)-1] }

func (v *JSONViewer) resize(width, height int) {
	v.width, v.height = width, height
	v.help.Width = width
	for i := range v.stack {
		v.stack[i].Resize(width, height)
	}
}

func (v *JSONViewer) handleKey(msg tea.KeyMsg) (bool, error) {
	switch {
	case key.Matches(msg, keys.Quit):
		return true, nil
	case key.Matches(msg, keys.Enter):
		v.navigateForward()
	case key.Matches(msg, keys.Back):
		v.navigateBack()
	case key.Matches(msg, keys.Raw):
		return false, v.toggleRaw()
	case key.Matches(msg, keys.PrintValue):
		v.message = v.getSelectedContent()
		return true, nil
	default:
		v.current().Update(msg, v.rawMode)
	}
	return false, nil
}

func (v *JSONViewer) getSelectedContent() string {
	tableView, ok := v.current().(*TableView)
	if !ok {
		return v.current().GetData().Raw
	}

	selected := tableView.rowData[tableView.table.Cursor()]
	if selected.Type == gjson.String {
		return selected.String()
	}
	return selected.Raw
}

func (v *JSONViewer) navigateForward() {
	tableView, ok := v.current().(*TableView)
	if !ok {
		return
	}

	if len(tableView.rowData) < 1 {
		return
	}

	cursor := tableView.table.Cursor()
	selected := tableView.rowData[cursor]
	if !v.canNavigateInto(selected) {
		return
	}

	path := v.buildNavigationPath(tableView, cursor)
	forwardView, err := newView(path, selected, v.rawMode)
	if err != nil {
		return
	}

	v.stack = append(v.stack, forwardView)
	v.resize(v.width, v.height)
}

func (v *JSONViewer) buildNavigationPath(tableView *TableView, cursor int) string {
	if tableView.data.IsArray() {
		return fmt.Sprintf("%s[%d]", tableView.path, cursor)
	}
	key := tableView.data.Get("@keys").Array()[cursor].Str
	return fmt.Sprintf("%s[%s]", tableView.path, quoteString(key))
}

func quoteString(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	return stringLiteralStyle.Render("\"" + s + "\"")
}

func (v *JSONViewer) canNavigateInto(data gjson.Result) bool {
	switch {
	case data.IsArray():
		return len(data.Array()) > 0
	case data.IsObject():
		return len(data.Map()) > 0
	case data.Type == gjson.String:
		str := data.String()
		return strings.Contains(str, "\n") || lipgloss.Width(str) >= maxStringLength
	}
	return false
}

func (v *JSONViewer) navigateBack() {
	if len(v.stack) > 1 {
		v.stack = v.stack[:len(v.stack)-1]
	}
}

func (v *JSONViewer) toggleRaw() error {
	v.rawMode = !v.rawMode

	for i, view := range v.stack {
		viewWithRaw, err := newView(view.GetPath(), view.GetData(), v.rawMode)
		if err != nil {
			return err
		}
		v.stack[i] = viewWithRaw
	}

	v.resize(v.width, v.height)
	return nil
}

func (v *JSONViewer) View() string {
	view := v.current()
	title := v.buildTitle(view)
	content := titleStyle.Render(title)
	style := v.getStyleForData(view.GetData())
	content += "\n" + style.Render(view.View())
	content += "\n" + v.help.View(keys)
	return content
}

func (v *JSONViewer) buildTitle(view JSONView) string {
	title := v.root
	if len(view.GetPath()) > 0 {
		title += " → " + view.GetPath()
	}
	if v.rawMode {
		title += " (JSON)"
	}
	return title
}

func (v *JSONViewer) getStyleForData(data gjson.Result) lipgloss.Style {
	switch {
	case data.Type == gjson.String:
		return stringStyle
	case data.IsArray():
		return arrayStyle
	default:
		return objectStyle
	}
}

func newView(path string, data gjson.Result, raw bool) (JSONView, error) {
	if data.Type == gjson.String {
		return newTextView(path, data)
	}
	return newTableView(path, data, raw)
}

func newTextView(path string, data gjson.Result) (*TextView, error) {
	if !data.Exists() || data.Type != gjson.String {
		return nil, fmt.Errorf("invalid text JSON")
	}
	return &TextView{path: path, data: data}, nil
}

func newTableView(path string, data gjson.Result, raw bool) (*TableView, error) {
	if !data.Exists() || data.Type != gjson.JSON {
		return nil, fmt.Errorf("invalid table JSON")
	}

	switch {
	case data.IsArray():
		array := data.Array()
		if isArrayOfObjects(array) {
			return newArrayOfObjectsTableView(path, data, array, raw), nil
		} else {
			return newArrayTableView(path, data, array, raw), nil
		}
	case data.IsObject():
		return newObjectTableView(path, data, raw), nil
	default:
		return nil, fmt.Errorf("unsupported JSON type")
	}
}

func newArrayTableView(path string, data gjson.Result, array []gjson.Result, raw bool) *TableView {
	columns := []table.Column{{Title: "Items", Width: defaultColumnWidth}}
	rows := make([]table.Row, 0, len(array))
	rowData := make([]gjson.Result, 0, len(array))

	for _, item := range array {
		rows = append(rows, table.Row{formatValue(item, raw)})
		rowData = append(rowData, item)
	}

	t := createTable(columns, rows, arrayColor)
	return &TableView{
		path:    path,
		data:    data,
		table:   t,
		rowData: rowData,
		columns: columns,
	}
}

func newArrayOfObjectsTableView(path string, data gjson.Result, array []gjson.Result, raw bool) *TableView {
	keySet := make(map[string]struct{})
	var columns []table.Column

	for _, item := range array {
		for _, key := range item.Get("@keys").Array() {
			if _, exists := keySet[key.Str]; !exists {
				keySet[key.Str] = struct{}{}
				title := key.Str
				columns = append(columns, table.Column{Title: title, Width: defaultColumnWidth})
			}
		}
	}

	rows := make([]table.Row, 0, len(array))
	rowData := make([]gjson.Result, 0, len(array))

	for _, item := range array {
		row := make(table.Row, len(columns))
		for i, col := range columns {
			row[i] = formatValue(item.Get(col.Title), raw)
		}
		rows = append(rows, row)
		rowData = append(rowData, item)
	}

	t := createTable(columns, rows, arrayColor)
	return &TableView{
		path:    path,
		data:    data,
		table:   t,
		rowData: rowData,
		columns: columns,
	}
}

func newObjectTableView(path string, data gjson.Result, raw bool) *TableView {
	columns := []table.Column{{Title: "Object"}, {}}

	keys := data.Get("@keys").Array()
	rows := make([]table.Row, 0, len(keys))
	rowData := make([]gjson.Result, 0, len(keys))

	for _, key := range keys {
		value := data.Get(key.Str)
		title := key.Str
		rows = append(rows, table.Row{title, formatValue(value, raw)})
		rowData = append(rowData, value)
	}

	for _, row := range rows {
		for i, cell := range row {
			if i < len(columns) {
				columns[i].Width = max(columns[i].Width, lipgloss.Width(cell))
			}
		}
	}

	t := createTable(columns, rows, objectColor)
	return &TableView{
		path:    path,
		data:    data,
		table:   t,
		rowData: rowData,
		columns: columns,
	}
}

func createTable(columns []table.Column, rows []table.Row, bgColor lipgloss.Color) table.Model {
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(bgColor).
		Bold(false)
	t.SetStyles(s)

	return t
}

func formatValue(value gjson.Result, raw bool) string {
	if raw {
		return value.Get("@ugly").Raw
	}

	switch {
	case value.IsObject():
		return formatObject(value)
	case value.IsArray():
		return formatArray(value)
	case value.Type == gjson.String:
		return value.Str
	default:
		return value.Raw
	}
}

func formatObject(value gjson.Result) string {
	keys := value.Get("@keys").Array()
	keyStrs := make([]string, len(keys))

	for i, key := range keys {
		val := value.Get(key.Str)
		keyStrs[i] = formatObjectKey(key.Str, val)
	}

	return "{" + strings.Join(keyStrs, ", ") + "}"
}

func formatObjectKey(key string, val gjson.Result) string {
	switch {
	case val.IsObject():
		return key + ":{…}"
	case val.IsArray():
		return key + ":[…]"
	case val.Type == gjson.String:
		str := val.Str
		if lipgloss.Width(str) <= maxPreviewLength {
			return fmt.Sprintf(`%s:"%s"`, key, str)
		}
		return fmt.Sprintf(`%s:"%s…"`, key, truncate.String(str, uint(maxPreviewLength)))
	default:
		return key + ":" + val.Raw
	}
}

func formatArray(value gjson.Result) string {
	switch count := len(value.Array()); count {
	case 0:
		return "[]"
	case 1:
		return "[1 item]"
	default:
		return fmt.Sprintf("[%d items]", count)
	}
}

func isArrayOfObjects(array []gjson.Result) bool {
	for _, item := range array {
		if !item.IsObject() {
			return false
		}
	}
	return len(array) > 0
}

func sum(ints []int) int {
	total := 0
	for _, n := range ints {
		total += n
	}
	return total
}
