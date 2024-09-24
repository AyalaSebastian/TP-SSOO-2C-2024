package server

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/sisoputnfrba/tp-golang/kernel/utils"
	"github.com/sisoputnfrba/tp-golang/utils/conexiones"
	"github.com/sisoputnfrba/tp-golang/utils/server"
)

func Iniciar_kernel(logger *slog.Logger) {
	mux := http.NewServeMux()

	// Endpoints
	mux.HandleFunc("/handshake", server.Recibir_handshake(logger))

	conexiones.LevantarServidor(strconv.Itoa(utils.Configs.Port), mux, logger)

}
