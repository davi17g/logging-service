package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/davi17g/logging-service/records"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

//=============================================================================
type dataBaseBroker struct {
	client *mongo.Client
}

//=============================================================================
func (dbb *dataBaseBroker) setRecord(
	database string, record interface{}) error {

	var collection *mongo.Collection
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	switch record.(type) {
	case *records.Impression:
		collection = dbb.client.Database(database).Collection("impression")
		input, _ := record.(*records.Impression)
		_, err := collection.InsertOne(ctx, *input.Record())
		if err != nil {
			log.Errorf("Unable to insert a document: %s", err)
		}
	case *records.Click:
		collection = dbb.client.Database(database).Collection("click")
		input, _ := record.(*records.Click)
		_, err := collection.InsertOne(ctx, *input.Record())
		if err != nil {
			log.Errorf("Unable to insert a document: %s", err)
		}
	case *records.Completion:
		collection = dbb.client.Database(database).Collection("completion")
		input, _ := record.(*records.Completion)
		_, err := collection.InsertOne(ctx, *input.Record())
		if err != nil {
			log.Errorf("Unable to insert a document: %s", err)
		}
	default:
		return errors.New("unknown object")
	}
	return nil
}

//=============================================================================
func getNewDataBaseBroker(addr string, port int) (*dataBaseBroker, error) {
	client, err := mongo.NewClient(
		options.Client().ApplyURI(
			fmt.Sprintf("mongodb://%s:%d", addr, port)))
	if err != nil {
		return nil, err
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	if err = client.Connect(ctx); err != nil {
		return nil, err
	}

	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	return &dataBaseBroker{client: client}, nil
}
