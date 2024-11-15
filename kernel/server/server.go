package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/sisoputnfrba/tp-golang/kernel/client"
	"github.com/sisoputnfrba/tp-golang/kernel/planificador"
	"github.com/sisoputnfrba/tp-golang/kernel/utils"
	"github.com/sisoputnfrba/tp-golang/utils/conexiones"
	"github.com/sisoputnfrba/tp-golang/utils/types"
)

func Iniciar_kernel(logger *slog.Logger) {
	mux := http.NewServeMux()

	// Endpoints
	mux.HandleFunc("POST /PROCESS_CREATE", PROCESS_CREATE(logger))
	mux.HandleFunc("PUT /PROCESS_EXIT", PROCESS_EXIT(logger))
	mux.HandleFunc("POST /THREAD_CREATE", THREAD_CREATE(logger))
	mux.HandleFunc("PATCH /THREAD_JOIN/{tid}", THREAD_JOIN(logger))
	mux.HandleFunc("DELETE /THREAD_CANCEL/{tid}", THREAD_CANCEL(logger))
	mux.HandleFunc("DELETE /THREAD_EXIT", THREAD_EXIT(logger))
	mux.HandleFunc("POST /DUMP_MEMORY", DUMP_MEMORY(logger))
	mux.HandleFunc("PATCH /dump_response/{response}", Respuesta_dump(logger))
	mux.HandleFunc("POST /MUTEX_CREATE/{mutex}", MUTEX_CREATE(logger))
	mux.HandleFunc("PATCH /MUTEX_LOCK/{mutex}", MUTEX_LOCK(logger))
	mux.HandleFunc("PATCH /MUTEX_UNLOCK/{mutex}", MUTEX_UNLOCK(logger))
	mux.HandleFunc("PUT /IO/{ms}", IO(logger))

	mux.HandleFunc("PUT /recibir-desalojo", Recibir_desalojo(logger))

	conexiones.LevantarServidor(strconv.Itoa(utils.Configs.Port), mux, logger)

}

// Syscalls referidas a procesos

func PROCESS_CREATE(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info(fmt.Sprintf("## (%d:%d) - Solicitó syscall: PROCESS_CREATE", utils.Execute.PID, utils.Execute.TID))
		decoder := json.NewDecoder(r.Body)
		var magic types.ProcessCreateParams
		err := decoder.Decode(&magic)
		if err != nil {
			logger.Error(fmt.Sprintf("Error al decodificar mensaje: %s", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error al decodificar mensaje"))
			return
		}
		planificador.Crear_proceso(magic.Path, magic.Tamanio, magic.Prioridad, logger)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}

func PROCESS_EXIT(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info(fmt.Sprintf("## (%d:%d) - Solicitó syscall: PROCESS_EXIT", utils.Execute.PID, utils.Execute.TID))
		planificador.Finalizar_proceso(utils.Execute.PID, logger)
		utils.Execute = nil
		utils.Planificador.Signal()

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}

func DUMP_MEMORY(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info(fmt.Sprintf("## (%d:%d) - Solicitó syscall: DUMP_MEMORY", utils.Execute.PID, utils.Execute.TID))
		parametros := types.PIDTID{TID: utils.Execute.TID, PID: utils.Execute.PID} // Saco el pid y el tid del hilo que esta ejecutando
		utils.Execute = nil
		bloqueado := utils.Bloqueado{PID: parametros.PID, TID: parametros.TID, Motivo: utils.DUMP}
		utils.Encolar(&planificador.ColaBlocked, bloqueado)
		utils.Planificador.Signal()

		client.Enviar_Body(parametros, utils.Configs.IpMemory, utils.Configs.PortMemory, "memory-dump", logger)
		logger.Info(fmt.Sprintf("## (%d:%d) - Bloqueado por: DUMP MEMORY", utils.Execute.PID, utils.Execute.TID))

		w.WriteHeader(http.StatusOK)
	}
}

func Respuesta_dump(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := r.PathValue("response")
		defer w.WriteHeader(http.StatusOK)

		if response == "OK" {
			desbloqueado := utils.Desencolar_Por_Motivo(&planificador.ColaBlocked, utils.DUMP)
			pcb := utils.Obtener_PCB_por_PID(desbloqueado.PID)
			tcb := pcb.TCBs[desbloqueado.TID]
			utils.Encolar_ColaReady(planificador.ColaReady, tcb)
		} else {
			desbloqueado := utils.Desencolar_Por_Motivo(&planificador.ColaBlocked, utils.DUMP)
			pcb := utils.Obtener_PCB_por_PID(desbloqueado.PID)
			tcb := pcb.TCBs[desbloqueado.TID]
			planificador.Finalizar_proceso(tcb.PID, logger)
			logger.Info(fmt.Sprintf("## Finaliza el proceso %d", desbloqueado.PID))
		}
	}
}

// Syscalls referidas a hilos

// BODY - VERBO POST
func THREAD_CREATE(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info(fmt.Sprintf("## (%d:%d) - Solicitó syscall: THREAD_CREATE", utils.Execute.PID, utils.Execute.TID))
		// Agarramos los parametros del body
		var params types.ThreadCreateParams
		err := json.NewDecoder(r.Body).Decode(&params)
		if err != nil {
			logger.Error(fmt.Sprintf("Error al decodificar mensaje: %s\n", err.Error()))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Creamos el hilo
		planificador.Crear_hilo(params.Path, params.Prioridad, logger)

		// Respondemos con un OK
		respuesta, err := json.Marshal("OK")
		if err != nil {
			http.Error(w, "Error al codificar mensaje como JSON", http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(respuesta)

	}
}

// QUERY PATH - VERBO DELETE
func THREAD_EXIT(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info(fmt.Sprintf("## (%d:%d) - Solicitó syscall: THREAD_EXIT", utils.Execute.PID, utils.Execute.TID))
		// Eliminamos el hilo que esta ejecutando actualmente
		logger.Info("Se recibio un THREAD_EXIT")

		// Finalizamos el hilo
		planificador.Finalizar_hilo(utils.Execute.TID, utils.Execute.PID, logger)

		// Respondemos con un OK
		respuesta, err := json.Marshal("OK")
		if err != nil {
			http.Error(w, "Error al codificar mensaje como JSON", http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(respuesta)
	}
}

// QUERY PATH - VERBO DELETE
// EL TID A ELIMINAR SE MANDA POR PARAMETRO DE LA URL EJ: /THREAD_CANCEL/1
func THREAD_CANCEL(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info(fmt.Sprintf("## (%d:%d) - Solicitó syscall: THREAD_CANCEL", utils.Execute.PID, utils.Execute.TID))
		// Tomamos el valor del tid de la variable de la URL
		tid := r.PathValue("tid")

		// Finalizamos el hilo
		tidNum, err := strconv.Atoi(tid)
		if err != nil {
			http.Error(w, "Error al convertir el TID a numero", http.StatusBadRequest)
			return
		}

		planificador.Finalizar_hilo(uint32(tidNum), utils.Execute.PID, logger)

		// Respondemos con un OK
		respuesta, err := json.Marshal("OK")
		if err != nil {
			http.Error(w, "Error al codificar mensaje como JSON", http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(respuesta)
	}
}

// QUERY PATH - VERBO PATCH
// EL TID A QUE BLOQUEA SE MANDA POR PARAMETRO DE LA URL EJ: /THREAD_JOIN/1
// SI NO EXISTE EL TID, SE RESPONDE CON "CONTINUAR_EJECUCION", SI EXISTE SE RESPONDE CON "OK"
func THREAD_JOIN(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info(fmt.Sprintf("## (%d:%d) - Solicitó syscall: THREAD_JOIN", utils.Execute.PID, utils.Execute.TID))
		// Tomamos el valor del tid de la variable de la URL
		tid := r.PathValue("tid")

		// Verificamos  que el tid exista actualmente
		// planificador.Crear_hilo(params.Path, params.Prioridad, logger)
		tidNum, err := strconv.Atoi(tid)
		if err != nil {
			http.Error(w, "Error al convertir el TID a numero", http.StatusBadRequest)
		}
		_, existe := utils.MapaPCB[utils.Execute.PID].TCBs[uint32(tidNum)]

		if !existe {
			respuesta, err := json.Marshal("CONTINUAR_EJECUCION")
			if err != nil {
				http.Error(w, "Error al codificar mensaje como JSON", http.StatusInternalServerError)
			}
			w.WriteHeader(http.StatusOK)
			w.Write(respuesta)
			return
		}

		// Mandamos el hilo a block
		utils.Execute = nil
		bloqueado := utils.Bloqueado{PID: utils.Execute.PID, TID: utils.Execute.TID, Motivo: utils.THREAD_JOIN, QuienFue: tid}
		utils.Encolar(&planificador.ColaBlocked, bloqueado)
		utils.Planificador.Signal()
		logger.Info(fmt.Sprintf("## (%d:%d) - Bloqueado por: THREAD_JOIN", utils.Execute.PID, utils.Execute.TID))

		respuesta, err := json.Marshal("OK")
		if err != nil {
			http.Error(w, "Error al codificar mensaje como JSON", http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(respuesta)
	}
}

// Syscall referida a MUTEX

// QUERY PATH - VERBO POST
// Si ya existe responde con "MUTEX_YA_EXISTE", si no existe lo crea y responde con "OK"
func MUTEX_CREATE(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info(fmt.Sprintf("## (%d:%d) - Solicitó syscall: MUTEX_CREATE", utils.Execute.PID, utils.Execute.TID))
		// Tomamos el valor del tid de la variable de la URL
		mutexName := r.PathValue("mutex")

		// Creamos el mutex
		_, existe := utils.MapaPCB[utils.Execute.PID].Mutexs[mutexName]
		if existe {
			respuesta, err := json.Marshal("MUTEX_YA_EXISTE")
			if err != nil {
				http.Error(w, "Error al codificar mensaje como JSON", http.StatusInternalServerError)
			}
			w.WriteHeader(http.StatusOK)
			w.Write(respuesta)
		}

		// Creamos el mutex y lo agregamos al mapa de mutexs del PCB
		utils.MapaPCB[utils.Execute.PID].Mutexs[mutexName] = "LIBRE"
		respuesta, err := json.Marshal("OK")
		if err != nil {
			http.Error(w, "Error al codificar mensaje como JSON", http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(respuesta)
	}
}

// QUERY PATH - VERBO PATCH
// 3 CASOS:
// 1. Si el mutex no existe, finaliza el hilo y responde con "HILO_FINALIZADO"
// 2. Si el mutex esta libre, lo toma y responde con "MUTEX_TOMADO"
// 3. Si el mutex esta ocupado, bloquea el hilo y responde con "HILO_BLOQUEADO"
func MUTEX_LOCK(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info(fmt.Sprintf("## (%d:%d) - Solicitó syscall: MUTEX_LOCK", utils.Execute.PID, utils.Execute.TID))
		// Tomamos el valor del tid de la variable de la URL
		mutexName := r.PathValue("mutex")

		// Verificamos que el mutex exista - si NO existe mandamos el hilo a Exit
		_, existe := utils.MapaPCB[utils.Execute.PID].Mutexs[mutexName]
		if !existe {
			planificador.Finalizar_hilo(utils.Execute.TID, utils.Execute.PID, logger)
			respuesta, err := json.Marshal("HILO_FINALIZADO")
			if err != nil {
				http.Error(w, "Error al codificar mensaje como JSON", http.StatusInternalServerError)
			}
			w.WriteHeader(http.StatusOK)
			w.Write(respuesta)
			utils.Execute = nil
			utils.Planificador.Signal()
			return
		}

		// Tomamos el mutex si esta libre
		if utils.MapaPCB[utils.Execute.PID].Mutexs[mutexName] == "LIBRE" {
			utils.MapaPCB[utils.Execute.PID].Mutexs[mutexName] = strconv.Itoa(int(utils.Execute.TID))
			respuesta, err := json.Marshal("MUTEX_TOMADO")
			if err != nil {
				w.Write([]byte("Error al codificar mensaje como JSON"))
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
			w.WriteHeader(http.StatusOK)
			w.Write(respuesta)
			return
		}

		// Si no esta libre, bloqueamos el hilo
		if utils.MapaPCB[utils.Execute.PID].Mutexs[mutexName] != "LIBRE" {
			bloqueado := utils.Bloqueado{PID: utils.Execute.PID, TID: utils.Execute.TID, Motivo: utils.Mutex, QuienFue: mutexName}
			utils.Encolar(&planificador.ColaBlocked, bloqueado)
			logger.Info(fmt.Sprintf("## (%d:%d) - Bloqueado por: MUTEX", utils.Execute.PID, utils.Execute.TID))

			// Respondemos con un HILO_BLOQUEADO
			respuesta, err := json.Marshal("HILO_BLOQUEADO")
			if err != nil {
				http.Error(w, "Error al codificar mensaje como JSON", http.StatusInternalServerError)
			}
			w.WriteHeader(http.StatusOK)
			w.Write(respuesta)
			utils.Execute = nil
			utils.Planificador.Signal()
			return
		}

	}
}

// Query path - Verbo PATCH
// Si el mutex no existe responde con "HILO_FINALIZADO" y finaliza el hilo
// Si el mutex se le asigna a un hilo responde "MUTEX_ASIGNADO"
// Si el mutex queda libre responde "MUTEX_LIBRE"
// Si el hilo no posee el mutex responde "HILO_NO_POSEE_MUTEX"
func MUTEX_UNLOCK(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info(fmt.Sprintf("## (%d:%d) - Solicitó syscall: MUTEX_UNLOCK", utils.Execute.PID, utils.Execute.TID))
		// Tomamos el valor del tid de la variable de la URL
		mutexName := r.PathValue("mutex")

		// Verificamos que el mutex exista caso contrario mandamos el hilo a exit
		_, existe := utils.MapaPCB[utils.Execute.PID].Mutexs[mutexName]
		if !existe {
			planificador.Finalizar_hilo(utils.Execute.TID, utils.Execute.PID, logger)
			respuesta, err := json.Marshal("HILO_FINALIZADO")
			if err != nil {
				http.Error(w, "Error al codificar mensaje como JSON", http.StatusInternalServerError)
			}
			w.WriteHeader(http.StatusOK)
			w.Write(respuesta)
			utils.Execute = nil
			utils.Planificador.Signal()
			return
		}

		// Si el mutex existe, lo asignamos o liberamos segun corresponda
		if utils.MapaPCB[utils.Execute.PID].Mutexs[mutexName] == strconv.Itoa(int(utils.Execute.TID)) {
			for _, bloqueado := range planificador.ColaBlocked {
				// Si alguien quiere el mutex
				count := 0
				if bloqueado.PID == utils.Execute.PID && bloqueado.Motivo == utils.Mutex && bloqueado.QuienFue == mutexName {
					count++
					utils.MapaPCB[utils.Execute.PID].Mutexs[mutexName] = strconv.Itoa(int(bloqueado.TID))
					utils.Desencolar(&planificador.ColaBlocked) //! Acá creo que hay que cambiarla por la funcion de desencolar por motivo (Consultar con lucas)
					respuesta, err := json.Marshal("MUTEX_ASIGNADO")
					if err != nil {
						http.Error(w, "Error al codificar mensaje como JSON", http.StatusInternalServerError)
					}
					w.WriteHeader(http.StatusOK)
					w.Write(respuesta)
					return
				}
				// Si el mutex no lo necesita nadie
				if count == 0 {
					utils.MapaPCB[utils.Execute.PID].Mutexs[mutexName] = "LIBRE"
					respuesta, err := json.Marshal("MUTEX_LIBRE")
					if err != nil {
						http.Error(w, "Error al codificar mensaje como JSON", http.StatusInternalServerError)
					}
					w.WriteHeader(http.StatusOK)
					w.Write(respuesta)
					return
				}
			}
		}

		// Si el mutex existe y no esta tomado por el hilo q invoca la syscall
		respuesta, err := json.Marshal("HILO_NO_POSEE_MUTEX")
		if err != nil {
			http.Error(w, "Error al codificar mensaje como JSON", http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(respuesta)
	}
}

func IO(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info(fmt.Sprintf("## (%d:%d) - Solicitó syscall: IO", utils.Execute.PID, utils.Execute.TID))
		valor := r.PathValue("ms")
		tiempo, _ := strconv.Atoi(valor) //Convierto el tiempo a numero
		solicitud := utils.SolicitudIO{
			PID:       utils.Execute.PID,
			TID:       utils.Execute.TID,
			Duracion:  tiempo,
			Timestamp: time.Now(),
		}
		utils.Encolar(&planificador.ColaIO, solicitud)
		utils.Encolar(&planificador.ColaBlocked, utils.Bloqueado{PID: utils.Execute.PID, TID: utils.Execute.TID}) // Acá me falta el motivo pero no se como ponerlo
		logger.Info(fmt.Sprintf("## (%d:%d) - Bloqueado por: IO", utils.Execute.PID, utils.Execute.TID))
		utils.Execute = nil
		utils.Planificador.Signal()

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}

func Recibir_desalojo(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var magic types.HiloDesalojado
		err := decoder.Decode(&magic)
		if err != nil {
			logger.Error(fmt.Sprintf("Error al decodificar hilo desalojado: %s", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error al decodificar mensaje"))
			return
		}

		switch magic.Motivo {
		case "FIN QUANTUM":
			utils.Execute = nil
			logger.Info(fmt.Sprintf("## (%d:%d) - Desalojado por fin de Quantum", magic.PID, magic.TID))
			utils.Encolar_ColaReady(planificador.ColaReady, utils.Obtener_PCB_por_PID(magic.PID).TCBs[magic.TID])
			utils.Planificador.Signal()
		case "SEGMENTATION FAULT":
			planificador.Finalizar_proceso(magic.PID, logger)
			utils.Execute = nil
			utils.Planificador.Signal()
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}
