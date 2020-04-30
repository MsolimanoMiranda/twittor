package usuarioController

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/MsolimanoMiranda/twittor/bd"
	"github.com/MsolimanoMiranda/twittor/helpers"
	"github.com/MsolimanoMiranda/twittor/models"
	"github.com/MsolimanoMiranda/twittor/services"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func RegistroUsuario(w http.ResponseWriter, r *http.Request) {

	var t models.Usuario
	err := json.NewDecoder(r.Body).Decode(&t)

	if err != nil {
		helpers.MessageResponse("Error en los datos recibidos "+err.Error(), w, http.StatusBadRequest)
		return
	}

	if len(t.Email) == 0 {
		helpers.MessageResponse("El Email de usuario es requerido", w, http.StatusBadRequest)
		return
	}
	if len(t.Password) < 6 {
		helpers.MessageResponse("El Password debe tener minimo 6 caracters", w, http.StatusBadRequest)
		return
	}

	_, encontrado, _ := ChequeoYaExisteUsuario(t.Email)
	if encontrado == true {
		helpers.MessageResponse("Ya existe un usuario con es Email", w, http.StatusBadRequest)
		return
	}

	new_id, status, err := InsertarRegistro(t)
	if err != nil {
		helpers.MessageResponse("Error al insertar el usuario"+err.Error(), w, http.StatusBadRequest)
		return
	}
	if status == false {
		helpers.MessageResponse("No se ha logrado insertr el registro del Usuario", w, http.StatusBadRequest)
		return
	}

	helpers.MessageResponse("Se registro el usuario con ID "+new_id, w, http.StatusCreated)

}

func LoginUsuario(w http.ResponseWriter, r *http.Request) {

	var t models.Usuario
	err := json.NewDecoder(r.Body).Decode(&t)

	if err != nil {
		helpers.MessageResponse("Error en los datos recibidos "+err.Error(), w, http.StatusBadRequest)
		return
	}

	if len(t.Email) == 0 {
		helpers.MessageResponse("El Email de usuario es requerido", w, http.StatusBadRequest)
		return
	}
	if len(t.Password) == 0 {
		helpers.MessageResponse("El Password debe tener minimo 6 caracters", w, http.StatusBadRequest)
		return
	}

	_, encontrado, _ := ChequeoYaExisteUsuario(t.Email)
	if encontrado == false {
		helpers.MessageResponse("No existe un usuario con es Email", w, http.StatusBadRequest)
		return
	}

	usuario, status := Login(t.Email, t.Password)
	if err != nil {
		helpers.MessageResponse("Error al insertar el usuario"+err.Error(), w, http.StatusBadRequest)
		return
	}
	if status == false {
		helpers.MessageResponse("Usuario y/o Contraseña inválidos", w, http.StatusBadRequest)
		return
	}

	token, _ := services.GenerarToken(usuario)
	res := struct {
		User  models.Usuario `json:"user"`
		Token string         `json:"token"`
	}{
		usuario,
		string(token),
	}
	helpers.ResponseObject([]interface{}{res}, w)

	expirationTime := time.Now().Add(24 * time.Hour)
	http.SetCookie(w, &http.Cookie{
		Name:    "User",
		Value:   `` + string(token) + ``,
		Expires: expirationTime,
	})
}

func VerPerfil(w http.ResponseWriter, r *http.Request) {

	ID := r.URL.Query().Get("id")
	if len(ID) < 1 {
		helpers.MessageResponse("Debe enviar el parámetro ID", w, http.StatusBadRequest)
		return
	}

	perfil, err := BuscoPerfil(ID)
	if err != nil {
		helpers.MessageResponse("Ocurrió un error al intentar buscar el registro "+err.Error(), w, http.StatusBadRequest)
		return
	}

	helpers.ResponseObject([]interface{}{perfil}, w)
}

//InsertarRegistro funcion que perimite insentaar registro
func InsertarRegistro(u models.Usuario) (string, bool, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := bd.MongoCN.Database("twittor")
	col := db.Collection("usuarios")

	u.Password, _ = helpers.EncriptarPassword(u.Password)

	result, err := col.InsertOne(ctx, u)
	if err != nil {
		return "", false, err
	}

	ObjID, _ := result.InsertedID.(primitive.ObjectID)
	return ObjID.String(), true, nil
}

// ChequeoYaExisteUsuario es la funcion que verifica que exista
func ChequeoYaExisteUsuario(email string) (models.Usuario, bool, string) {

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := bd.MongoCN.Database("twittor")
	col := db.Collection("usuarios")

	condicion := bson.M{"email": email}

	var resultado models.Usuario

	err := col.FindOne(ctx, condicion).Decode(&resultado)
	ID := resultado.ID.Hex()

	if err != nil {
		return resultado, false, ID
	}

	return resultado, true, ID
}

func Login(email string, password string) (models.Usuario, bool) {

	usu, encontrado, _ := ChequeoYaExisteUsuario(email)

	if encontrado == false {
		return usu, false
	}

	passwordBytes := []byte(password)
	passwordBD := []byte(usu.Password)
	err := bcrypt.CompareHashAndPassword(passwordBD, passwordBytes)

	if err != nil {
		return usu, false
	}

	return usu, true

}

func BuscoPerfil(ID string) (models.Usuario, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := bd.MongoCN.Database("twittor")
	col := db.Collection("usuarios")

	var resultado models.Usuario

	objId, _ := primitive.ObjectIDFromHex(ID)
	condicion := bson.M{"_id": objId}

	err := col.FindOne(ctx, condicion).Decode(&resultado)
	resultado.Password = ""

	if err != nil {
		return resultado, err
	}

	return resultado, nil

}
