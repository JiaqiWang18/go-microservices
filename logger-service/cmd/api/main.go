package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"logger-service/data"
	"net/http"
	"os"
	"time"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017"
	gRpcPort = "50001"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	// connect to mongo
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}
	log.Println("Connected to MongoDB")
	client = mongoClient
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	// create a new instance of the config
	app := Config{
		Models: data.New(client),
	}
	// start the web server
	// go app.serve()
	srv := &http.Server{
		Addr:    ":" + webPort,
		Handler: app.routes(),
	}
	err = srv.ListenAndServe()
	log.Println("Starting logger service")
	if err != nil {
		log.Panic(err)
	}
}

//func (app *Config) serve() {
//	// start the web server
//	srv := &http.Server{
//		Addr:    ":" + webPort,
//		Handler: app.routes(),
//	}
//	err := srv.ListenAndServe()
//	if err != nil {
//		log.Panic(err)
//	}
//}

func connectToMongo() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: os.Getenv("MONGO_INITDB_ROOT_USERNAME"),
		Password: os.Getenv("MONGO_INITDB_ROOT_PASSWORD"),
	})
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Panic(err)
		return nil, err
	}
	return c, nil
}
