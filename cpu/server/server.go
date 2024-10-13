package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/sisoputnfrba/tp-golang/cpu/utils"
	"github.com/sisoputnfrba/tp-golang/utils/conexiones"
	"github.com/sisoputnfrba/tp-golang/utils/server"
	"github.com/sisoputnfrba/tp-golang/utils/types"
)

func Iniciar_cpu(logger *slog.Logger) {
	mux := http.NewServeMux()

	// Endpoints
	mux.HandleFunc("/mensaje", server.Recibir_handshake(logger))
	mux.HandleFunc("POST /EJECUTAR_KERNEL", WAIT_FOR_TID_PID(logger))
	mux.HandleFunc("POST /INTERRUPCION_FIN_QUANTUM", ReciboInterrupcionTID(logger))

	//mux.HandleFunc("POST /comunicacion-memoria", ComunicacionMemoria(logger))

	conexiones.LevantarServidor(strconv.Itoa(utils.Configs.Port), mux, logger)

}

// Variable global para almacenar el PID y TID
var ReceivedPIDTID *types.PIDTID = nil

var ReceivedInterrupt uint32

// Getter para acceder a la variable global ReceivedPIDTID
//func GetReceivedPIDTID() *types.PIDTID {
//	return receivedPIDTID
//}

// Handler para esperar el PID y TID del Kernel
func WAIT_FOR_TID_PID(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var pidtid types.PIDTID

		// Intentamos decodificar el cuerpo de la solicitud
		err := decoder.Decode(&pidtid)
		if err != nil {
			// Log de error en caso de fallo al decodificar
			logger.Error(fmt.Sprintf("Error al decodificar mensaje: %s\n", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error al decodificar mensaje"))
			return
		}

		// Log de la información recibida si la decodificación fue exitosa
		logger.Info(fmt.Sprintf("Recibido TID: %d, PID: %d", pidtid.TID, pidtid.PID))

		// Asignar a la variable global
		ReceivedPIDTID = &pidtid

		// Responder con éxito
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("TID y PID recibidos"))
	}
}

func ReciboInterrupcionTID(Logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var interrupt uint32

		// Intentamos decodificar el cuerpo de la solicitud
		err := decoder.Decode(&interrupt)
		if err != nil {
			// Log de error en caso de fallo al decodificar
			Logger.Error(fmt.Sprintf("Error al decodificar mensaje: %s\n", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error al decodificar mensaje"))
			return
		}

		// Log de la información recibida si la decodificación fue exitosa
		Logger.Info(fmt.Sprintf("Recibido interrupcion: %d", interrupt))

		// Asignar a la variable global
		ReceivedInterrupt = interrupt

		// Responder con éxito
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("TID  recibido"))
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
