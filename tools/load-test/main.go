package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/davi17g/logging-service/records"
	log "github.com/sirupsen/logrus"
)

//=============================================================================
func sendRequest(
	recordType records.HandlerType, record records.Recorder) error {

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(record)
	if err != nil {
		return err
	}
	var url string
	switch recordType {

	case records.CompletionHandler:
		url = "http://localhost:8080/completion"
	case records.ClickHandler:
		url = "http://localhost:8080/click"
	case records.ImpressionHandler:
		url = "http://localhost:8080/impression"
	default:
		return errors.New("unknown handler")
	}
	_, err = http.Post(url, "application/json; charset=utf-8", b)
	if err != nil {
		return err
	}
	return nil
}

//=============================================================================
func doWorkImpression(shutdown chan struct{}) {

	for {
		select {
		case <-shutdown:
			return
		default:
			now := time.Now()
			n := rand.Int63n(9)
			obj := &records.Impression{
				DateTime:      now.String(),
				TransactionID: strconv.Itoa(int(time.Now().Unix())),
				Adtype:        records.AdType(n),
				UserID:        strconv.Itoa(int(time.Now().UnixNano())),
			}
			if err := sendRequest(records.ImpressionHandler, obj); err != nil {
				log.Errorf("Got an error while sending request: %s", err)
			}
		}
	}
}

//=============================================================================
func doWorkClick(shutdown chan struct{}) {
	for {
		select {
		case <-shutdown:
			return
		default:
			now := time.Now()
			n := rand.Int63n(9)
			obj := &records.Click{
				DateTime:      now.String(),
				TransactionID: strconv.Itoa(int(time.Now().Unix())),
				Adtype:        records.AdType(n),
				TimeToClick:   time.Now().String(),
				UserId:        strconv.Itoa(int(time.Now().UnixNano())),
			}
			if err := sendRequest(records.ClickHandler, obj); err != nil {
				log.Errorf("Got an error while sending request: %s", err)
			}
		}
	}
}

//=============================================================================
func doWorkCompletion(shutdown chan struct{}) {
	for {
		select {
		case <-shutdown:
			return
		default:
			obj := &records.Completion{
				DateTime:      time.Now().String(),
				TransactionID: strconv.Itoa(int(time.Now().Unix())),
			}
			if err := sendRequest(
				records.CompletionHandler, obj); err != nil {
				log.Errorf(
					"Got an error while sending request: %s", err)
			}
		}
	}
}

//=============================================================================
func main() {
	duration := flag.Int("time", 2, "Specify load-test duration in minutes")
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
	shutdown := make(chan struct{})
	go doWorkImpression(shutdown)
	go doWorkClick(shutdown)
	go doWorkCompletion(shutdown)
	time.Sleep(time.Duration(*duration) * time.Minute)
	close(shutdown)
}
