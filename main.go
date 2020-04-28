package main
import (

	"log"
	
	"github.com/MsolimanoMiranda/twittor/handlers"
	"github.com/MsolimanoMiranda/twittor/bd"

)

func main() {

	if bd.CheckConnection() == 0 {
		log.Fatal("Sin conexion a la base de Datos")
	}
	
	handlers.Manejadores()



}
