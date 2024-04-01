/*
Written by Efdal Sancak (aka z3ntl3)

github.com/z3ntl3

Disclaimer: Educational purposes only
License: GNU
*/
package db

import (
	"context"
	"log"
	"time"

	"github.com/spf13/viper"
	"github.com/z3ntl3/VidmolySpoof/globals"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func Connect(connURI string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	c, err := mongo.Connect(ctx, options.Client().ApplyURI(connURI).SetMinPoolSize(500).SetMaxPoolSize(500))
	if err != nil {
		log.Fatal(err)
	}

	if err := c.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	globals.MongoClient = c.Database(viper.GetStringMap("database")["name"].(string))
}
