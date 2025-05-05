package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/MoazamAli/go-mongo-crud.git/controllers"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	r := httprouter.New()
	uc := controllers.NewUserController(getSession())
	r.GET("/user/:id", uc.GetUser)
	r.GET("/users", uc.GetUsers)
	r.POST("/users", uc.CreateUser)
	r.DELETE("/user/:id", uc.DeleteUser)
	http.ListenAndServe(":8080", r)
}

func getSession() *mongo.Client {

	// Set client options
	clientOpts := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatal("Connection error:", err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("MongoDB not reachable:", err)
	}
	fmt.Println("Connected to MongoDB!")
	return client

	// s, err := mgo.Dial("mongodb://127.0.0.1:27017")
	// if err != nil {
	// 	panic(err)
	// }
	// return s
}
