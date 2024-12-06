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
	"github.com/sisoputnfrba/tp-golang/utils/types"
)

func Inicializar_cpu(logger *slog.Logger) {
	mux := http.NewServeMux()

	// Endpoints de kernel
	mux.HandleFunc("POST /EJECUTAR_KERNEL", Recibir_PIDTID(logger))
	mux.HandleFunc("POST /INTERRUPCION_FIN_QUANTUM", RecibirInterrupcion(logger))
	mux.HandleFunc("POST /INTERRUPT", RecibirInterrupcion(logger))
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

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("PID y TID almacenados y CPU iniciada"))

		utils.Control = true
		// Llamar a Comenzar_cpu para iniciar el proceso de CPU
		cicloDeInstruccion.Comenzar_cpu(logger)

	}
}

// Función para recibir la interrupción y el TID desde la solicitud
func RecibirInterrupcion(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var bodyInterrupcion types.InterruptionInfo
		err := json.NewDecoder(r.Body).Decode(&bodyInterrupcion)
		if err != nil {
			logger.Error(fmt.Sprintf("Error al decodificar mensaje: %s\n", err.Error()))
		}

		// Almacenar el nombre de la interrupción y el TID en la variable
		var interrupcion = types.InterruptionInfo{
			NombreInterrupcion: bodyInterrupcion.NombreInterrupcion,
			TID:                bodyInterrupcion.TID,
		}

		cicloDeInstruccion.InterrupcionRecibida = &interrupcion

		// Log de confirmación
		logger.Info("## Llega interrupcion al puerto Interrupt")

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Interrupción y TID almacenados"))
	}
}
