package utils

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/sisoputnfrba/tp-golang/utils/conexiones"
	"github.com/sisoputnfrba/tp-golang/utils/server"
)

func Iniciar_fileSystem(logger *slog.Logger) {
	mux := http.NewServeMux()

	// Endpoints
	mux.HandleFunc("/handshake", server.Recibir_handshake(logger))

	conexiones.LevantarServidor(strconv.Itoa(Configs.Port), mux, logger)
}
