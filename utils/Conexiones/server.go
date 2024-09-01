package server

// Importamos librerias
import (
	"log"
	"net/http"
)

func initServer(port string, handler http.Handler) {
	log.Printf("Levantando servidor en el puerto %s...\n", port)
	err := http.ListenAndServe(":"+port, handler)

	//Manejo de errores
	if err != nil {
		log.Fatalf("Error al levantar el servidor: %v", err)
	}
}
