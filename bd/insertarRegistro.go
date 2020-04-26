package bd


import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/MsolimanoMiranda/twittor/models"
)


func InsertarRegistro(u models.Usuario) (string,bool,error){

	ctx, cancel := context.WithTimeout(context.Background(),15*time.Second)
	defer cancel()

	db := MongoCN.Database("twittor")
	col := db.Database("usuarios")

	u.Password,_ = EncriptarPassword(u.Password)

	result, err := col.InsertOne(ctx,u)
	if err != nil {
		return "",false,err
	}

	ObjID,_ := result.InsertedID.(primitive.ObjectID)
	return ObjID.String(),true,nil
}