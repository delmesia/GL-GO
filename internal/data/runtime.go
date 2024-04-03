package data

import (
	"fmt"
	"strconv"
)

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
