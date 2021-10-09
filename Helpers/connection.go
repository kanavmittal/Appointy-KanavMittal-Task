package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Database {
	adminOptions := options.Client().ApplyURI("mongodb+srv://kanav:kanavmittal@cluster0.rne9l.mongodb.net/UserDB?retryWrites=true&w=majority")
	admin, err := mongo.Connect(context.TODO(), adminOptions)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to Database is successful")
	collection := admin.Database("UserDB")
	return collection
}

type ErrorResponse struct {
	StatusCode   int    `json:"status"`
	ErrorMessage string `json:"message"`
}

func GetError(err error, w http.ResponseWriter) {
	log.Fatal(err.Error())
	var response = ErrorResponse{
		ErrorMessage: err.Error(),
		StatusCode:   http.StatusInternalServerError,
	}
	message, _ := json.Marshal(response)
	w.WriteHeader(response.StatusCode)
	w.Write(message)
}
