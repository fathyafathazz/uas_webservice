package main

import (
	"fmt"
	"uas_musik/controller/auth"
	"uas_musik/controller/artist"
	"uas_musik/controller/album"
	"uas_musik/database"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	database.InitDB()
	fmt.Println("Hello World")

	router := mux.NewRouter()

	router.HandleFunc("/regis", auth.Registration).Methods("POST")
	router.HandleFunc("/login", auth.Login).Methods("POST")

	//Router handler artist
	router.HandleFunc("/artists", artist.GetArtist).Methods("GET")
	router.HandleFunc("/artists", auth.JWTAuth(artist.PostArtist)).Methods("POST")
	router.HandleFunc("/artists/{id}", auth.JWTAuth(artist.PutArtist)).Methods("PUT")
	router.HandleFunc("/artists/{id}", auth.JWTAuth(artist.DeleteArtist)).Methods("DELETE")

	//Router handler Album
	router.HandleFunc("/albums", album.GetAlbum).Methods("GET")
	router.HandleFunc("/albums", auth.JWTAuth(album.PostAlbum)).Methods("POST")
	router.HandleFunc("/albums/{id}", auth.JWTAuth(album.PutAlbum)).Methods("PUT")
	router.HandleFunc("/albums/{id}", auth.JWTAuth(album.DeleteAlbum)).Methods("DELETE")

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://127.0.0.1:5500"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
		Debug: true,
	})

	handler := c.Handler(router)

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}