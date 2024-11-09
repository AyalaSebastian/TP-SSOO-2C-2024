package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/sisoputnfrba/tp-golang/kernel/utils"
	"github.com/sisoputnfrba/tp-golang/memoria/memSistema"
	"github.com/sisoputnfrba/tp-golang/memoria/memUsuario"
	"github.com/sisoputnfrba/tp-golang/utils/conexiones"
	"github.com/sisoputnfrba/tp-golang/utils/types"
)

func Iniciar_memoria(logger *slog.Logger) {
	mux := http.NewServeMux()

	// Endpoints
	mux.HandleFunc("POST /CREAR-PROCESO", Crear_proceso(logger))
	mux.HandleFunc("PATCH /FINALIZAR-PROCESO/{pid}", Finalizar_proceso(logger))
	mux.HandleFunc("POST /CREAR_HILO", Crear_hilo(logger))
	mux.HandleFunc("POST /FINALIZAR_HILO", Finalizar_hilo(logger))
	mux.HandleFunc("POST /MEMORY-DUMP", Memory_dump(logger))

	// Comunicacion con CPU
	//pasa el contexto de ejecucion a cpu
	//mux.HandleFunc("POST /contexto", Obtener_Contexto_De_Ejecucion(logger))
	mux.HandleFunc("/contexto", Obtener_Contexto_De_Ejecucion(logger))

	//envia proxima instr a cpu fase fetch
	mux.HandleFunc("GET /instruccion /{tid}/{pc}", Obtener_Instrucción(logger))

	//recibo msj de cpu para que haga la instruccion read mem
	mux.HandleFunc("/read_mem / {direccionFisica}", Read_Mem(logger))

	//recibo msj de cpu para que haga la instruccion write mem
	mux.HandleFunc("POST /write_mem", Write_Mem(logger))

	conexiones.LevantarServidor(strconv.Itoa(utils.Configs.Port), mux, logger)

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
		memUsuario.CrearProceso(utils.Execute.TID)
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
		var err error
		var req struct {
		}
		err = json.NewDecoder(r.Body).Decode(&req)

		if err != nil {
			logger.Error(fmt.Sprintf("Error al decodificar mensaje: %s\n", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error al decodificar mensaje"))
			return
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}

// Comunicacion con CPU

//chat gpt me dice que funcionan bien "Obtener_Contexto_De_Ejecucion" de mem y "SolicitarContextoEjecucion" de cpu

// Función HTTP para obtener el contexto de ejecución completo para un PID-TID

// modificar w http.ResponseWriter, r *http.Request, y listo
func Obtener_Contexto_De_Ejecucion(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Decodificar la solicitud para obtener el PID y TID
		var pidTid struct {
			PID uint32 `json:"pid"`
			TID uint32 `json:"tid"`
		}

		err := json.NewDecoder(r.Body).Decode(&pidTid)
		if err != nil {
			http.Error(w, "Error al decodificar la solicitud", http.StatusBadRequest)
			return
		}

		// Buscar los contextos para el PID y TID en los mapas
		contextoPID, existePID := memSistema.ContextosPID[int(pidTid.PID)]
		contextoTID, existeTID := memSistema.ContextosTID[int(pidTid.TID)]

		// Verificar que ambos contextos existen en los mapas
		if !existePID || !existeTID {
			http.Error(w, "PID o TID no encontrado", http.StatusNotFound)
			return
		}

		// Log de solicitud de contexto OBLIGATORIO
		fmt.Printf("Solicitud / actualización de Contexto: “## Contexto Solicitado - (PID:TID) - (%d:%d)”\n", pidTid.PID, pidTid.TID)

		// Crear el contexto completo usando la estructura que CPU espera (RegCPU)
		contextoCompleto := types.RegCPU{
			PC:     contextoTID.PC,
			AX:     contextoTID.AX,
			BX:     contextoTID.BX,
			CX:     contextoTID.CX,
			DX:     contextoTID.DX,
			EX:     contextoTID.EX,
			FX:     contextoTID.FX,
			GX:     contextoTID.GX,
			HX:     contextoTID.HX,
			Base:   contextoPID.Base,
			Limite: contextoPID.Limite,
		}

		// Codificar el contexto completo como JSON y enviarlo como respuesta
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(contextoCompleto)
		if err != nil {
			http.Error(w, "Error al codificar la respuesta", http.StatusInternalServerError)
			return
		}

		fmt.Printf("Contexto completo enviado para PID %d y TID %d\n", pidTid.PID, pidTid.TID)
	}
}

func Actualizar_Contexto(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			TID uint32
			PID uint32
		}
		err := json.NewDecoder(r.Body).Decode(&req)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		//	Actualizar_Contexto(req.PID, req.TID)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))

	}
}

func Obtener_Instrucción(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info(fmt.Sprintf("## (%d:%d) - Solicitó syscall: OBTENER INSTRUCCION", utils.Execute.PID, utils.Execute.TID))
		
		tid := r.PathValue("tid")
		pc := r.PathValue("pc")
		instruccion := memSistema.BuscarSiguienteInstruccion(tid, pc)
		client.Enviar_QueryPath(instruccion, utils.Configs.IpCPU, utils.Config.PortCPU, "obtener-instruccion", "GET", logger)
		return
	}
}

func Read_Mem(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//	direccion_fisica := r.PathValue("direccion_fisica")
		//	valor := memUsuario.BuscarPorDireccion(direccion_fisica)
		//	client.Enviar_QueryPath(valor, utils.Configs.IpCPU, utils.Config.PortCPU, "readMem", "GET", logger)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))

	}
}

//recibo 4 bytes y los escribo a partir del byte enviado como direccion fisica
//dentro de la Memoria de Usuario y se responderá como OK.

// falta hacer la logica con el TID recibido
// Función Write_Mem para manejar la escritura en la memoria a partir de la API
func Write_Mem(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestData struct {
			DireccionFisica uint32 `json:"direccion_fisica"`
			Valor           uint32 `json:"valor"`
			TID             uint32 `json:"tid"`
		}

		// Decodificar el JSON
		err := json.NewDecoder(r.Body).Decode(&requestData)
		if err != nil {
			logger.Error("Error al decodificar JSON en Write_Mem", slog.Any("error", err))
			http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
			return
		}

		// Verificar si la dirección física está dentro de alguna partición
		encontrado := false
		/*	for _, particion := range memSistema.Particiones {
				if requestData.DireccionFisica >= particion.Base && requestData.DireccionFisica < particion.Base+particion.Limite-4 {
					encontrado = true
					break
				}
			}
		*/
		if !encontrado {
			logger.Error("Dirección física fuera de rango de particiones")
			http.Error(w, "Dirección física fuera de rango de particiones", http.StatusBadRequest)
			return
		}

		// Escribir el valor en little-endian en la memoria
		//	binary.LittleEndian.PutUint32(memSistema.Memoria[requestData.DireccionFisica:], requestData.Valor)

		// Log obligatorio de Escritura en espacio de usuario
		logger.Info(fmt.Sprintf("Escritura / lectura en espacio de usuario: ## Escritura - (PID:TID) - (N/A:%d) - Dir. Física: %d - Tamaño: %d",
			requestData.TID, requestData.DireccionFisica, 4))

		// Confirmar la operación
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))

		// Log de escritura exitosa
		logger.Info(fmt.Sprintf("Escritura en memoria exitosa: TID %d - Dirección Física: %d - Valor: %d",
			requestData.TID, requestData.DireccionFisica, requestData.Valor))
	}
}
