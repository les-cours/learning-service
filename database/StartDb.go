package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/les-cours/learning-service/env"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"

	_ "github.com/lib/pq"
)

type MongoClient struct {
	MongoDB *mongo.Client
}

func New(mongo *mongo.Client) *MongoClient {
	return &MongoClient{
		MongoDB: mongo,
	}
}

func StartDatabase() (*sql.DB, error) {
	dataSourceName := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=%s",
		env.Settings.Database.PSQLConfig.Host,
		env.Settings.Database.PSQLConfig.Port,
		env.Settings.Database.PSQLConfig.Username,
		env.Settings.Database.PSQLConfig.Password,
		env.Settings.Database.PSQLConfig.DbName,
		env.Settings.Database.PSQLConfig.SslMode,
	)
	log.Println("dataSourceName: ", dataSourceName)

	db, err := sql.Open("postgres", dataSourceName)

	if err != nil {
		log.Fatalf("Failed to connect to postgres database: %v", err)
	}
	fmt.Println("Connected to postgres!")

	db.SetConnMaxLifetime(0)
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(7)

	return db, nil
}

func StartMongoDB() (*mongo.Client, error) {
	var mongoURI = env.Settings.Database.MongoConfig.URI

	log.Println("mongoURI: ", mongoURI)
	clientOptions := options.Client().ApplyURI(mongoURI)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Printf("Error while connecting to db")
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Printf("Error while connecting to db 2 ")
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	return client, nil
}
