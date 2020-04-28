package routers

import (
	"encoding/json"
	"net/http"

	"github.com/MsolimanoMiranda/twittor/controller"
	"github.com/MsolimanoMiranda/twittor/models"	
	"github.com/MsolimanoMiranda/twittor/helpers"

)


func Registro(w http.ResponseWriter,r *http.Request){

	var t models.Usuario
	err := json.NewDecoder(r.Body).Decode(&t)
	

	if err != nil {
		helpers.MessageResponse("Error en los datos recibidos "+err.Error(),w,http.StatusBadRequest)
		return
	}

	if len(t.Email)==0{
		helpers.MessageResponse("El Email de usuario es requerido",w,http.StatusBadRequest)
		return
	}
	if len(t.Password)<6{
		helpers.MessageResponse("El Password debe tener minimo 6 caracters",w,http.StatusBadRequest)
		return
	}


	_,encontrado,_ := usuarioController.ChequeoYaExisteUsuario(t.Email)
	if encontrado == true {
		helpers.MessageResponse("Ya existe un usuario con es Email",w,http.StatusBadRequest)
		return
	}

	new_id,status,err := usuarioController.InsertarRegistro(t)
	if err != nil {
		helpers.MessageResponse("Error al insertar el usuario"+err.Error(),w,http.StatusBadRequest)
		return
	}
	if status == false {
		helpers.MessageResponse("No se ha logrado insertr el registro del Usuario",w,http.StatusBadRequest)
		return
	}

	helpers.MessageResponse("Se registro el usuario con ID "+new_id,w,http.StatusCreated)

}

