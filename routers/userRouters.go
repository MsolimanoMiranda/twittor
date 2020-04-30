package routers

import (
	usuarioController "github.com/MsolimanoMiranda/twittor/controller"
	"github.com/MsolimanoMiranda/twittor/middlew"
	"github.com/gorilla/mux"
)

func UserRouters(router *mux.Router) {
	s := router.PathPrefix("/user").Subrouter()
	s.HandleFunc("/login", middlew.ChequeoBD(usuarioController.LoginUsuario)).Methods("POST")
	s.HandleFunc("/registro", middlew.ChequeoBD(middlew.ValidoJWT(usuarioController.RegistroUsuario))).Methods("POST")
	s.HandleFunc("/perfil", middlew.ChequeoBD(middlew.ValidoJWT(usuarioController.VerPerfil))).Methods("GET")

}
