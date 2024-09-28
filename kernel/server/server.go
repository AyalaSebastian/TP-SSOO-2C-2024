package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/sisoputnfrba/tp-golang/kernel/planificador"
	"github.com/sisoputnfrba/tp-golang/kernel/utils"
	"github.com/sisoputnfrba/tp-golang/utils/conexiones"
	"github.com/sisoputnfrba/tp-golang/utils/types"
)

func Iniciar_kernel(logger *slog.Logger) {
	mux := http.NewServeMux()

	// Endpoints
	mux.HandleFunc("POST /PROCESS_CREATE", PROCESS_CREATE(logger))
	mux.HandleFunc("PUT /PROCESS_EXIT", PROCESS_DESTROY(logger))
	conexiones.LevantarServidor(strconv.Itoa(utils.Configs.Port), mux, logger)

}

func PROCESS_CREATE(logger *slog.Logger) http.HandlerFunc {
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
		logger.Info(fmt.Sprintf("Me llegaron los siguientes parametros: %+v", magic))

		// Aca va el desarrollo

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}

func PROCESS_DESTROY(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pid := r.PathValue("pid")                //Recibimos el pid a finalizar
		val, _ := strconv.ParseUint(pid, 10, 32) //Convierto el pid a uint32 ya que viene en String
		parsePid := uint32(val)
		planificador.Finalizar_proceso(parsePid, logger)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}
