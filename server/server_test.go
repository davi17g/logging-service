package server

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/davi17g/logging-service/records"
	"github.com/davi17g/logging-service/utils"
	"github.com/stretchr/testify/assert"
)

//=============================================================================
func recordToJson(
	t *testing.T, recorder records.Recorder) *bytes.Buffer {
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(recorder)
	if err != nil {
		t.Fail()
	}
	return b
}

//=============================================================================
func TestServer(t *testing.T) {

	getRequest := make(chan *records.Request, 1)
	server := GetNewHttpServer("test", 0000, getRequest)

	impression := &records.Impression{
		DateTime:      "date-time",
		TransactionID: "tid",
		Adtype:        0,
		UserID:        "uid",
	}

	click := &records.Click{
		DateTime:      "date-time",
		TransactionID: "tid",
		Adtype:        0,
		TimeToClick:   "time-to-click",
		UserId:        "uid",
	}

	completion := &records.Completion{
		DateTime:      "date-time",
		TransactionID: "uid",
	}

	testcases := map[string]struct {
		record   *bytes.Buffer
		handler  func(w http.ResponseWriter, r *http.Request)
		expected records.Recorder
	}{
		"impression": {
			record: recordToJson(t, impression),
			handler: server.impression,
			expected: impression},
		"click":      {
			record: recordToJson(t, click),
			handler: server.click,
			expected: click},
		"completion": {
			record: recordToJson(t, completion),
			handler: server.completion,
			expected: completion},
	}

	for tc, val := range testcases {

		t.Run(tc, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/"+tc, val.record)
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(val.handler)
			handler.ServeHTTP(rr, req)
			rc := <-getRequest
			actual, err := utils.ConvJsonToObject(rc)
			if err != nil {
				log.Fatal(err)
			}

			assert.Equal(t, val.expected.String(), actual.String())
			assert.Equal(t, http.StatusOK, rr.Code)
		})
	}
}
