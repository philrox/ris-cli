package parser

import (
	"github.com/philrox/ris-cli/internal/model"
)

// ParseDocumentResponse parses a search API response and extracts the first
// document, intended for single-document retrieval via the dokument command.
func ParseDocumentResponse(data []byte) (model.Document, error) {
	result, err := ParseSearchResponse(data)
	if err != nil {
		return model.Document{}, err
	}
	if len(result.Documents) == 0 {
		return model.Document{}, nil
	}
	return result.Documents[0], nil
}
