package handlers

import (
	"net/http"
	"os"
	"log"

	"github.com/MsolimanoMiranda/twittor/middlew"
	"github.com/MsolimanoMiranda/twittor/routers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

//MANEJADORES SETEO PUERTO,Handler y pongo a escuchar el server
func Manejadores(){

		router := mux.NewRouter()
		router.HandleFunc("/registro",middlew.ChequeoBD(routers.Registro)).Methods("POST")

		PORT := os.Getenv("PORT")

		if PORT == ""{
			PORT = "8080"
		}

		handler := cors.AllowAll().Handler(router)
		log.Fatal(http.ListenAndServe(":"+PORT,handler))

}