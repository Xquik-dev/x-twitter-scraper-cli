// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

// Package apiquery encodes request parameters as the URL query values of
// https://pkg.go.dev/net/url#Values.
package apiquery

import (
	"net/url"
	"strings"
)

type ArrayQueryFormat int

const (
	ArrayQueryFormatComma ArrayQueryFormat = iota
	ArrayQueryFormatRepeat
)

type Param struct {
	values []string
	list   bool
}

func Scalar(value string) Param {
	return Param{values: []string{value}}
}

func List(values []string) Param {
	return Param{values: values, list: true}
}

func (p Param) Map(transform func(string) (string, error)) (Param, error) {
	values := make([]string, len(p.values))
	for i, value := range p.values {
		transformed, err := transform(value)
		if err != nil {
			return Param{}, err
		}
		values[i] = transformed
	}
	return Param{values: values, list: p.list}, nil
}

func Encode(params map[string]Param, format ArrayQueryFormat) url.Values {
	encoded := url.Values{}
	for key, param := range params {
		if param.list && format == ArrayQueryFormatComma {
			encoded.Add(key, strings.Join(param.values, ","))
			continue
		}
		for _, value := range param.values {
			encoded.Add(key, value)
		}
	}
	return encoded
}
