package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/sisoputnfrba/tp-golang/cpu/cicloDeInstruccion"
	"github.com/sisoputnfrba/tp-golang/cpu/utils"
	"github.com/sisoputnfrba/tp-golang/utils/conexiones"
	"github.com/sisoputnfrba/tp-golang/utils/server"
	"github.com/sisoputnfrba/tp-golang/utils/types"
)

func Inicializar_cpu(logger *slog.Logger) {
	mux := http.NewServeMux()

	//seguro borraremos esto
	mux.HandleFunc("/mensaje", server.Recibir_handshake(logger))

	// Endpoints de kernel
	mux.HandleFunc("POST /EJECUTAR_KERNEL", Recibir_PIDTID(logger))
	mux.HandleFunc("POST /INTERRUPCION_FIN_QUANTUM", ReciboInterrupcionTID(logger))
	mux.HandleFunc("POST /INTERRUPT", ReciboInterrupcionTID(logger))
	//mux.HandleFunc("POST /comunicacion-memoria", ComunicacionMemoria(logger))

	conexiones.LevantarServidor(strconv.Itoa(utils.Configs.Port), mux, logger)

}

// SetGlobalPIDTID recibe un PIDTID y actualiza las variables globales PID y TID.
func Recibir_PIDTID(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var pidtid types.PIDTID

		// Decodificar el cuerpo de la solicitud JSON
		if err := json.NewDecoder(r.Body).Decode(&pidtid); err != nil {
			http.Error(w, "Error al decodificar el JSON de la solicitud", http.StatusBadRequest)
			logger.Error("Error al decodificar JSON", slog.String("error", err.Error()))
			return
		}

		// Almacenar el PID y TID en la variable global
		cicloDeInstruccion.GlobalPIDTID = pidtid

		// Log de confirmación de la actualización
		logger.Info("PID y TID actualizados", slog.Any(
			"PID", pidtid.PID), slog.Any("TID", pidtid.TID))

		// Llamar a Comenzar_cpu para iniciar el proceso de CPU

		cicloDeInstruccion.Comenzar_cpu(logger)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("PID y TID almacenados y CPU iniciada"))
	}
}

var ReceivedInterrupt uint32

func ReciboInterrupcionTID(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var interrupt uint32

		// Intentamos decodificar el cuerpo de la solicitud
		err := decoder.Decode(&interrupt)
		if err != nil {
			// Log de error en caso de fallo al decodificar
			logger.Error(fmt.Sprintf("Error al decodificar mensaje: %s\n", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error al decodificar mensaje"))
			return
		}

		// Log de la información recibida si la decodificación fue exitosa
		logger.Info(("## Llega interrupcion al puerto Interrupt"))

		// Asignar a la variable global
		ReceivedInterrupt = interrupt
		// Responder con éxito
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("TID  recibido"))

	}
}
