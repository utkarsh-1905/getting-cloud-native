package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testGo struct {
	Age     int    `bson:"age"`
	Message string `bson:"message"`
}

type JSONOutput struct {
	Status string `json:"status"`
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", homeRoute).Methods("GET")
	_ = godotenv.Load()
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGODB")))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	collection := client.Database("test").Collection("testgos")
	defer client.Disconnect(ctx)

	r.HandleFunc("/db/create", func(w http.ResponseWriter, r *http.Request) {
		for i := 0; i < 100; i++ {
			rand := testGo{Age: rand.Intn(100), Message: "abc" + strconv.Itoa(rand.Intn(100))}
			_, err := collection.InsertOne(ctx, rand)
			if err != nil {
				fmt.Println(err)
				_ = json.NewEncoder(w).Encode(JSONOutput{Status: "error"})
				panic(err)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(JSONOutput{Status: "OK"})
	}).Methods("GET")

	r.HandleFunc("/db/read", func(w http.ResponseWriter, r *http.Request) {
		query := mux.Vars(r)
		_, err := collection.Find(ctx, bson.M{"age": query["age"]})
		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			fmt.Println(err)
			_ = json.NewEncoder(w).Encode(JSONOutput{Status: "error"})
		} else {
			_ = json.NewEncoder(w).Encode(JSONOutput{Status: "OK"})
		}
	}).Methods("GET")

	r.HandleFunc("/db/delete", func(w http.ResponseWriter, r *http.Request) {
		_, err := collection.DeleteMany(ctx, bson.M{})
		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			fmt.Println(err)
			_ = json.NewEncoder(w).Encode(JSONOutput{Status: "error"})
		} else {
			_ = json.NewEncoder(w).Encode(JSONOutput{Status: "OK"})
		}
	}).Methods("GET")

	fmt.Printf("Server is running on port 3005")
	log.Fatal(http.ListenAndServe(":3005", r))
}

func homeRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")
}
