package config

import (
	"context"
	"log"
	"os"

	//	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Client :
var Client *mongo.Client

const (
	local = "localhost"
)

func init() {
	// get a mongo sessions
	mongoURL := os.Getenv("mongo_uri")
	c, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURL))
	if err != nil {
		log.Fatal("Error in connecting DB ", err)
	}
	if err := c.Ping(context.TODO(), nil); err != nil {
		log.Fatal("Error in pinging DB", err)
	}
	log.Println("coneected to DB:::::::")
	Client = c
	return

}
