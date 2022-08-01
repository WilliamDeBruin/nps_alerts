package common

import (
	"encoding/json"
	"fmt"
)

// ToJSON returns a the JSON form of obj. If unable to Marshal obj, a JSON error message is
// returned with the %#v formatted string of the object
func ToJSON(obj interface{}) string {
	b, err := json.Marshal(obj)
	if err != nil {
		return fmt.Sprintf(`{"error":"failed to marshal into JSON","obj":%q}`, fmt.Sprintf("%#v", obj))
	}
	return string(b)
}
