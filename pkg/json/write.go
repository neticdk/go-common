// Deprecated: json is deprecated - use github.com/neticdk/go-stdlib/xjson
package json

import (
	"io"

	"github.com/neticdk/go-stdlib/xjson"
)

// PrettyPrintJSON pretty prints JSON
// Deprecated: json is deprecated - use github.com/neticdk/go-stdlib/xjson.PrettyPrintJSON
func PrettyPrintJSON(body []byte, writer io.Writer) error {
	return xjson.PrettyPrintJSON(body, writer)
}
