package main

import (
	"encoding/json"

	"github.com/davi17g/logging-service/records"
	log "github.com/sirupsen/logrus"
)

//=============================================================================
func convJsonToObject(req *records.Request) (records.Recorder, error) {

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

//=============================================================================
type dataBaseWriter struct {
	dataBase *dataBaseBroker
}

//=============================================================================
func (dbw *dataBaseWriter) writeToDB(
	setObjectCH chan interface{}, shutdown chan struct{}) {

	for {
		select {
		case <-shutdown:
			return
		case obj := <-setObjectCH:
			if err := dbw.dataBase.setRecord("logs", obj); err != nil {
				log.Errorf("Got an insertion error: %s", err)
			}
		}
	}
}

//=============================================================================
type workerPool struct {
	getRequestCH chan *records.Request
	setObjectCH  chan interface{}
}

//=============================================================================
func (w *workerPool) doWork(shutdown chan struct{}) {
	for {
		select {
		case <-shutdown:
			return
		case request := <-w.getRequestCH:
			obj, err := convJsonToObject(request)
			if err != nil {
				log.Errorf("Unable to unmarshal json object: %s", err)
				continue
			}
			w.setObjectCH <- obj
		}
	}
}

//=============================================================================
func getNewWorkerPool(
	getRequestCH chan *records.Request,
	setObjectCH chan interface{}) *workerPool {
	return &workerPool{getRequestCH: getRequestCH, setObjectCH: setObjectCH}
}

//=============================================================================
func getNewDataBaseWriter(addr string, port int) (*dataBaseWriter, error) {
	db, err := getNewDataBaseBroker(addr, port)
	if err != nil {
		return nil, err
	}
	return &dataBaseWriter{dataBase: db}, nil
}
