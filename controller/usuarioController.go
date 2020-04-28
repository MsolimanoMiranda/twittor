package usuarioController


import (
	"context"
	"time"
	
	"github.com/MsolimanoMiranda/twittor/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/MsolimanoMiranda/twittor/helpers"
	"github.com/MsolimanoMiranda/twittor/bd"
	"go.mongodb.org/mongo-driver/bson"
	
	
)

//InsertarRegistro funcion que perimite insentaar registro
func InsertarRegistro(u models.Usuario) (string,bool,error){

	ctx, cancel := context.WithTimeout(context.Background(),15*time.Second)
	defer cancel()

	db := bd.MongoCN.Database("twittor")
	col := db.Collection("usuarios")

	u.Password,_ = helpers.EncriptarPassword(u.Password)

	result, err := col.InsertOne(ctx,u)
	if err != nil {
		return "",false,err
	}

	ObjID,_ := result.InsertedID.(primitive.ObjectID)
	return ObjID.String(),true,nil
}

// ChequeoYaExisteUsuario es la funcion que verifica que exista
func ChequeoYaExisteUsuario(email string) (models.Usuario,bool,string){

	ctx, cancel := context.WithTimeout(context.Background(),15*time.Second)
	defer cancel()

	db := bd.MongoCN.Database("twittor")
	col := db.Collection("usuarios")

	condicion := bson.M{"email":email}

	var resultado models.Usuario

	err := col.FindOne(ctx, condicion).Decode(&resultado)
	ID 	:= resultado.ID.Hex()
	
	if err != nil {
		return resultado,false,ID
	}

	return resultado,true,ID
}