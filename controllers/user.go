package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/MoazamAli/go-mongo-crud.git/models"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserController struct {
	session *mongo.Client
}

func isValidObjectId(id string) bool {
	_, err := primitive.ObjectIDFromHex(id)
	return err == nil
}
func NewUserController(s *mongo.Client) *UserController {
	return &UserController{s}
}

func (uc *UserController) GetUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := uc.session.Database("go-mongo-crud").Collection("users").Find(ctx, bson.M{})
	if err != nil {
		log.Fatal("Find failed:", err)
	}
	var users []models.User
	if err = cursor.All(ctx, &users); err != nil {
		log.Fatal(err)
	}
	uj, err := json.Marshal(users)
	if err != nil {
		fmt.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(uj)
}

func (uc *UserController) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	if !isValidObjectId(id) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	oid, _ := primitive.ObjectIDFromHex(id)
	u := models.User{}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// var result bson.M
	err := uc.session.Database("go-mongo-crud").Collection("users").FindOne(ctx, bson.M{"_id": oid}).Decode(&u)
	if err != nil {
		log.Fatal("Find failed:", err)
	}
	// fmt.Println("Found document:", result)

	// if err := uc.session.Database("go-mongo-crud").Collection("users").FindOne(ctx,).One(&u); err != nil {
	// 	w.WriteHeader(404)
	// 	return
	// }
	uj, err := json.Marshal(u)
	if err != nil {
		fmt.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(uj)
}

func (uc *UserController) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	u := models.User{}
	json.NewDecoder(r.Body).Decode(&u)
	u.Id = primitive.NewObjectID()
	uc.session.Database("go-mongo-crud").Collection("users").InsertOne(context.TODO(), u)
	uj, err := json.Marshal(u)
	if err != nil {
		fmt.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(uj)
}

func (uc *UserController) DeleteUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	if !isValidObjectId(id) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	oid, _ := primitive.ObjectIDFromHex(id)
	if _, err := uc.session.Database("go-mongo-crud").Collection("users").DeleteOne(context.TODO(), bson.M{"_id": oid}); err != nil {
		w.WriteHeader(404)
		return
	}
	w.WriteHeader(http.StatusOK)
}
