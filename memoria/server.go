package main

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/sisoputnfrba/tp-golang/utils/conexiones"
	"github.com/sisoputnfrba/tp-golang/utils/server"
)

func iniciarMemoria(logger *slog.Logger) {
	mux := http.NewServeMux()

	mux.HandleFunc("/handshake", server.Recibir_handshake(logger))
	//mux.HandleFunc("/leer", leerHandler)
	//mux.HandleFunc("/escribir", escribirHandler)

	conexiones.LevantarServidor(strconv.Itoa(config.Port), mux, logger)

}
