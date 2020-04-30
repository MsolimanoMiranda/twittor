package middlew

import (
	"net/http"

	"github.com/MsolimanoMiranda/twittor/bd"
	"github.com/MsolimanoMiranda/twittor/helpers"
	"github.com/MsolimanoMiranda/twittor/services"
)

func ChequeoBD(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if bd.CheckConnection() == 0 {

			helpers.MessageResponse("Conxion perdida con la Base de Datos", w, 500)
			return
		}
		next.ServeHTTP(w, r)
	}
}

/*ValidoJWT permite validar el JWT que nos viene en la petición */
func ValidoJWT(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _, _, err := services.ValidarToken(r.Header.Get("Authorization"))
		if err != nil {
			helpers.MessageResponse("Error en el Token ! "+err.Error(), w, http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	}
}
