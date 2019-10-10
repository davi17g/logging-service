package utils

import (
	"encoding/json"

	"github.com/davi17g/logging-service/records"
)

//=============================================================================
func ConvJsonToObject(req *records.Request) (records.Recorder, error) {

	var rc records.Recorder
	switch req.Handler {
	case records.ImpressionHandler:
		rc = &records.Impression{}
	case records.ClickHandler:
		rc = &records.Click{}
	case records.CompletionHandler:
		rc = &records.Completion{}
	}
	if err := json.Unmarshal(req.Body, rc); err != nil {
		return nil, err
	}
	return rc, nil
}