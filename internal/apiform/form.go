// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

// Package apiform encodes request bodies as multipart/form-data, which
// https://www.rfc-editor.org/rfc/rfc7578 defines.
package apiform

type Marshaler interface {
	MarshalMultipart() ([]byte, string, error)
}

type FormFormat int

const (
	FormatRepeat FormFormat = iota
	FormatComma
	FormatBrackets
	FormatIndicesDots
	FormatIndicesBrackets
)
