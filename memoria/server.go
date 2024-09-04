package main

import (
	"net/http"

	"github.com/sisoputnfrba/tp-golang/utils/conexiones"
)

func iniciarMemoria() {
	mux := http.NewServeMux()

	//mux.HandleFunc("/leer", leerHandler)
	//mux.HandleFunc("/escribir", escribirHandler)

	conexiones.LevantarServidor("8081", mux)

}
