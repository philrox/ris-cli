package parser

import (
	"encoding/json"
	"fmt"
)

// FlexibleString handles JSON fields that can be either a plain string
// or an object with a "#text" property.
type FlexibleString string

func (f *FlexibleString) UnmarshalJSON(data []byte) error {
	// Try plain string first.
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		*f = FlexibleString(s)
		return nil
	}

	// Try object with #text.
	var obj struct {
		Text string `json:"#text"`
	}
	if err := json.Unmarshal(data, &obj); err == nil {
		*f = FlexibleString(obj.Text)
		return nil
	}

	return fmt.Errorf("cannot unmarshal %s into FlexibleString", string(data))
}

func (f FlexibleString) String() string {
	return string(f)
}

// FlexibleArray handles JSON fields that can be either a single object
// or an array of objects of type T.
type FlexibleArray[T any] []T

func (f *FlexibleArray[T]) UnmarshalJSON(data []byte) error {
	// Try array first.
	var arr []T
	if err := json.Unmarshal(data, &arr); err == nil {
		*f = arr
		return nil
	}

	// Try single object.
	var single T
	if err := json.Unmarshal(data, &single); err == nil {
		*f = []T{single}
		return nil
	}

	return fmt.Errorf("cannot unmarshal %s into FlexibleArray", string(data))
}

// FlexibleInt handles JSON fields that can be a number, a string containing
// a number, or an object with "#text" containing a number.
type FlexibleInt int

func (f *FlexibleInt) UnmarshalJSON(data []byte) error {
	// Try plain number.
	var n int
	if err := json.Unmarshal(data, &n); err == nil {
		*f = FlexibleInt(n)
		return nil
	}

	// Try string.
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		var num int
		if _, err := fmt.Sscanf(s, "%d", &num); err == nil {
			*f = FlexibleInt(num)
			return nil
		}
	}

	// Try object with #text.
	var obj struct {
		Text string `json:"#text"`
	}
	if err := json.Unmarshal(data, &obj); err == nil {
		var num int
		if _, err := fmt.Sscanf(obj.Text, "%d", &num); err == nil {
			*f = FlexibleInt(num)
			return nil
		}
	}

	return fmt.Errorf("cannot unmarshal %s into FlexibleInt", string(data))
}
