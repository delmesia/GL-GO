package validator

import "regexp"

// Declare a regular expression for sanity checking the format of the email address.
// This pattern is taken from https://html.spec.whatwg.org/#valid-e-mail-address.
var (
	EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\. [a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

// A validator struct which contains a map of validator errors.
type Validator struct {
	Errors map[string]string
}

// A helper which creates a new validator instance with an empty errors map.
func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

// Valid() returns true if the Errors map doesn't contain any entries.
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// AddError adds an error message to map (so long as no entry
// already exist for the given key.)
func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// Check adds an error message to the map only if validation check is not 'ok'
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// Generic function which returns true if a specific value is in a list
func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}
	return false
}

// Matches return true if a string value matches a specific regexp patterns.
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// Generic function which returns true if all values in a slices are unique.
func Unique[t comparable](values []t) bool {
	uniqueValues := make(map[t]bool)

	for _, value := range values {
		uniqueValues[value] = true
	}
	return len(values) == len(uniqueValues)
}
