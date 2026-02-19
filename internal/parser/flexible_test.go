package parser

import (
	"encoding/json"
	"testing"
)

// ---------------------------------------------------------------------------
// FlexibleString
// ---------------------------------------------------------------------------

func TestFlexibleStringUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		json string
		want string
	}{
		{"plain string", `"hello"`, "hello"},
		{"empty string", `""`, ""},
		{"object with #text", `{"#text":"from object"}`, "from object"},
		{"object with empty #text", `{"#text":""}`, ""},
		{"raw fallback number", `42`, "42"},
		{"raw fallback bool", `true`, "true"},
		{"null becomes empty", `null`, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fs FlexibleString
			if err := json.Unmarshal([]byte(tt.json), &fs); err != nil {
				t.Fatalf("UnmarshalJSON(%s) returned error: %v", tt.json, err)
			}
			if string(fs) != tt.want {
				t.Errorf("UnmarshalJSON(%s) = %q, want %q", tt.json, string(fs), tt.want)
			}
		})
	}
}

func TestFlexibleStringString(t *testing.T) {
	tests := []struct {
		name string
		fs   FlexibleString
		want string
	}{
		{"non-empty", FlexibleString("hello"), "hello"},
		{"empty", FlexibleString(""), ""},
		{"unicode", FlexibleString("\u00d6sterreich"), "\u00d6sterreich"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fs.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestFlexibleStringInStruct verifies that FlexibleString works correctly
// when embedded in a parent struct being unmarshalled from JSON.
func TestFlexibleStringInStruct(t *testing.T) {
	type Doc struct {
		Title FlexibleString `json:"title"`
	}

	tests := []struct {
		name string
		json string
		want string
	}{
		{"plain string field", `{"title":"My Title"}`, "My Title"},
		{"object field", `{"title":{"#text":"My Title"}}`, "My Title"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var doc Doc
			if err := json.Unmarshal([]byte(tt.json), &doc); err != nil {
				t.Fatalf("Unmarshal returned error: %v", err)
			}
			if string(doc.Title) != tt.want {
				t.Errorf("Title = %q, want %q", string(doc.Title), tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// FlexibleArray
// ---------------------------------------------------------------------------

func TestFlexibleArrayUnmarshalJSON(t *testing.T) {
	type Item struct {
		Name string `json:"name"`
	}

	tests := []struct {
		name      string
		json      string
		wantLen   int
		wantNames []string
	}{
		{
			"array of items",
			`[{"name":"a"},{"name":"b"},{"name":"c"}]`,
			3,
			[]string{"a", "b", "c"},
		},
		{
			"single item",
			`{"name":"only"}`,
			1,
			[]string{"only"},
		},
		{
			"array with one item",
			`[{"name":"solo"}]`,
			1,
			[]string{"solo"},
		},
		{
			"empty array",
			`[]`,
			0,
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fa FlexibleArray[Item]
			if err := json.Unmarshal([]byte(tt.json), &fa); err != nil {
				t.Fatalf("UnmarshalJSON(%s) returned error: %v", tt.json, err)
			}
			if len(fa) != tt.wantLen {
				t.Fatalf("len = %d, want %d", len(fa), tt.wantLen)
			}
			for i, wantName := range tt.wantNames {
				if fa[i].Name != wantName {
					t.Errorf("fa[%d].Name = %q, want %q", i, fa[i].Name, wantName)
				}
			}
		})
	}
}

func TestFlexibleArrayStrings(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		wantLen int
		want    []string
	}{
		{
			"array of strings",
			`["a","b","c"]`,
			3,
			[]string{"a", "b", "c"},
		},
		{
			"single string",
			`"only"`,
			1,
			[]string{"only"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fa FlexibleArray[string]
			if err := json.Unmarshal([]byte(tt.json), &fa); err != nil {
				t.Fatalf("UnmarshalJSON(%s) returned error: %v", tt.json, err)
			}
			if len(fa) != tt.wantLen {
				t.Fatalf("len = %d, want %d", len(fa), tt.wantLen)
			}
			for i, w := range tt.want {
				if fa[i] != w {
					t.Errorf("fa[%d] = %q, want %q", i, fa[i], w)
				}
			}
		})
	}
}

func TestFlexibleArrayUnmarshalError(t *testing.T) {
	var fa FlexibleArray[struct {
		Name string `json:"name"`
	}]
	// A plain number cannot be unmarshalled into a struct or a slice of structs.
	if err := json.Unmarshal([]byte(`42`), &fa); err == nil {
		t.Error("expected error for invalid input, got nil")
	}
}

// TestFlexibleArrayInStruct verifies FlexibleArray works when embedded in
// a parent struct.
func TestFlexibleArrayInStruct(t *testing.T) {
	type Tag struct {
		Value string `json:"value"`
	}
	type Doc struct {
		Tags FlexibleArray[Tag] `json:"tags"`
	}

	tests := []struct {
		name    string
		json    string
		wantLen int
	}{
		{"array field", `{"tags":[{"value":"a"},{"value":"b"}]}`, 2},
		{"single object field", `{"tags":{"value":"a"}}`, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var doc Doc
			if err := json.Unmarshal([]byte(tt.json), &doc); err != nil {
				t.Fatalf("Unmarshal returned error: %v", err)
			}
			if len(doc.Tags) != tt.wantLen {
				t.Errorf("len(Tags) = %d, want %d", len(doc.Tags), tt.wantLen)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// FlexibleInt
// ---------------------------------------------------------------------------

func TestFlexibleIntUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		json string
		want int
	}{
		{"plain number", `42`, 42},
		{"zero", `0`, 0},
		{"negative number", `-7`, -7},
		{"string number", `"123"`, 123},
		{"string negative", `"-5"`, -5},
		{"object with #text", `{"#text":"99"}`, 99},
		{"object with #text negative", `{"#text":"-3"}`, -3},
		{"invalid string fallback zero", `"not-a-number"`, 0},
		{"object with non-numeric #text", `{"#text":"abc"}`, 0},
		{"bool fallback zero", `true`, 0},
		{"null fallback zero", `null`, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fi FlexibleInt
			if err := json.Unmarshal([]byte(tt.json), &fi); err != nil {
				t.Fatalf("UnmarshalJSON(%s) returned error: %v", tt.json, err)
			}
			if int(fi) != tt.want {
				t.Errorf("UnmarshalJSON(%s) = %d, want %d", tt.json, int(fi), tt.want)
			}
		})
	}
}

// TestFlexibleIntInStruct verifies that FlexibleInt works correctly when
// embedded in a parent struct.
func TestFlexibleIntInStruct(t *testing.T) {
	type Page struct {
		Number FlexibleInt `json:"number"`
	}

	tests := []struct {
		name string
		json string
		want int
	}{
		{"number field", `{"number":5}`, 5},
		{"string field", `{"number":"10"}`, 10},
		{"object field", `{"number":{"#text":"20"}}`, 20},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var page Page
			if err := json.Unmarshal([]byte(tt.json), &page); err != nil {
				t.Fatalf("Unmarshal returned error: %v", err)
			}
			if int(page.Number) != tt.want {
				t.Errorf("Number = %d, want %d", int(page.Number), tt.want)
			}
		})
	}
}
