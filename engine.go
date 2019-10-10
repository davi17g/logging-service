package main

import (
	"github.com/davi17g/logging-service/database"
	"github.com/davi17g/logging-service/records"
	"github.com/davi17g/logging-service/utils"
	log "github.com/sirupsen/logrus"
)

//=============================================================================
type dataBaseWriter struct {
	dataBase *database.DataBaseBroker
}

//=============================================================================
func (dbw *dataBaseWriter) writeToDB(
	setObjectCH chan interface{}, shutdown chan struct{}) {

	for {
		select {
		case <-shutdown:
			return
		case obj := <-setObjectCH:
			if err := dbw.dataBase.SetRecord("logs", obj); err != nil {
				log.Errorf("Got an insertion error: %s", err)
			}
		}
	}
}

//=============================================================================
func (dbw *dataBaseWriter) Close() error {
	if err := dbw.dataBase.Close(); err != nil {
		return err
	}
	return nil
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
			obj, err := utils.ConvJsonToObject(request)
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
	db, err := database.GetNewDataBaseBroker(addr, port)
	if err != nil {
		return nil, err
	}
	return &dataBaseWriter{dataBase: db}, nil
}
