package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

// Character struct to represent the data fetched from Rick and Morty API
type Character struct {
	ID      int    `json:"id" bson:"_id"`
	Name    string `json:"name" bson:"name"`
	Status  string `json:"status" bson:"status"`
	Species string `json:"species" bson:"species"`
	Type    string `json:"type" bson:"type"`
	Gender  string `json:"gender" bson:"gender"`
	Origin  struct {
		Name string `json:"name" bson:"name"`
	} `json:"origin" bson:"origin"`
	Location struct {
		Name string `json:"name" bson:"name"`
	} `json:"location" bson:"location"`
	Image string `json:"image" bson:"image"`
}

// connectToMongoDB connects to MongoDB database
func connectToMongoDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Ping the database to ensure the connection is successful
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	mongoClient = client
	log.Println("Successfully connected to MongoDB")
}

// fetchAndInsertCharacters fetches data from the Rick and Morty API and inserts it into MongoDB
func fetchAndInsertCharacters() {
	// Fetch data from the API
	resp, err := http.Get("https://rickandmortyapi.com/api/character")
	if err != nil {
		log.Fatalf("Failed to fetch data from Rick and Morty API: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		Info struct {
			Pages int `json:"pages"`
		} `json:"info"`
		Results []Character `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatalf("Failed to decode API response: %v", err)
	}

	// Insert the characters into MongoDB, ensuring no duplicates are inserted
	collection := mongoClient.Database("rickandmorty").Collection("characters")
	var insertModels []mongo.WriteModel

	for _, character := range result.Results {
		// Use an upsert to insert if not already present, or update if the character exists
		upsertModel := mongo.NewUpdateOneModel()
		upsertModel.Filter = bson.M{"_id": character.ID} // Ensure that the _id is unique
		upsertModel.Update = bson.M{"$set": character}
		// Initialize Upsert pointer if it's nil
		if upsertModel.Upsert == nil {
			upsertModel.Upsert = new(bool)
		}
		*upsertModel.Upsert = true // Set the value to true
		insertModels = append(insertModels, upsertModel)
	}

	// Perform the bulk write (insert/update) with upsert
	_, err = collection.BulkWrite(context.Background(), insertModels)
	if err != nil {
		log.Fatalf("Failed to insert/update data in MongoDB: %v", err)
		return
	}

	log.Printf("Successfully inserted %d characters into MongoDB", len(result.Results))
}

// searchCharacters handles the search request to search for characters by name
func searchCharacters(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Missing 'name' query parameter", http.StatusBadRequest)
		return
	}

	collection := mongoClient.Database("rickandmorty").Collection("characters")
	filter := bson.M{"name": bson.M{"$regex": name, "$options": "i"}}

	var results []bson.M
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &results); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func main() {
	connectToMongoDB()

	// Fetch and insert Rick and Morty characters from the API
	fetchAndInsertCharacters()

	// Set up the router and define routes
	r := mux.NewRouter()
	r.HandleFunc("/search", searchCharacters).Methods("GET")

	// Enable CORS middleware for the router
	http.ListenAndServe(":8080", handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000"}),        // Allow the React frontend origin
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}), // Allow specific HTTP methods
	)(r))
}
