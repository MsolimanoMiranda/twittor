package bd

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var host = "localhost"
var port = 27017
var MongoCN = ConectarBD()
var clienteOptions = options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%d", host, port))

//ConectarBD  es la funcion que me permite conectar la bd
func ConectarBD() *mongo.Client {

	client, err := mongo.Connect(context.TODO(), clienteOptions)
	if err != nil {
		log.Fatal(err.Error())
		return client
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err.Error())
		return client

	}
	log.Println("Conexion Exitosa con la DB")

	return client

}

//CheckConnection db es la funcion que me permite conectar la bd
func CheckConnection() int {

	err := MongoCN.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err.Error())
		return 0
	}
	return 1

}
