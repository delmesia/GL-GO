package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Define an error that the UnmarshalJSON() method can return if we're
// unable to parse or convert the JSON string successfully.
var ErrInvalidRuntimeFormat = errors.New("invalid runtime format")

// Declare a custom Runtime type, which has the type int32
// same as the Movie struct field.
type Runtime int32

// Implement a MarshalJSON() method on the Runtime type so that
// it satisfies the json.Marshaler() interface. This should return
// the json encoded value for the movie runtime.
// (in my case, it will return a string in the format "<runtimes> mins").
func (r Runtime) MarshalJSON() ([]byte, error) {
	// Generate a string containing the movie runtime in the required format.
	jsonValue := fmt.Sprintf("%d mins", r)
	// Use the strconv.Quote() function to wrap it in double quotes.
	// It needs to be surrounred by double quotes in order to be a valid *JSON string*.
	quotedJSONValue := strconv.Quote(jsonValue)
	// Convert the quoted string value to a byteslice and return it.
	return []byte(quotedJSONValue), nil
}

func (r *Runtime) UnmarshalJSON(jsonValue []byte) error {

	// We expect that the incoming JSON value will be a string in the format
	// "<runtime> mins", and the first thing we need to do is remove the
	// surrounding double-quotes from string. If we cannot unquote, return
	// ErrInvalidRuntimeFormat error.
	unquotedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidRuntimeFormat
	}
	// Split the string to isolate the part containing the number.
	parts := strings.Split(unquotedJSONValue, " ")

	// Sanity check the the parts of the string to make sure it was in the expected format.
	// If it isn't, we return the ErrInvalidRuntimeFormat error again.
	if len(parts) != 2 || parts[1] != "mins" {
		return ErrInvalidRuntimeFormat
	}
	// Otherwise, parse the string containing the number into an int32. Again,
	// if this fails, return the ErrInvalidRuntimeFormat.
	i, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return ErrInvalidRuntimeFormat
	}
	// Convert the int32 to Runtime type and assign this to the receiver.
	// Note the usage of * operator to dereference the receiver (which is a pointer
	// to a Runtime type) in order to set the underlying value of the pointer
	*r = Runtime(i)

	return nil
}
