// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

package jsonview

import (
	"errors"
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

type sliceAnyIterator struct {
	items   []any
	current int
	err     error
}

func (iterator *sliceAnyIterator) Next() bool {
	if iterator.current >= len(iterator.items) {
		return false
	}
	iterator.current++
	return true
}

func (iterator *sliceAnyIterator) Err() error {
	return iterator.err
}

func (iterator *sliceAnyIterator) Current() any {
	return iterator.items[iterator.current-1]
}

type invalidView struct {
	data gjson.Result
}

func (view *invalidView) GetPath() string              { return "invalid" }
func (view *invalidView) GetData() gjson.Result        { return view.data }
func (view *invalidView) Update(tea.Msg, bool) tea.Cmd { return nil }
func (view *invalidView) View() string                 { return "invalid" }
func (view *invalidView) Resize(width, height int)     {}

func keyMessage(value string) tea.KeyMsg {
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(value)}
}

func TestKeyMapAndStaticFormatting(t *testing.T) {
	t.Parallel()

	require.Len(t, keys.ShortHelp(), 7)
	require.Len(t, keys.FullHelp(), 1)
	assert.Contains(t, formatJSON(gjson.Result{}, 80), "Invalid JSON")
	assert.Contains(t, formatJSON(gjson.Parse(`{"name":"xquik"}`), 80), "name")
	assert.Contains(t, formatResult(gjson.Parse(`"value"`), 0, 80), "value")
	assert.Contains(t, formatResult(gjson.Parse(`""`), 0, 80), "empty")
	assert.NotEmpty(t, formatResult(gjson.Parse(`"a very long value"`), 0, 5))
	assert.NotEmpty(t, formatResult(gjson.Parse("42"), 0, 80))
	assert.NotEmpty(t, formatResult(gjson.Parse("true"), 0, 80))
	assert.NotEmpty(t, formatResult(gjson.Parse("false"), 0, 80))
	assert.NotEmpty(t, formatResult(gjson.Parse("null"), 0, 80))
	assert.NotEmpty(t, formatResult(gjson.Parse(`[1,{"a":2}]`), 0, 80))
	assert.Empty(t, formatResult(gjson.Result{Type: gjson.Type(99)}, 0, 80))
	assert.True(t, isSingleLine(gjson.Parse("1"), 0))
	assert.False(t, isSingleLine(gjson.Parse("[1]"), 0))
	assert.Contains(t, formatJSONArray(gjson.Parse("[]"), 0, 80), "none")
	assert.NotEmpty(t, formatJSONArray(gjson.Parse(`[1,{"a":2}]`), 0, 80))
	assert.Contains(t, formatJSONObject(gjson.Parse("{}"), 0, 80), "empty")
	assert.NotEmpty(t, formatJSONObject(gjson.Parse(`{"a":1,"b":[2]}`), 0, 80))
	assert.Equal(t, "    ", getIndent(2))
	assert.Contains(t, RenderJSON("Title", gjson.Parse(`{"a":1}`)), "Title")
}

func TestViewConstructionAndFormatting(t *testing.T) {
	t.Parallel()

	text, err := newTextView("text", gjson.Parse(`"line one\nline two"`))
	require.NoError(t, err)
	require.Equal(t, "text", text.GetPath())
	require.Equal(t, "line one\nline two", text.GetData().String())
	text.Resize(30, 10)
	assert.True(t, text.ready)
	text.Resize(40, 12)
	assert.Equal(t, 40, text.viewport.Width)
	_ = text.View()
	assert.Nil(t, text.Update(keyMessage("j"), false))

	_, err = newTextView("", gjson.Parse("1"))
	require.Error(t, err)
	_, err = newTableView("", gjson.Parse("1"), false)
	require.Error(t, err)
	_, err = newView("", gjson.Result{}, false)
	require.Error(t, err)

	primitive, err := newTableView("array", gjson.Parse(`[1,"two",{"a":3}]`), false)
	require.NoError(t, err)
	require.Equal(t, "array", primitive.GetPath())
	require.Len(t, primitive.rowData, 3)
	primitive.Resize(20, 8)
	primitive.Resize(200, 20)
	_ = primitive.View()
	_ = primitive.GetData()
	assert.Nil(t, primitive.Update(keyMessage("j"), false))

	objects, err := newTableView("", gjson.Parse(`[{"a":1},{"a":2,"b":"x"}]`), false)
	require.NoError(t, err)
	require.Len(t, objects.columns, 2)

	object, err := newTableView("", gjson.Parse(`{"a":1,"b":[2]}`), false)
	require.NoError(t, err)
	require.Len(t, object.rowData, 2)

	assert.Equal(t, "raw", formatValue(gjson.Parse(`"raw"`), false))
	assert.NotEmpty(t, formatValue(gjson.Parse(`{"a":1}`), true))
	assert.Contains(t, formatObject(gjson.Parse(`{"nested":{"a":1},"items":[1],"short":"x","long":"abcdefghijklmnopqrstuvwxyz"}`)), "nested:{…}")
	assert.Equal(t, "[]", formatArray(gjson.Parse("[]")))
	assert.Equal(t, "[1 item]", formatArray(gjson.Parse("[1]")))
	assert.Equal(t, "[2 items]", formatArray(gjson.Parse("[1,2]")))
	assert.True(t, isArrayOfObjects(gjson.Parse(`[{"a":1}]`).Array()))
	assert.False(t, isArrayOfObjects(gjson.Parse(`[{"a":1},2]`).Array()))
	assert.False(t, isArrayOfObjects(nil))
	assert.Equal(t, 6, sum([]int{1, 2, 3}))
	assert.Contains(t, quoteString(`a\"b`), `\\`)
}

func TestTableViewLoadsIteratorResults(t *testing.T) {
	t.Parallel()

	view, err := newTableView("", gjson.Parse(`[{"a":1,"b":2}]`), false)
	require.NoError(t, err)
	view.iterator = &sliceAnyIterator{items: []any{map[string]any{"a": 3, "b": 4}}}
	view.Resize(80, 20)
	command := view.Update(keyMessage("j"), false)
	require.NotNil(t, command)
	require.Nil(t, command())
	require.Len(t, view.rowData, 2)
	require.Len(t, view.table.Rows(), 2)

	view.iterator = &sliceAnyIterator{err: errors.New("iterator failed")}
	view.isLoading = false
	view.table.SetCursor(len(view.table.Rows()) - 1)
	command = view.Update(keyMessage("j"), false)
	require.EqualError(t, command().(error), "iterator failed")

	view.iterator = &sliceAnyIterator{items: []any{make(chan int)}}
	view.isLoading = false
	command = view.loadMoreData(false)
	require.Error(t, command().(error))

	view.iterator = &sliceAnyIterator{items: []any{nil}}
	view.isLoading = false
	command = view.loadMoreData(false)
	require.Nil(t, command())

	view.iterator = nil
	require.Nil(t, view.loadMoreData(false)())
}

func TestJSONViewerNavigationAndModes(t *testing.T) {
	t.Parallel()

	root, err := newTableView("", gjson.Parse(`{"nested":{"value":1},"plain":"x"}`), false)
	require.NoError(t, err)
	viewer := &JSONViewer{
		stack: []JSONView{root},
		root:  "Root",
		help:  help.New(),
	}

	require.Nil(t, viewer.Init())
	model, command := viewer.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	require.Same(t, viewer, model)
	require.Nil(t, command)
	assert.Equal(t, 98, viewer.width)
	assert.NotEmpty(t, viewer.View())
	assert.Equal(t, "Root", viewer.buildTitle(root))

	model, command = viewer.Update(keyMessage("l"))
	require.Same(t, viewer, model)
	require.Nil(t, command)
	require.Len(t, viewer.stack, 2)
	assert.Contains(t, viewer.buildTitle(viewer.current()), "nested")

	model, command = viewer.Update(keyMessage("h"))
	require.Same(t, viewer, model)
	require.Nil(t, command)
	require.Len(t, viewer.stack, 1)

	model, command = viewer.Update(keyMessage("r"))
	require.Same(t, viewer, model)
	require.Nil(t, command)
	assert.True(t, viewer.rawMode)
	assert.Contains(t, viewer.buildTitle(viewer.current()), "JSON")

	model, command = viewer.Update(keyMessage("p"))
	require.Same(t, viewer, model)
	require.NotNil(t, command)
	assert.NotEmpty(t, viewer.message)

	model, command = viewer.Update(keyMessage("q"))
	require.Same(t, viewer, model)
	require.NotNil(t, command)

	model, command = viewer.Update(keyMessage("j"))
	require.Same(t, viewer, model)
	_ = command

	text, err := newTextView("", gjson.Parse(`"content"`))
	require.NoError(t, err)
	viewer.stack = []JSONView{text}
	assert.Equal(t, `"content"`, viewer.getSelectedContent())
	assert.True(t, viewer.canNavigateInto(gjson.Parse(`"line one\nline two"`)))
	assert.True(t, viewer.canNavigateInto(gjson.Parse(`"`+strings.Repeat("x", maxStringLength)+`"`)))
	assert.True(t, viewer.canNavigateInto(gjson.Parse(`[1]`)))
	assert.True(t, viewer.canNavigateInto(gjson.Parse(`{"a":1}`)))
	assert.False(t, viewer.canNavigateInto(gjson.Parse(`[]`)))
	assert.False(t, viewer.canNavigateInto(gjson.Parse(`{}`)))
	assert.False(t, viewer.canNavigateInto(gjson.Parse(`1`)))

	viewer.stack = []JSONView{&invalidView{data: gjson.Parse("1")}}
	_, command = viewer.toggleRaw()
	require.NotNil(t, command)

	assert.Equal(t, stringStyle, viewer.getStyleForData(gjson.Parse(`"x"`)))
	assert.Equal(t, arrayStyle, viewer.getStyleForData(gjson.Parse(`[1]`)))
	assert.Equal(t, objectStyle, viewer.getStyleForData(gjson.Parse(`{"a":1}`)))
}

func TestIteratorAdaptersAndEarlyExplorerErrors(t *testing.T) {
	t.Parallel()

	source := &sliceAnyIterator{items: []any{"first"}, err: errors.New("done")}
	adapted := genericToAnyIterator[any](source)
	require.True(t, adapted.Next())
	assert.Equal(t, "first", adapted.Current())
	require.False(t, adapted.Next())
	assert.EqualError(t, adapted.Err(), "done")

	badIterator := &sliceAnyIterator{err: errors.New("stream failed")}
	require.EqualError(t, ExploreJSONStream[any]("test", badIterator), "stream failed")

	marshalIterator := &sliceAnyIterator{items: []any{make(chan int)}}
	require.Error(t, ExploreJSONStream[any]("test", marshalIterator))
	require.Error(t, ExploreJSON("test", gjson.Result{}))
}

func TestTextViewUpdateUsesViewport(t *testing.T) {
	t.Parallel()

	text := &TextView{
		data:     gjson.Parse(`"content"`),
		viewport: viewport.New(20, 4),
		ready:    true,
	}
	assert.Nil(t, text.Update(keyMessage("j"), false))
}

func TestTableViewNarrowColumns(t *testing.T) {
	t.Parallel()

	model := table.New(
		table.WithColumns([]table.Column{{Title: "long-header"}, {Title: "b"}}),
		table.WithRows([]table.Row{{"long-value", "another-long-value"}}),
	)
	view := &TableView{table: model, columns: model.Columns()}
	view.Resize(8, 4)
	require.Len(t, view.table.Columns(), 2)
}
