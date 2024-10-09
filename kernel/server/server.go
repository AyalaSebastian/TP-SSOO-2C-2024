package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

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
	// mux.HandleFunc("VERBO /THREAD_CREATE", planificador.Crear_hilo(logger))
	// mux.HandleFunc("VERBO /THREAD_JOIN", planificador.Crear_hilo(logger))
	// mux.HandleFunc("VERBO /THREAD_CANCEL", planificador.Finalizar_hilo(logger))
	mux.HandleFunc("POST /DUMP_MEMORY", DUMP_MEMORY(logger))

	conexiones.LevantarServidor(strconv.Itoa(utils.Configs.Port), mux, logger)

}

func PROCESS_CREATE(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var magic types.ProcessCreateParams
		err := decoder.Decode(&magic)
		if err != nil {
			logger.Error(fmt.Sprintf("Error al decodificar mensaje: %s\n", err.Error()))
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
		// pid := r.PathValue("pid")                //Recibimos el pid a finalizar
		// val, _ := strconv.ParseUint(pid, 10, 32) //Convierto el pid a uint32 ya que viene en String
		// parsePid := uint32(val)
		// Comenté lo de arriba porque el enunciado no dice que se le pasa el pid como parametro, entonces si el que hace la syscall es el hilo ejecutando, lo sacamos de la variable execute
		pid := utils.Execute.PID
		planificador.Finalizar_proceso(pid, logger)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}

func DUMP_MEMORY(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parametros := types.PIDTID{TID: utils.Execute.TID, PID: utils.Execute.PID} // Saco el pid y el tid del hilo que esta ejecutando
		success := client.Enviar_Body(parametros, utils.Configs.IpMemory, utils.Configs.PortMemory, "memory-dump", logger)
		bloqueado := utils.Bloqueado{PID: parametros.PID, TID: parametros.TID, Motivo: utils.THREAD_JOIN} // Motivo hay que cambiarlo
		utils.Encolar(&planificador.ColaBlocked, bloqueado)
		// Esto va?? logger.Info(fmt.Sprintf("## (%d:%d) - Bloqueado por: DUMP MEMORY", utils.Execute.PID, utils.Execute.TID))
		if success {

			//utils.Encolar(&planificador.ColaReady, proceso:=types.TCB{}) // Desencolo el hilo que estaba en la cola de blocked y lo paso a ready (tengo que darle algunas vueltas a esto)
		} else {
			planificador.Finalizar_proceso(utils.Execute.PID, logger)
			logger.Info(fmt.Sprintf("## Finaliza el proceso %d", parametros.PID))
		}
	}
}
