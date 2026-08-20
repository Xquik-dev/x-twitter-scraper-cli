// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

package requestflag

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/goccy/go-yaml"
	"github.com/urfave/cli/v3"
)

// formatForFlagSet prepares piped values for flag.Set.
// Strings stay raw, scalars use %v, and complex values use JSON.
func formatForFlagSet(val any) (string, error) {
	switch v := val.(type) {
	case string:
		return v, nil
	case bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return fmt.Sprintf("%v", val), nil
	default:
		b, err := json.Marshal(val)
		if err != nil {
			return "", fmt.Errorf("cannot format value %T for flag.Set: %w", val, err)
		}
		return string(b), nil
	}
}

// Flag implements urfave/cli flags and maps values into HTTP requests.
// Pointer types preserve unset, null, and concrete values.
type Flag[
	T []any | []map[string]any | []DateTimeValue | []DateValue | []TimeValue | []string |
		[]float64 | []int64 | []bool | any | map[string]any | DateTimeValue | DateValue | TimeValue |
		string | float64 | int64 | bool |
		*string | *float64 | *int64 | *bool | *DateTimeValue | *DateValue | *TimeValue,
] struct {
	Name        string
	Category    string
	DefaultText string
	HideDefault bool
	Usage       string
	Sources     cli.ValueSourceChain
	Required    bool
	Hidden      bool
	Default     T
	Aliases     []string
	Validator   func(T) error

	QueryPath  string
	HeaderPath string
	BodyPath   string
	BodyRoot   bool
	PathParam  string

	// Const always includes Default but still permits overrides.
	Const bool

	// FileInput treats values as paths without an @ prefix.
	FileInput bool

	// DataAliases maps piped input names to the canonical API name.
	DataAliases []string

	count      int
	hasBeenSet bool
	applied    bool
	value      cli.Value
}

var _ cli.CategorizableFlag = (*Flag[any])(nil)

// InRequest exposes a flag's HTTP request location.
type InRequest interface {
	GetQueryPath() string
	GetHeaderPath() string
	GetBodyPath() string
	GetPathParam() string
	IsBodyRoot() bool
	IsFileInput() bool
	GetDataAliases() []string
}

func (f Flag[T]) GetQueryPath() string {
	return f.QueryPath
}

func (f Flag[T]) GetHeaderPath() string {
	return f.HeaderPath
}

func (f Flag[T]) GetBodyPath() string {
	return f.BodyPath
}

func (f Flag[T]) GetPathParam() string {
	return f.PathParam
}

func (f Flag[T]) IsBodyRoot() bool {
	return f.BodyRoot
}

func (f Flag[T]) IsFileInput() bool {
	return f.FileInput
}

func (f Flag[T]) GetDataAliases() []string {
	return f.DataAliases
}

// RequestContents stores values by HTTP request location.
type RequestContents struct {
	Queries map[string]any
	Headers map[string]any
	Body    any
}

// ApplyStdinDataToFlags maps piped data to unset path, query, and header flags.
// Body flags use the separate body merge. Nested flags read their outer map.
func ApplyStdinDataToFlags(cmd *cli.Command, data map[string]any) error {
	for _, flag := range cmd.Flags {
		if flag.IsSet() {
			continue
		}

		if inner, ok := flag.(HasOuterFlag); ok {
			outer, outerOk := inner.GetOuterFlag().(InRequest)
			if !outerOk || outer.GetBodyPath() == "" {
				continue
			}
			nested, ok := data[outer.GetBodyPath()].(map[string]any)
			if !ok {
				continue
			}
			innerField := inner.GetInnerField()
			val, found := nested[innerField]
			if !found {
				for _, alias := range inner.GetDataAliases() {
					if alias != "" && alias != innerField {
						if v, ok := nested[alias]; ok {
							val, found = v, true
							break
						}
					}
				}
			}
			if !found {
				continue
			}
			setVal, err := formatForFlagSet(val)
			if err != nil {
				return fmt.Errorf("cannot format piped value for flag %q: %w", flag.Names()[0], err)
			}
			if err := flag.Set(flag.Names()[0], setVal); err != nil {
				return fmt.Errorf("cannot set flag %q from piped data: %w", flag.Names()[0], err)
			}
			continue
		}

		inReq, ok := flag.(InRequest)
		if !ok {
			continue
		}

		// Body parameters use the body merge instead.
		for _, path := range []string{inReq.GetQueryPath(), inReq.GetHeaderPath(), inReq.GetPathParam()} {
			if path == "" {
				continue
			}
			var val any
			var found bool
			for _, key := range append([]string{path}, inReq.GetDataAliases()...) {
				if v, ok := data[key]; ok {
					val, found = v, true
					break
				}
			}
			if !found {
				continue
			}
			setVal, err := formatForFlagSet(val)
			if err != nil {
				return fmt.Errorf("cannot format piped value for flag %q: %w", flag.Names()[0], err)
			}
			if err := flag.Set(flag.Names()[0], setVal); err != nil {
				return fmt.Errorf("cannot set flag %q from piped data: %w", flag.Names()[0], err)
			}
			break
		}
	}
	return nil
}

func ExtractRequestContents(cmd *cli.Command) RequestContents {
	bodyMap := make(map[string]any)
	res := RequestContents{
		Queries: make(map[string]any),
		Headers: make(map[string]any),
		Body:    bodyMap,
	}

	for _, flag := range cmd.Flags {
		if !flag.IsSet() {
			continue
		}

		value := flag.Get()
		if toSend, ok := flag.(InRequest); ok {
			if queryPath := toSend.GetQueryPath(); queryPath != "" {
				res.Queries[queryPath] = value
			}
			if headerPath := toSend.GetHeaderPath(); headerPath != "" {
				res.Headers[headerPath] = value
			}
			if toSend.IsBodyRoot() {
				res.Body = value
			} else if bodyPath := toSend.GetBodyPath(); bodyPath != "" {
				bodyMap[bodyPath] = value
			}
		}
	}
	return res
}

func GetMissingRequiredFlags(cmd *cli.Command, body any) []cli.Flag {
	missing := []cli.Flag{}
	for _, flag := range cmd.Flags {
		if flag.IsSet() {
			continue
		}

		if required, ok := flag.(cli.RequiredFlag); ok && required.IsRequired() {
			missing = append(missing, flag)
			continue
		}

		if r, ok := flag.(RequiredFlagOrStdin); !ok || !r.IsRequiredAsFlagOrStdin() {
			continue
		}

		if toSend, ok := flag.(InRequest); ok {
			if toSend.IsBodyRoot() {
				if body != nil {
					continue
				}
			} else if bodyPath := toSend.GetBodyPath(); bodyPath != "" {
				if bodyMap, ok := body.(map[string]any); ok {
					if _, found := bodyMap[bodyPath]; found {
						continue
					}
				}
			}
		}
		missing = append(missing, flag)
	}
	return missing
}

var _ cli.Flag = (*Flag[any])(nil)

func (f *Flag[T]) PreParse() error {
	newVal := f.Default
	f.value = &cliValue[T]{newVal}

	// Validate defaults and external values.
	if f.Validator != nil {
		if err := f.Validator(f.value.Get().(T)); err != nil {
			return err
		}
	}
	f.applied = true
	return nil
}

func (f *Flag[T]) PostParse() error {
	if !f.hasBeenSet {
		if val, source, found := f.Sources.LookupWithSource(); found {
			defaultType := reflect.TypeOf(f.Default)
			if defaultType != nil && defaultType.Kind() == reflect.Pointer {
				defaultType = defaultType.Elem()
			}
			isString := defaultType != nil && defaultType.Kind() == reflect.String
			isBool := defaultType != nil && defaultType.Kind() == reflect.Bool

			if val != "" || isString {
				if err := f.Set(f.Name, val); err != nil {
					return fmt.Errorf(
						"could not parse %[1]q as %[2]T value from %[3]s for flag %[4]s: %[5]s",
						val, f.value, source, f.Name, err,
					)
				}
			} else if isBool {
				_ = f.Set(f.Name, "false")
			}

			f.hasBeenSet = true
		}
	}
	return nil
}

func (f *Flag[T]) Set(name string, val string) error {
	if !f.applied {
		if err := f.PreParse(); err != nil {
			return err
		}
		f.applied = true
	}

	f.count++

	// Do not append the first explicit value to a default slice.
	if f.count == 1 && f.value != nil {
		typ := reflect.TypeOf(f.Default)
		if typ != nil && typ.Kind() == reflect.Slice {
			emptySlice := reflect.MakeSlice(typ, 0, 0).Interface()
			f.value = &cliValue[T]{emptySlice.(T)}
		}
	}

	if err := f.value.Set(val); err != nil {
		return err
	}

	f.hasBeenSet = true

	if f.Validator != nil {
		if err := f.Validator(f.value.Get().(T)); err != nil {
			return err
		}
	}
	return nil
}

func (f *Flag[T]) Get() any {
	if f.value != nil {
		return f.value.Get()
	}
	return f.Default
}

func (f *Flag[T]) String() string {
	return cli.FlagStringer(f)
}

func (f *Flag[T]) IsSet() bool {
	return f.hasBeenSet || f.Const
}

func (f *Flag[T]) Names() []string {
	return cli.FlagNames(f.Name, f.Aliases)
}

var _ cli.VisibleFlag = (*Flag[any])(nil)

func (f *Flag[T]) IsVisible() bool {
	return !f.Hidden
}

func (f *Flag[T]) GetCategory() string {
	return f.Category
}

func (f *Flag[T]) SetCategory(c string) {
	f.Category = c
}

var _ cli.RequiredFlag = (*Flag[any])(nil)

func (f *Flag[T]) IsRequired() bool {
	// Const flags never require user input.
	if f.Const {
		return false
	}
	// Stdin may satisfy required request flags.
	if f.BodyPath != "" || f.BodyRoot || f.PathParam != "" || f.QueryPath != "" || f.HeaderPath != "" {
		return false
	}
	return f.Required
}

type RequiredFlagOrStdin interface {
	IsRequiredAsFlagOrStdin() bool
}

func (f *Flag[T]) IsRequiredAsFlagOrStdin() bool {
	// Const flags never require user input.
	if f.Const {
		return false
	}
	return f.Required
}

var _ cli.DocGenerationFlag = (*Flag[any])(nil)

func (f *Flag[T]) TakesValue() bool {
	var t T
	return reflect.TypeOf(t) == nil || reflect.TypeOf(t).Kind() != reflect.Bool
}

func (f *Flag[T]) GetUsage() string {
	return f.Usage
}

func (f *Flag[T]) GetValue() string {
	if f.value == nil {
		return ""
	}
	return f.value.String()
}

func (f *Flag[T]) GetDefaultText() string {
	return f.DefaultText
}

// GetEnvVars returns the flag's environment variables.
func (f *Flag[T]) GetEnvVars() []string {
	return f.Sources.EnvKeys()
}

func (f *Flag[T]) IsDefaultVisible() bool {
	return !f.HideDefault
}

func (f *Flag[T]) TypeName() string {
	ty := reflect.TypeOf(f.Default)
	if ty == nil {
		return ""
	}
	// Show the pointee type in help output.
	if ty.Kind() == reflect.Pointer {
		ty = ty.Elem()
	}

	getTypeName := func(t reflect.Type) string {
		switch t.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return "int"
		case reflect.Float32, reflect.Float64:
			return "float"
		case reflect.Bool:
			return "boolean"
		case reflect.String:
			switch t.Name() {
			case "DateTimeValue":
				return "datetime"
			case "DateValue":
				return "date"
			case "TimeValue":
				return "time"
			default:
				return "string"
			}
		default:
			if t.Name() == "" {
				return "any"
			}
			return strings.ToLower(t.Name())
		}
	}

	switch ty.Kind() {
	case reflect.Slice:
		elemType := ty.Elem()
		return getTypeName(elemType)
	case reflect.Map:
		keyType := ty.Key()
		valueType := ty.Elem()
		return fmt.Sprintf("%s=%s", getTypeName(keyType), getTypeName(valueType))
	default:
		return getTypeName(ty)
	}
}

var _ cli.DocGenerationMultiValueFlag = (*Flag[any])(nil)

func (f *Flag[T]) IsMultiValueFlag() bool {
	if reflect.TypeOf(f.Default) == nil {
		return false
	}
	kind := reflect.TypeOf(f.Default).Kind()
	return kind == reflect.Slice || kind == reflect.Map
}

func (f *Flag[T]) IsBoolFlag() bool {
	// Pointer booleans need an explicit value to preserve 3 states.
	_, isBool := any(f.Default).(bool)
	return isBool
}

var _ cli.Countable = (*Flag[any])(nil)

func (f *Flag[T]) Count() int {
	return f.count
}

var _ cli.LocalFlag = (*Flag[any])(nil)

func (f Flag[T]) IsLocal() bool {
	// Request flags work at any command level.
	return true
}

// cliValue implements cli.Value for supported request types.
type cliValue[
	T []any | []map[string]any | []DateTimeValue | []DateValue | []TimeValue | []string | []float64 |
		[]int64 | []bool | any | map[string]any | DateTimeValue | DateValue | TimeValue | string |
		float64 | int64 | bool |
		*string | *float64 | *int64 | *bool | *DateTimeValue | *DateValue | *TimeValue,
] struct {
	value T
}

// parseCLIArg converts one argument into type T.
func parseCLIArg[
	T []any | []map[string]any | []DateTimeValue | []DateValue | []TimeValue | []string | []float64 |
		[]int64 | []bool | any | map[string]any | DateTimeValue | DateValue | TimeValue | string |
		float64 | int64 | bool |
		*string | *float64 | *int64 | *bool | *DateTimeValue | *DateValue | *TimeValue,
](value string) (T, error) {
	var parsedValue any
	var err error

	var empty T

	if value == "null" {
		switch any(empty).(type) {
		// Pointer flags preserve unset, null, and concrete values.
		case *string, *int64, *float64, *bool, *DateValue, *DateTimeValue, *TimeValue:
			return empty, nil
		// Nil maps already marshal as JSON null.
		case map[string]any:
			return empty, nil
		}
	}

	switch any(empty).(type) {
	case string:
		parsedValue = value
	case int64:
		parsedValue, err = strconv.ParseInt(value, 0, 64)
	case float64:
		parsedValue, err = strconv.ParseFloat(value, 64)
	case bool:
		parsedValue, err = strconv.ParseBool(value)
	case DateTimeValue:
		var dt DateTimeValue
		err = (&dt).Parse(value)
		if err == nil {
			parsedValue = dt
		}

	case DateValue:
		var d DateValue
		err = (&d).Parse(value)
		if err == nil {
			parsedValue = d
		}

	case TimeValue:
		var t TimeValue
		err = (&t).Parse(value)
		if err == nil {
			parsedValue = t
		}

	// Non-null pointer flags parse their pointee value.
	case *string:
		v := value
		parsedValue = &v
	case *int64:
		var v int64
		v, err = strconv.ParseInt(value, 0, 64)
		if err == nil {
			parsedValue = &v
		}
	case *float64:
		var v float64
		v, err = strconv.ParseFloat(value, 64)
		if err == nil {
			parsedValue = &v
		}
	case *bool:
		var v bool
		v, err = strconv.ParseBool(value)
		if err == nil {
			parsedValue = &v
		}
	case *DateTimeValue:
		var dt DateTimeValue
		err = (&dt).Parse(value)
		if err == nil {
			parsedValue = &dt
		}
	case *DateValue:
		var d DateValue
		err = (&d).Parse(value)
		if err == nil {
			parsedValue = &d
		}
	case *TimeValue:
		var t TimeValue
		err = (&t).Parse(value)
		if err == nil {
			parsedValue = &t
		}

	default:
		if strings.HasPrefix(value, "@") {
			// Preserve @file references for later expansion.
			parsedValue = value
		} else {
			var yamlValue T
			err = yaml.Unmarshal([]byte(value), &yamlValue)
			if err == nil {
				parsedValue = yamlValue
			} else if allowAsLiteralString(value) {
				parsedValue = value
			} else {
				parsedValue = nil
				err = fmt.Errorf("failed to parse as YAML: %w", err)
			}
		}
	}

	// Preserve YAML null through the type assertion.
	if parsedValue == nil {
		parsedValue = (*struct{})(nil)
	}

	if err == nil {
		if typedValue, ok := parsedValue.(T); ok {
			return typedValue, nil
		} else {
			expectedType := reflect.TypeFor[T]()
			err = fmt.Errorf("cannot convert %q (%v) to %v", value, parsedValue, expectedType)
		}
	}
	return empty, err

}

// Ptr returns a pointer for nullable flag defaults.
func Ptr[T any](v T) *T {
	return &v
}

// allowAsLiteralString accepts identifier-like values after YAML parsing fails.
// It rejects punctuation that may indicate malformed YAML.
func allowAsLiteralString(s string) bool {
	for _, c := range s {
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) &&
			c != '_' && c != '-' && c != '.' && c != '=' {
			return false
		}
	}
	return true
}

// Set parses a CLI value.
func (c *cliValue[T]) Set(value string) error {
	valueType := reflect.TypeOf(c.value)
	// Repeated flags append to slices.
	if valueType != nil && valueType.Kind() == reflect.Slice {
		elemType := valueType.Elem()

		var singleElem any
		var err error
		switch elemType.Kind() {
		case reflect.String:
			singleElem, err = parseCLIArg[string](value)
		case reflect.Int64:
			singleElem, err = parseCLIArg[int64](value)
		case reflect.Float64:
			singleElem, err = parseCLIArg[float64](value)
		case reflect.Bool:
			singleElem, err = parseCLIArg[bool](value)
		default:
			switch elemType.Name() {
			case "DateTimeValue":
				singleElem, err = parseCLIArg[DateTimeValue](value)
			case "DateValue":
				singleElem, err = parseCLIArg[DateValue](value)
			case "TimeValue":
				singleElem, err = parseCLIArg[TimeValue](value)
			default:
				if elemType.Kind() == reflect.Map && elemType.Key().Kind() == reflect.String {
					singleElem, err = parseCLIArg[map[string]any](value)
				} else {
					singleElem, err = parseCLIArg[any](value)
				}
			}
		}

		if err != nil {
			return err
		}

		sliceValue := reflect.ValueOf(c.value)
		if !sliceValue.IsValid() || sliceValue.IsNil() {
			sliceValue = reflect.MakeSlice(valueType, 0, 1)
		}

		newElem := reflect.ValueOf(singleElem)
		sliceValue = reflect.Append(sliceValue, newElem)

		c.value = sliceValue.Interface().(T)
	} else {
		if parsedValue, err := parseCLIArg[T](value); err != nil {
			return err
		} else {
			c.value = parsedValue
		}
	}

	return nil
}

func (c *cliValue[T]) Get() any {
	return c.value
}

func (c *cliValue[T]) String() string {
	switch v := any(c.value).(type) {
	case string, int, int64, float64, bool, DateTimeValue, DateValue, TimeValue,
		[]string, []int, []int64, []float64, []bool, []DateTimeValue, []DateValue, []TimeValue:
		return fmt.Sprintf("%v", v)

	case *string, *int64, *float64, *bool, *DateTimeValue, *DateValue, *TimeValue:
		// Render nil pointers as the CLI null literal.
		rv := reflect.ValueOf(v)
		if rv.IsNil() {
			return "null"
		}
		return fmt.Sprintf("%v", rv.Elem().Interface())

	default:
		yamlBytes, err := yaml.MarshalWithOptions(c.value, yaml.Flow(true))
		if err != nil {
			return fmt.Sprintf("%v", c.value)
		}
		return string(yamlBytes)
	}
}

func (c *cliValue[T]) IsBoolFlag() bool {
	_, ok := any(c.value).(bool)
	return ok
}

type DateValue string
type DateTimeValue string
type TimeValue string

func (d DateValue) String() string {
	return string(d)
}

func (d DateTimeValue) String() string {
	return string(d)
}

func (t TimeValue) String() string {
	return string(t)
}

// parseTimeWithFormats returns the first successful parse.
func parseTimeWithFormats(s string, formats []string) (time.Time, error) {
	var lastErr error
	for _, format := range formats {
		t, err := time.Parse(format, s)
		if err == nil {
			return t, nil
		}
		lastErr = err
	}
	return time.Time{}, lastErr
}

func (d *DateValue) Parse(s string) error {
	formats := []string{
		"2006-01-02",
		"01/02/2006",
		"Jan 2, 2006",
		"January 2, 2006",
		"2-Jan-2006",
	}

	t, err := parseTimeWithFormats(s, formats)
	if err != nil {
		return fmt.Errorf("unable to parse date: %v", err)
	}

	*d = DateValue(t.Format("2006-01-02"))
	return nil
}

func (d *DateTimeValue) Parse(s string) error {
	formats := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		time.RFC1123,
		time.RFC822,
		time.ANSIC,
	}

	t, err := parseTimeWithFormats(s, formats)
	if err != nil {
		return fmt.Errorf("unable to parse datetime: %v", err)
	}

	*d = DateTimeValue(t.Format(time.RFC3339))
	return nil
}

func (t *TimeValue) Parse(s string) error {
	formats := []string{
		"15:04:05",
		"15:04:05.999999999Z07:00",
		"3:04:05PM",
		"3:04 PM",
		"15:04",
		time.Kitchen,
	}

	parsedTime, err := parseTimeWithFormats(s, formats)
	if err != nil {
		return fmt.Errorf("unable to parse time: %v", err)
	}

	*t = TimeValue(parsedTime.Format("15:04:05"))
	return nil
}

// SettableInnerField accepts nested flag values.
type SettableInnerField interface {
	SetInnerField(string, any)
}

// InnerFieldSeeder initializes untyped nullable containers before assignment.
type InnerFieldSeeder interface {
	SeedInnerCollection(isArrayOfObjects bool)
}

func (f *Flag[T]) SetInnerField(field string, val any) {
	if f.value == nil {
		f.value = &cliValue[T]{}
	}

	if settableInnerField, ok := f.value.(SettableInnerField); ok {
		settableInnerField.SetInnerField(field, val)
		f.hasBeenSet = true
	} else {
		panic(fmt.Sprintf("cannot set an inner field on %v", f.value))
	}
}

// SeedInnerCollection seeds an untyped flag with a map or list.
func (f *Flag[T]) SeedInnerCollection(isArrayOfObjects bool) {
	if f.value == nil {
		f.value = &cliValue[T]{}
	}
	cv, ok := f.value.(*cliValue[T])
	if !ok {
		return
	}
	if reflect.ValueOf(cv.value).Kind() != reflect.Invalid {
		return
	}
	if isArrayOfObjects {
		if seed, ok := any([]map[string]any{}).(T); ok {
			cv.value = seed
		}
		return
	}
	if seed, ok := any(map[string]any{}).(T); ok {
		cv.value = seed
	}
}

func (c *cliValue[T]) SetInnerField(field string, val any) {
	flagVal := c.value
	flagValReflect := reflect.ValueOf(flagVal)
	switch flagValReflect.Kind() {
	case reflect.Slice:
		if flagValReflect.Type().Elem().Kind() != reflect.Map {
			return
		}

		sliceLen := flagValReflect.Len()
		if sliceLen > 0 {
			lastElement := flagValReflect.Index(sliceLen - 1).Interface().(map[string]any)
			if _, hasInnerField := lastElement[field]; !hasInnerField {
				lastElement[field] = val
				return
			}
		}

		newMap := map[string]any{field: val}
		switch sliceVal := any(c.value).(type) {
		case []map[string]any:
			c.value = any(append(sliceVal, newMap)).(T)
		case []any:
			c.value = any(append(sliceVal, newMap)).(T)
		}

	case reflect.Map:
		mapVal, ok := any(flagVal).(map[string]any)
		if !ok || mapVal == nil {
			mapVal = map[string]any{field: val}
			c.value = any(mapVal).(T)
		} else {
			mapVal[field] = val
		}
	}
}
