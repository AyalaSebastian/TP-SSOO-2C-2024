package utils

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/sisoputnfrba/tp-golang/utils/conexiones"
	"github.com/sisoputnfrba/tp-golang/utils/server"
	"github.com/sisoputnfrba/tp-golang/utils/types"
)

func Iniciar_memoria(logger *slog.Logger) {
	mux := http.NewServeMux()

	// Endpoints
	mux.HandleFunc("/handshake", server.Recibir_handshake(logger))
	mux.HandleFunc("POST /crear-proceso", Crear_proceso(logger))

	conexiones.LevantarServidor(strconv.Itoa(Configs.Port), mux, logger)

}

func Crear_proceso(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var magic types.PathTamanio
		err := decoder.Decode(&magic)
		if err != nil {
			logger.Error(fmt.Sprintf("Error al decodificar mensaje: %s\n", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error al decodificar mensaje"))
			return
		}
		logger.Info(fmt.Sprintf("Me llegaron los siguientes datos: %+v", magic))

		// IMPORTANTE: Acá tiene que ir todo para que la memoria cree el proceso (Está en pagina 20 y 21 del enunciado)

		// Si memoria pudo asignar el espacio necesario para el proceso responde con OK a Kernel
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}
