package config

import (
	"context"
	"log"

	//	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Client :
var Client *mongo.Client

const (
	local = "mongodb://localhost:27017"
)

func init() {
	// get a mongo sessions
	//mongoURL := os.Getenv("mongo_uri")
	c, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(local))
	if err != nil {
		log.Fatal("Error in connecting DB ", err)
	}
	if err := c.Ping(context.TODO(), nil); err != nil {
		log.Fatal("Error in pinging DB", err)
	}
	log.Println("connected to DB:::::::")
	Client = c
	return

}
