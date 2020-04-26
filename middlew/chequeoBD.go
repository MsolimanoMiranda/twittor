package middlew

import (
	"net/http"
	
	"github.com/MsolimanoMiranda/twittor/bd"
)


func ChequeoBD(next http.HandlerFunc) http.HandlerFunc{
	return func(w http.ResponseWriter,r *http.Request){
		if bd.CheckConnection() == 0 {
			http.Error(w,"Conxion perdida con la Base de Datos",500)
			return
		}
		next.ServeHTTP(w,r)
	}
}