package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

var rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Validator struct {
	GeneralErrors    []string
	ValidationErrors map[string][]string
}

// Valid checks whether there are any validation errors present
func (v *Validator) Valid() bool {
	return len(v.ValidationErrors) == 0 && len(v.GeneralErrors) == 0
}

// AddValidationError adds an error message to the FieldErrors map (so long as no
// entry already exists for the given key).
func (v *Validator) AddValidationError(key, message string) {
	if v.ValidationErrors == nil {
		v.ValidationErrors = make(map[string][]string)
	}
	if _, exists := v.ValidationErrors[key]; exists {
		v.ValidationErrors[key] = append(v.ValidationErrors[key], message)
	} else {
		v.ValidationErrors[key] = []string{message}
	}
}

// AddGeneralError helper for adding error messages to the new
// GeneralErrors slice.
func (v *Validator) AddGeneralError(message string) {
	v.GeneralErrors = append(v.GeneralErrors, message)
}

// CheckField adds an error message to the ValidationErrors map only if a
// validation check is not 'ok'.
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddValidationError(key, message)
	}
}

// NotBlank returns true if a value is not an empty string
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// MaxChars returns true if a value contains no more than n characters.
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// MinChars returns true if a value contains at least n characters.
func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

// IsEmailAddress returns true if a value matches a provided compiled regular
// expression pattern.
func IsEmailAddress(value string) bool {
	return rxEmail.MatchString(value)
}

// PermittedValue returns true if a value is in a list of permitted integers.
func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}
	return false
}
