package utils

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/sisoputnfrba/tp-golang/utils/conexiones"
	"github.com/sisoputnfrba/tp-golang/utils/server"
)

func Iniciar_cpu(logger *slog.Logger) {
	mux := http.NewServeMux()

	// Aca van a ir todos los handlers

	mux.HandleFunc("/handshake", server.Recibir_handshake(logger))
	//mux.HandleFunc("/escribir", escribirHandler)

	conexiones.LevantarServidor(strconv.Itoa(Configs.Port), mux, logger)
}
