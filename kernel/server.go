package main

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/sisoputnfrba/tp-golang/utils/conexiones"
	"github.com/sisoputnfrba/tp-golang/utils/server"
)

func iniciar_kernel(logger *slog.Logger) {
	mux := http.NewServeMux()

	// Ejemplos de como agregar handlers
	mux.HandleFunc("/handshake", server.Recibir_handshake(logger))
	//mux.HandleFunc("/leer", leerHandler)
	//mux.HandleFunc("/escribir", escribirHandler)

	conexiones.LevantarServidor(strconv.Itoa(config.Port), mux, logger)

}
