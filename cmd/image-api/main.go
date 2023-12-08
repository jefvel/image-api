package main

import (
	"image-api/pkg/handler"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()

	apiv1 := r.PathPrefix("v1").Subrouter()

	apiv1.HandleFunc("/images", handler.ListImages).Methods("GET")
	apiv1.HandleFunc("/images/{id}", handler.GetImage).Methods("GET")

	apiv1.HandleFunc("/images/{id}/data", handler.GetImageData).Methods("GET")

	apiv1.HandleFunc("/images", handler.UploadImage).Methods("POST")
	apiv1.HandleFunc("/images/{id}", handler.UpdateImage).Methods("PUT")

	loggedRouter := handlers.LoggingHandler(os.Stdout, r)

	err := http.ListenAndServe(":8000", loggedRouter)

	if err != nil {
		log.Fatalln(err)
	}
}
