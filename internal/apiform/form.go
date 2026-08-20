// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

package apiform

type Marshaler interface {
	MarshalMultipart() ([]byte, string, error)
}

type FormFormat int

const (
	// FormatRepeat uses repeated keys.
	FormatRepeat FormFormat = iota
	// FormatComma joins values with commas.
	FormatComma
	// FormatBrackets uses key[].
	FormatBrackets
	// FormatIndicesDots uses key.0 and key.1.
	FormatIndicesDots
	// FormatIndicesBrackets uses key[0] and key[1].
	FormatIndicesBrackets
)
