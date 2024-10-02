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

func Iniciar_cpu(logger *slog.Logger) {
	mux := http.NewServeMux()

	// Endpoints
	mux.HandleFunc("/mensaje", server.Recibir_handshake(logger))
	//mux.HandleFunc("POST /comunicacion-memoria", ComunicacionMemoria(logger))

	conexiones.LevantarServidor(strconv.Itoa(Configs.Port), mux, logger)

	//nuevo endpoint
	mux.HandleFunc("/wait_tid_pid", WAIT_FOR_TID_PID(logger))
}

// funcion que espera recibir tid y pid de kernel
func WAIT_FOR_TID_PID(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parseamos el cuerpo de la solicitud para obtener un PCB
		var pcb types.PCB
		if err := json.NewDecoder(r.Body).Decode(&pcb); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			logger.Error("Error decoding PCB: ", err)
			return
		}

		// Aquí podrías procesar los datos del PCB y sus TCBs
		for _, tcb := range pcb.TCBs {
			logger.Info(fmt.Sprintf("Recibido TID: %d, PID: %d", tcb.TID, tcb.PID))
			// Procesa el TCB según la lógica de tu simulador (ej. encolar, ejecutar, etc.)
		}

		// Responder con éxito
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("TID y PID recibidos"))
	}
}

/*
func ComunicacionMemoria(logger *slog.Logger) http.HandlerFunc {
	var request BodyRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	respuesta, err := json.Marshal(fmt.Sprintf("Hola %s! Como andas?", request.Name))
	if err != nil {
		http.Error(w, "Error al codificar los datos como JSON", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respuesta)
}
*/
