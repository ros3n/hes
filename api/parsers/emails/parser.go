package emails

import (
	"encoding/json"
	"errors"
	"github.com/ros3n/hes/api/models"
	"io"
	"log"
	"mime"
)

var (
	// ErrUnsupportedContentType indicates that the content type of the request body is not supported.
	ErrUnsupportedContentType = errors.New("unsupported content type")

	// ErrParsingFailed indicates that the request body is malformed
	ErrParsingFailed = errors.New("malformed request body")
)

// PayloadParser is an interface that declares methods needed to convert a request body to an EmailChangeSet.
type PayloadParser interface {
	Parse(io.Reader) error
	Data() *models.EmailChangeSet
}

// NewPayloadParses is a factory method that returns a PayloadParser instance matching a given content type.
func NewPayloadParser(contentType string) (PayloadParser, error) {
	parsedContentType := parseContentType(contentType)
	switch parsedContentType {
	case "application/json":
		return &JsonPayloadParser{}, nil
	}
	return nil, ErrUnsupportedContentType
}

func parseContentType(contentType string) string {
	parsed, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		log.Println(err)
		return ""
	}
	return parsed
}

type JsonPayloadParser struct {
	changeSet *models.EmailChangeSet
}

func (p *JsonPayloadParser) Parse(payload io.Reader) error {
	changeSet := models.EmailChangeSet{}
	decoder := json.NewDecoder(payload)

	err := decoder.Decode(&changeSet)
	if err != nil {
		log.Println(err)
		return ErrParsingFailed
	}

	p.changeSet = &changeSet

	return nil
}

func (p *JsonPayloadParser) Data() *models.EmailChangeSet {
	return p.changeSet
}
