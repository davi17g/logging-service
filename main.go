package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/davi17g/logging-service/records"
	"github.com/davi17g/logging-service/server"
	log "github.com/sirupsen/logrus"
)

//=============================================================================
const (
	poolSize = 5
)

//=============================================================================
func main() {

	srvAddr := flag.String("srvAddr", "localhost", "Specify service address")
	srvPort := flag.Int("srvPort", 8080, "Specify service port")
	dbAddr := flag.String("dbAddr", "localhost", "Specify data-base address")
	dbPort := flag.Int("dbPort", 27017, "Specify data-base port")
	flag.Parse()

	log.Infof(
		"Starting the service with following params: "+
			"server-address: %s  "+
			"server-port: %d "+
			"database-address: %s "+
			"database-port: %d", *srvAddr, *srvPort, *dbAddr, *dbPort)

	sig := make(chan os.Signal, 1)
	setObjectCH := make(chan interface{}, 100)
	getRequestCH := make(chan *records.Request, 100)
	shutdown := make(chan struct{})

	dbWriter, err := getNewDataBaseWriter(*dbAddr, *dbPort)
	if err != nil {
		log.Panic(err)
	}
	defer func() {
		if err := dbWriter.Close(); err != nil {
			log.Errorf("Unable to close dbWriter: %s", err)
		}
	}()

	pool := getNewWorkerPool(getRequestCH, setObjectCH)
	go dbWriter.writeToDB(setObjectCH, shutdown)

	for i := 0; i < poolSize; i++ {
		go pool.doWork(shutdown)
	}

	srv := server.GetNewHttpServer(*srvAddr, *srvPort, getRequestCH)
	if err := srv.Start(); err != nil {
		log.Panicf("Server was unable to start: %s", err)
	}

	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	close(shutdown)
	<-sig
	close(sig)
}
