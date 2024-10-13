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
	mux.HandleFunc("POST /FINALIZAR_HILO", Finalizar_hilo(logger))
	mux.HandleFunc("POST /memory-dump", Memory_dump(logger))
	conexiones.LevantarServidor(strconv.Itoa(Configs.Port), mux, logger)

}

//Coms con KERNEL

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

		// En caso de haberse finalizado el hilo
		respuesta, err := json.Marshal("OK")
		if err != nil {
			http.Error(w, "Error al codificar los datos como JSON", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		//w.Write([]byte("OK"))

		w.Write(respuesta)
		logger.Info(fmt.Sprintf("## Hilo Destruido - (PID:TID) - (%d:%d)", infoHilo.PID, infoHilo.TID))
	}
}

func Memory_dump(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//		decoder := json.NewDecoder(r.Body)
		/*
			if err != nil {
				logger.Error(fmt.Sprintf("Error al decodificar mensaje: %s\n", err.Error()))
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Error al decodificar mensaje"))
				return
			}
		*/
        var req
        err := json.NewDecoder(r.Body).Decode(&req)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}


        w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}

// Coms con CPU
func Actualizar_Contexto(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
        var req
        err := json.NewDecoder(r.Body).Decode(&req)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
        Actualizar_Contexto(req.PID, req.TID)
        
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
	
	}
}

func Obtener_Contexto(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RegCPU
        err := json.NewDecoder(r.Body).Decode(&req)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
        regCPU, err := Ver_Contexto(req.PID, req.TID)
        if err != nil {
            http.Error(w, err.Error(), http.StatusNotFound)
            return
        }
        ipCPU := "127.0.0.1" //Cambiar por referencia a archivo config
        puertoCPU := 8004




        exito := client.Enviar_Body(regCPU, conf, puertoCPU, endpointCPU, logger)
        if !exito {
            http.Error(w, "Error al enviar el contexto al CPU", http.StatusInternalServerError)
            return
        }
        
        w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))

	}
}

func Obtener_ID(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
        var req
        err := json.NewDecoder(r.Body).Decode(&req)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
        w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))

	}
}

func Read_Mem(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
        var req
        err := json.NewDecoder(r.Body).Decode(&req)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

    
        w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))

	}
}
func Write_Mem(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
    var req
        err := json.NewDecoder(r.Body).Decode(&req)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

        w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))

	}
}
