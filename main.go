package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	helpers "github.com/kanavmittal/simple-go-service/Helpers"
	models "github.com/kanavmittal/simple-go-service/Models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var collection = helpers.ConnectDB()

func main() {

	http.HandleFunc("/users/", getUser)
	http.HandleFunc("/users", createUser)
	http.HandleFunc("/posts/", getPost)
	http.HandleFunc("/posts", createPost)
	http.HandleFunc("/posts/users/", userPosts)

	err := http.ListenAndServe(":9090", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "POST" {
		var user models.Users
		_ = json.NewDecoder(r.Body).Decode(&user)
		hash, err := HashPassword(user.Password)
		user.Password = hash
		result, err := collection.Collection("Users").InsertOne(context.TODO(), user)
		if err != nil {
			helpers.GetError(err, w)
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(result)
	} else {
		w.WriteHeader(405) // Return 405 Method Not Allowed.
		return
	}
}

func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "GET" {
		id := r.URL.Path[len("/users/"):]
		objectId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			log.Printf("invalid id, %v", err)
			w.WriteHeader(400) // Return 400 Bad Request.
			return
		}
		var users []models.Users
		cur, err := collection.Collection("Users").Find(context.TODO(), bson.M{"_id": objectId})
		if err != nil {
			helpers.GetError(err, w)
			return
		}
		defer cur.Close(context.TODO())
		for cur.Next(context.TODO()) {
			var user models.Users
			err := cur.Decode(&user)
			if err != nil {
				log.Fatal(err)
			}
			users = append(users, user)
		}
		if err := cur.Err(); err != nil {
			log.Fatal(err)
		}
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(users)
	} else {
		w.WriteHeader(405) // Return 405 Method Not Allowed.
		return
	}
}

func createPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "POST" {
		var post models.Posts
		_ = json.NewDecoder(r.Body).Decode(&post)
		if post.Timestamp.IsZero() {
			post.Timestamp = time.Now()
		}
		result, err := collection.Collection("Posts").InsertOne(context.TODO(), post)
		if err != nil {
			helpers.GetError(err, w)
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(result)
	} else {
		w.WriteHeader(405) // Return 405 Method Not Allowed.
		return
	}
}

func getPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "GET" {
		id := r.URL.Path[len("/posts/"):]
		objectId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			log.Printf("invalid id, %v", err)
			w.WriteHeader(400) // Return 400 Bad Request.
			return
		}
		var posts []models.Posts
		cur, err := collection.Collection("Posts").Find(context.TODO(), bson.M{"_id": objectId})
		if err != nil {
			helpers.GetError(err, w)
			return
		}
		defer cur.Close(context.TODO())
		for cur.Next(context.TODO()) {
			var post models.Posts
			err := cur.Decode(&post)
			if err != nil {
				log.Fatal(err)
			}
			posts = append(posts, post)
		}
		if err := cur.Err(); err != nil {
			log.Fatal(err)
		}
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(posts)
	} else {
		w.WriteHeader(405) // Return 405 Method Not Allowed.
		return
	}
}

func userPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "GET" {
		id := r.URL.Path[len("/posts/users/"):]
		objectId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			log.Printf("invalid id, %v", err)
			w.WriteHeader(400) // Return 400 Bad Request.
			return
		}
		pages, ok := r.URL.Query()["key"]
		page := 1
		if !ok || len(pages[0]) < 1 {
			page = 1
		} else {
			page, _ = strconv.Atoi(pages[0])
		}
		var posts []models.Posts
		//Pagination Starts Here
		findOptions := options.Find()
		var perPage int64 = 2
		findOptions.SetSkip((int64(page) - 1) * perPage)
		findOptions.SetLimit(perPage)
		//Pagination Ends Here
		cur, err := collection.Collection("Posts").Find(context.TODO(), bson.M{"userid": objectId}, findOptions)
		if err != nil {
			helpers.GetError(err, w)
			return
		}
		defer cur.Close(context.TODO())
		for cur.Next(context.TODO()) {
			var post models.Posts
			err := cur.Decode(&post)
			if err != nil {
				log.Fatal(err)
			}
			posts = append(posts, post)
		}
		if err := cur.Err(); err != nil {
			log.Fatal(err)
		}
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(posts)
	} else {
		w.WriteHeader(405)
		return
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
