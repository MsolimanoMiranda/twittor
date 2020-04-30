package helpers

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type StringResponse struct {
	Status int    `json:"status"`
	Data   string `json:"data"`
	Error  bool   `json:"error"`
}

type ObjectResponse struct {
	Status int           `json:"status"`
	Data   []interface{} `json:"data"`
	Error  bool          `json:"error"`
}

//MessageResponse para devolver un string
func MessageResponse(valor string, w http.ResponseWriter, status int) {

	w.Header().Set("Content-Type", "application/json")

	val := status
	var error_bol = false
	if val == 200 {
		error_bol = false
	} else {
		error_bol = true
	}
	response := StringResponse{
		Status: status,
		Data:   valor,
		Error:  error_bol,
	}
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

func ResponseObject(valor []interface{}, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")

	response := ObjectResponse{
		Status: http.StatusOK,
		Data:   valor,
		Error:  true,
	}
	// log.Println("error",response)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

//EncriptarPassword  es la funcion que me permite encriptar
func EncriptarPassword(pass string) (string, error) {

	costo := 8
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), costo)
	return string(bytes), err

}
