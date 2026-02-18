package format

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/philrox/ris-cli/internal/model"
)

// JSON writes search results as pretty-printed JSON to the writer.
func JSON(w io.Writer, result model.SearchResult) error {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	_, err = fmt.Fprintln(w, string(data))
	return err
}

// JSONDocument writes a document with its text content as pretty-printed JSON.
func JSONDocument(w io.Writer, doc model.Document, content string) error {
	output := model.DocumentContent{
		Metadata: doc,
		Content:  content,
	}
	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	_, err = fmt.Fprintln(w, string(data))
	return err
}
