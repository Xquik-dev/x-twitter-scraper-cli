// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

package jsonview

import (
	"testing"

	"github.com/charmbracelet/bubbles/help"
	"github.com/tidwall/gjson"

	"github.com/stretchr/testify/require"
)

func TestNavigateForward_EmptyRowData(t *testing.T) {
	t.Parallel()

	emptyArray := gjson.Parse("[]")
	view, err := newTableView("", emptyArray, false)
	require.NoError(t, err)

	viewer := &JSONViewer{
		stack: []JSONView{view},
		root:  "test",
		help:  help.New(),
	}

	model, cmd := viewer.navigateForward()
	require.Equal(t, model, viewer, "expected same viewer model returned")
	require.Nil(t, cmd)

	require.Equal(t, 1, len(viewer.stack), "expected stack length 1, got %d", len(viewer.stack))
}
