package utils

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/sisoputnfrba/tp-golang/utils/conexiones"
	"github.com/sisoputnfrba/tp-golang/utils/types"
)

func Iniciar_memoria(logger *slog.Logger) {
	mux := http.NewServeMux()

	// Endpoints
	mux.HandleFunc("POST /crear-proceso", Crear_proceso(logger))
	mux.HandleFunc("PATCH /finalizar-proceso/{pid}", Finalizar_proceso(logger))
	mux.HandleFunc("POST /CREAR_HILO", Crear_hilo(logger))
	mux.HandleFunc("PATCH /FINALIZAR_HILO", Finalizar_hilo(logger))

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
		logger.Info(fmt.Sprintf("Me llegaron los siguientes parametros para crear proceso: %+v", magic))

		// IMPORTANTE: Acá tiene que ir todo para que la memoria CREE el proceso (Está en pagina 20 y 21 del enunciado)

		// Si memoria pudo asignar el espacio necesario para el proceso responde con OK a Kernel
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}

func Finalizar_proceso(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pid := r.PathValue("pid") //Recibimos el pid a finalizar

		logger.Info(fmt.Sprintf("Liberando memoria de Proceso con PID = %+v", pid))

		// IMPORTANTE: Acá tiene que ir todo para que la memoria FINALICE el proceso (Está en pagina 21 del enunciado)

		// Si memoria pudo finalizar el proceso responde con OK a Kernel
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}

}

func Crear_hilo(logger *slog.Logger) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var infoHilo types.EnviarHiloAMemoria
		err := json.NewDecoder(r.Body).Decode(&infoHilo)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Aca va toda la logica para crear el hilo pag(22)

		pidParceado := strconv.Itoa(int(infoHilo.PID))
		logger.Info("## Hilo Creado - (PID:TID) - (%d:%d)", pidParceado, infoHilo.TID)

		// En caso de haberse creado el hilo

		respuesta, err := json.Marshal("OK")
		if err != nil {
			http.Error(w, "Error al codificar los datos como JSON", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(respuesta)

		logger.Info(fmt.Sprintf("## Hilo Creado - (PID:TID) - (%d:%d)", infoHilo.PID, infoHilo.TID))
	}
}

func Finalizar_hilo(logger *slog.Logger) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var infoHilo types.PIDTID
		err := json.NewDecoder(r.Body).Decode(&infoHilo)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Aca va toda la logica para finalizar el hilo

		pidParceado := strconv.Itoa(int(infoHilo.PID))
		logger.Info("## Hilo Destruido - (PID:TID) - (%d:%d)", pidParceado, infoHilo.TID)

		// En caso de haberse finalizado el hilo
		respuesta, err := json.Marshal("OK")
		if err != nil {
			http.Error(w, "Error al codificar los datos como JSON", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(respuesta)

		logger.Info(fmt.Sprintf("## Hilo Finalizado - (PID:TID) - (%d:%d)", infoHilo.PID, infoHilo.TID))
	}
}
