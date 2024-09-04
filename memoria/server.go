package main

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/sisoputnfrba/tp-golang/utils/conexiones"
)

func iniciarMemoria(logger *slog.Logger) {
	mux := http.NewServeMux()

	//mux.HandleFunc("/leer", leerHandler)
	//mux.HandleFunc("/escribir", escribirHandler)

	conexiones.LevantarServidor(strconv.Itoa(config.Port), mux, logger)

}
