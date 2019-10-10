package database

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
type DataBaseBroker struct {
	client *mongo.Client
}

//=============================================================================
func (dbb *DataBaseBroker) SetRecord(
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
func GetNewDataBaseBroker(addr string, port int) (*DataBaseBroker, error) {
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

	return &DataBaseBroker{client: client}, nil
}

//=============================================================================
func (dbb *DataBaseBroker) Close() error {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	if err := dbb.client.Disconnect(ctx); err != nil {
		return err
	}
	return nil
}

