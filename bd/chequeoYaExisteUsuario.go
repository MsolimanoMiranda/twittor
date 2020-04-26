package bd


import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/MsolimanoMiranda/twittor/models"
)


func InsertarRegistro(email string) (models.Usuario,bool,string){

	ctx, cancel := context.WithTimeout(context.Background(),15*time.Second)
	defer cancel()

	db := MongoCN.Database("twittor")
	col := db.Database("usuarios")

	condicion := bson.M{"email":email}

	var resultado models.Usuario

	err := col.FindOne(ctx, condicion).Decode(&resultado)
	ID 	:= resultado.ID.hex()

	result, err := col.InsertOne(ctx,u)
	if err != nil {
		return resultado,false,ID
	}

	return resultado,true,ID
}