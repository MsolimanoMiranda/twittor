package routers

import (
	"encoding/json"
	"net/http"

	"github.com/MsolimanoMiranda/twittor/bd"
	"github.com/MsolimanoMiranda/twittor/models"	
)


func Registro(w http.ResponseWriter,r *http.Request){

	var t models.Usuario
	err := json.NewDecoder(r.Body).Decode(&t)

	if err != nil {
		http.Error(w,"Error en los datos recibidos "+err.Error(), 400)
		return
	}

	if len(t.Email)==0{
		http.Error(w,"El Email de usuario es requerido", 400)
		return
	}
	if len(t.Password)<6{
		http.Error(w,"El Password debe tener minimo 6 caracters ", 400)
		return
	}

	_,encontrado,_ := bd.ChequeoYaExisteUsuario(t.Email)
	if encontrado == true {
		http.Error(w,"Ya existe un usuario con es Email", 400)
		return
	}

	_,status,err := bd.InsertarRegistro(t)
	if err != nil {
		http.Error(w,"Error al insertar el usuario"+err.Error(), 400)
		return
	}
	if status == false {
		http.Error(w,"No se ha logrado insertr el registro del Usuario", 400)
		return
	}

	w.WriteHeader(http.StatusCreated)

}