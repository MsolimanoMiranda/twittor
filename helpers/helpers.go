package helpers

import (
	"encoding/json"
	"net/http"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type StringResponse struct {
	Status int    `json:"status"`
	Data   string `json:"data"`
	Error  bool   `json:"error"`
}

type ObjectResponse struct {
	Status int    `json:"status"`
	Data   bson.M `json:"data"`
	Error  bool   `json:"error"`
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

func ResponseObject(valor interface{}, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	ret := bson.M{"status": http.StatusOK, "data": valor, "error": false}

	// response := ObjectResponse{
	// 	Status: http.StatusOK,
	// 	Data:   ret,
	// 	Error:  true,
	// }
	// log.Println("error",response)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ret)
}

//EncriptarPassword  es la funcion que me permite encriptar
func EncriptarPassword(pass string) (string, error) {

	costo := 8
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), costo)
	return string(bytes), err

}

func ConvertToMap(model interface{}) bson.M {
	ret := bson.M{}

	modelReflect := reflect.ValueOf(model)

	if modelReflect.Kind() == reflect.Ptr {
		modelReflect = modelReflect.Elem()
	}

	modelRefType := modelReflect.Type()
	fieldsCount := modelReflect.NumField()

	var fieldData interface{}

	for i := 0; i < fieldsCount; i++ {
		field := modelReflect.Field(i)

		switch field.Kind() {
		case reflect.Struct:
			fallthrough
		case reflect.Ptr:
			fieldData = ConvertToMap(field.Interface())
		default:
			fieldData = field.Interface()
		}

		ret[modelRefType.Field(i).Name] = fieldData
	}

	return ret
}
