package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/davi17g/logging-service/records"
	"github.com/rcrowley/go-metrics"
	log "github.com/sirupsen/logrus"
)

//=============================================================================
type httpServer struct {
	m            metrics.Meter
	getRequestCH chan *records.Request
	addr         string
}

//=============================================================================
func (server *httpServer) Start() error {
	http.HandleFunc("/click", server.click)
	http.HandleFunc("/impression", server.impression)
	http.HandleFunc("/completion", server.completion)
	if err := http.ListenAndServe(server.addr, nil); err != nil {
		return err
	}
	return nil
}

//=============================================================================
func (server *httpServer) impression(
	w http.ResponseWriter, req *http.Request) {

	if req.Method == http.MethodPost {
		server.m.Mark(1)
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Errorf("Unable to read request content: %s", err)
			w.WriteHeader(http.StatusBadRequest)
		} else {
			server.getRequestCH <- &records.Request{
				Handler: records.ImpressionHandler,
				Body:    body,
			}
		}

		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

//=============================================================================
func (server *httpServer) click(w http.ResponseWriter, req *http.Request) {

	if req.Method == http.MethodPost {
		server.m.Mark(1)
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Errorf("Unable to read request content: %s", err)
			w.WriteHeader(http.StatusBadRequest)
		} else {
			server.getRequestCH <- &records.Request{
				Handler: records.ClickHandler,
				Body:    body,
			}
		}
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

//=============================================================================
func (server *httpServer) completion(
	w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		server.m.Mark(1)
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Errorf("Unable to read request content: %s", err)
			w.WriteHeader(http.StatusBadRequest)
		} else {
			server.getRequestCH <- &records.Request{
				Handler: records.CompletionHandler,
				Body:    body,
			}
			w.WriteHeader(http.StatusOK)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

//=============================================================================
func (server *httpServer) Close() {
	server.m.Stop()
}

//=============================================================================
func GetNewHttpServer(
	addr string, port int, getRequestCH chan *records.Request) *httpServer {
	go metrics.Log(metrics.DefaultRegistry, 5*time.Second, log.New())
	return &httpServer{
		m:            metrics.GetOrRegisterMeter("requests", nil),
		addr:         fmt.Sprintf("%s:%d", addr, port),
		getRequestCH: getRequestCH}
}
