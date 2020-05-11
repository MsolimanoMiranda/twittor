package usuarioController

import (
	"context"
	"encoding/json"
	"log"
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

	helpers.ResponseObject(res, w)

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

	// helpers.ResponseObject([]interface{}{perfil}, w)
	helpers.ResponseObject(perfil, w)

}

func ListarUsuario(w http.ResponseWriter, r *http.Request) {

	usuarios, err := ListUser()
	if err == false {
		helpers.MessageResponse("Ocurrió un error al listar Usuarios ", w, http.StatusBadRequest)
		return
	}

	// helpers.ResponseObject([]interface{}{perfil}, w)
	helpers.ResponseObject(usuarios, w)

}

func ModificarPerfil(w http.ResponseWriter, r *http.Request) {

	var t models.Usuario

	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, "Datos Incorrectos "+err.Error(), 400)
		return
	}

	var status bool

	status, err = ModificoRegistro(t, services.IDUsuario)
	log.Println(services.IDUsuario)
	if err != nil {
		helpers.MessageResponse("Ocurrión un error al intentar modificar el registro. Reintente nuevamente "+err.Error(), w, http.StatusBadRequest)
		return
	}

	if status == false {
		helpers.MessageResponse("No se ha logrado modificar el registro del usuario", w, http.StatusBadRequest)
		return
	}

	helpers.MessageResponse("Se Modifico con exito", w, http.StatusCreated)

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

func ListUser() ([]*models.Usuario, bool) {

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := bd.MongoCN.Database("twittor")
	col := db.Collection("usuarios")
	var results []*models.Usuario

	var incluir bool
	cur, err := col.Find(ctx, bson.M{})

	for cur.Next(ctx) {
		var s models.Usuario
		err := cur.Decode(&s)
		if err != nil {
			return results, false
		}
		incluir = true
		if incluir == true {
			s.Password = ""
			s.Biografia = ""
			s.SitioWeb = ""
			s.Ubicacion = ""
			s.Banner = ""
			// s.Email = ""

			results = append(results, &s)
		}
	}

	err = cur.Err()
	if err != nil {
		return results, false
	}
	cur.Close(ctx)
	return results, true
}

func ModificoRegistro(u models.Usuario, ID string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := bd.MongoCN.Database("twittor")
	col := db.Collection("usuarios")

	registro := make(map[string]interface{})
	if len(u.Nombre) > 0 {
		registro["nombre"] = u.Nombre
	}
	if len(u.Apellidos) > 0 {
		registro["apellidos"] = u.Apellidos
	}
	registro["fechaNacimiento"] = u.FechaNacimiento
	if len(u.Avatar) > 0 {
		registro["avatar"] = u.Avatar
	}
	if len(u.Banner) > 0 {
		registro["banner"] = u.Banner
	}
	if len(u.Biografia) > 0 {
		registro["biografia"] = u.Biografia
	}
	if len(u.Ubicacion) > 0 {
		registro["ubicacion"] = u.Ubicacion
	}
	if len(u.SitioWeb) > 0 {
		registro["sitioWeb"] = u.SitioWeb
	}

	updtString := bson.M{
		"$set": registro,
	}

	objID, _ := primitive.ObjectIDFromHex(ID)
	filtro := bson.M{"_id": bson.M{"$eq": objID}}

	_, err := col.UpdateOne(ctx, filtro, updtString)
	if err != nil {
		return false, err
	}

	return true, nil
}
