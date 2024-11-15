package server

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/sisoputnfrba/tp-golang/memoria/client"
	"github.com/sisoputnfrba/tp-golang/memoria/memSistema"
	"github.com/sisoputnfrba/tp-golang/memoria/memUsuario"
	"github.com/sisoputnfrba/tp-golang/memoria/utils"
	"github.com/sisoputnfrba/tp-golang/utils/conexiones"
	"github.com/sisoputnfrba/tp-golang/utils/types"
)

func Iniciar_memoria(logger *slog.Logger) {
	mux := http.NewServeMux()

	// Endpoints
	mux.HandleFunc("POST /CREAR-PROCESO", Crear_proceso(logger))
	mux.HandleFunc("PATCH /FINALIZAR-PROCESO/{pid}", FinalizarProceso(logger))
	mux.HandleFunc("POST /CREAR_HILO", Crear_hilo(logger))
	mux.HandleFunc("POST /FINALIZAR_HILO", FinalizarHilo(logger))
	mux.HandleFunc("POST /FINALIZAR_HILO", FinalizarHilo(logger))
	mux.HandleFunc("POST /MEMORY-DUMP", MemoryDump(logger))

	// Comunicacion con CPU
	//pasa el contexto de ejecucion a cpu
	//mux.HandleFunc("POST /contexto", Obtener_Contexto_De_Ejecucion(logger))
	mux.HandleFunc("/contexto", Obtener_Contexto_De_Ejecucion(logger))

	//envia proxima instr a cpu fase fetch
	mux.HandleFunc("GET /instruccion", Obtener_Instrucción(logger))

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
		//nuevo estructura de NuevoProcesoEnMemoria
		var magic types.NuevoProcesoEnMemoria
		err := decoder.Decode(&magic)
		if err != nil {
			logger.Error(fmt.Sprintf("Error al decodificar mensaje: %s\n", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error al decodificar mensaje"))
			return
		}
		logger.Info(fmt.Sprintf("Me llegaron los siguientes parametros para crear proceso: %+v", magic))

		// Llamar a Inicializar_proceso con los parámetros correspondientes
		memUsuario.AsignarPID(magic.PCB.PID, magic.Tamanio, magic.Pseudo)

		// Si la inicialización fue exitosa
		logger.Info("## Proceso Creado - PID: %d  - Tamaño: %d", magic.PCB.PID, magic.Tamanio)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}

func FinalizarProceso(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verificar que el método sea PATCH
		if r.Method != http.MethodPatch {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}

		// Decodificar la solicitud para obtener el PID
		var pid struct {
			PID uint32 `json:"pid"`
		}
		err := json.NewDecoder(r.Body).Decode(&pid)
		if err != nil {
			http.Error(w, "Error al decodificar la solicitud", http.StatusBadRequest)
			return
		}

		//marca la particion como libre en memoria de usuario
		memUsuario.LiberarParticionPorPID(pid.PID)

		// Ejecutar la función para eliminar el contexto del PID en Memoria de sistema
		memSistema.EliminarContextoPID(pid.PID)

		///////////////////necesito que me envien el tamaño del proceso para ponerlo en el log/////////////////////
		//como no lo tengo lo saco del log

		// Log de destrucción del proceso
		logger.Info(fmt.Sprintf("Destrucción de Proceso: ## Proceso Destruido - PID: %d - ", pid.PID))

		// Responder al Kernel con "OK" si la operación fue exitosa
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	}
}

func Crear_hilo(logger *slog.Logger) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		decoder := json.NewDecoder(r.Body)
		var magic types.EnviarHiloAMemoria
		err := decoder.Decode(&magic)
		if err != nil {
			logger.Error(fmt.Sprintf("Error al decodificar mensaje: %s\n", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error al decodificar mensaje"))
			return
		}
		logger.Info(fmt.Sprintf("Me llegaron los siguientes parametros para crear proceso: %+v", magic))

		memSistema.CrearContextoTID(magic.TID, magic.PID, magic.Path)

		logger.Info("## Hilo Creado - (PID:TID) - (%d:%d)", magic.PID, magic.PID)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}

// hecho
func FinalizarHilo(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verificar que el método sea POST
		if r.Method != http.MethodPost {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}

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

		// Log de solicitud de finalización del hilo
		logger.Info(fmt.Sprintf("Finalización de hilo: ## Finalizar hilo solicitado - (PID:TID) - (%d:%d)", pidTid.PID, pidTid.TID))

		// Ejecutar la función para eliminar el contexto del TID en Memoria
		memSistema.EliminarContextoTID(pidTid.PID, pidTid.TID)

		// Log de destrucción del hilo
		logger.Info(fmt.Sprintf("Destrucción de Hilo: ## Hilo Destruido - (PID:TID) - (%d:%d)", pidTid.PID, pidTid.TID))

		// Responder al Kernel con "OK" si la operación fue exitosa
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	}
}

// Función que maneja el endpoint de Memory Dump a partir del archivo recibido por file System
func MemoryDump(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}

		// Decodificar la solicitud para obtener el PID y TID
		var pidTid struct {
			TID uint32 `json:"tid"`
			PID uint32 `json:"pid"`
		}

		err := json.NewDecoder(r.Body).Decode(&pidTid)
		if err != nil {
			logger.Error("Error al decodificar la solicitud", slog.Any("error", err))
			http.Error(w, "Error al decodificar la solicitud", http.StatusBadRequest)
			return
		}

		// Obtener la partición correspondiente al PID
		particion, existe := memUsuario.PidAParticion[pidTid.PID]
		if !existe {
			logger.Error("PID no encontrado", slog.Any("pid", pidTid.PID))
			http.Error(w, "PID no encontrado", http.StatusNotFound)
			return
		}

		// Obtener la memoria del proceso a partir de la partición
		memoriaProceso := memUsuario.MemoriaDeUsuario[memUsuario.Particiones[particion].Base : memUsuario.Particiones[particion].Base+memUsuario.Particiones[particion].Limite]

		// Generar el timestamp actual
		timestamp := time.Now().Unix()

		// Crear la estructura con la memoria, timestamp, PID y TID
		memoryDumpRequest := types.DumpFile{
			Nombre:  fmt.Sprintf("%d-%d-%d.dmp", pidTid.PID, pidTid.TID, timestamp),
			Tamanio: len(memoriaProceso),
			Datos:   memoriaProceso,
		}

		// Enviar la estructura al FileSystem para crear el archivo con el dump
		// Usamos el endpoint "dump" para enviar la estructura
		exito := client.Enviar_QueryPath(memoryDumpRequest, utils.Configs.IpFilesystem, utils.Configs.PortFilesystem, "dump", "POST", logger)
		if !exito {
			logger.Error("Error al enviar el dump al FileSystem")
			http.Error(w, "Error al enviar el dump al FileSystem", http.StatusInternalServerError)
			return
		}

		// Responder con un mensaje OK si la operación fue exitosa
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))

		// Log de éxito
		logger.Info(fmt.Sprintf("Memory Dump realizado con éxito: %d-%d-%d.dmp", pidTid.PID, pidTid.TID, timestamp))
	}
}

///////////////////////////////////////////////////////////////////////////////////////////
//////////////////////     COMUNICACION CON CPU      //////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////

// Función HTTP para obtener el contexto de ejecución completo para un PID-TID

// modificar w http.ResponseWriter, r *http.Request, y listo

func Obtener_Contexto_De_Ejecucion(logger *slog.Logger) http.HandlerFunc {
	retardoDePeticion()
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

		// Buscar el contexto para el PID en el mapa ContextosPID
		contextoPID, existePID := memSistema.ContextosPID[pidTid.PID]

		// Verificar si el PID existe
		if !existePID {
			http.Error(w, "PID no encontrado", http.StatusNotFound)
			return
		}

		// Buscar el TID dentro del contexto del PID
		contextoTID, existeTID := contextoPID.TIDs[pidTid.TID]

		// Verificar si el TID existe dentro del PID
		if !existeTID {
			http.Error(w, "TID no encontrado en el PID", http.StatusNotFound)
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
	retardoDePeticion()
	return func(w http.ResponseWriter, r *http.Request) {

		var req struct {
			ContextoDeEjecucion types.ContextoEjecucionTID
			TID                 uint32
			PID                 uint32
		}
		err := json.NewDecoder(r.Body).Decode(&req)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		memSistema.Actualizar_TID(req.PID, req.TID, req.ContextoDeEjecucion)
		logger.Info(fmt.Sprintf("## Contexto Actualizado - (PID:TID) - (%d:%d) ", req.PID, req.TID))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))

	}
}

// mal
func Obtener_Instrucción(logger *slog.Logger) http.HandlerFunc {
	retardoDePeticion()
	return func(w http.ResponseWriter, r *http.Request) {

		var requestData struct {
			PC  uint32 `json:"pc"`
			TID uint32 `json:"tid"`
			PID uint32 `json:"pid"`
		}
		err := json.NewDecoder(r.Body).Decode(&requestData)
		if err != nil {
			// Si hay error al decodificar la solicitud, enviar una respuesta con error
			http.Error(w, fmt.Sprintf("Error al leer la solicitud: %v", err), http.StatusBadRequest)
			return
		}
		logger.Info(fmt.Sprintf("## (%d:%d) - Solicitó syscall: OBTENER INSTRUCCION", requestData.PID, requestData.TID))

		instruccion := memSistema.BuscarSiguienteInstruccion(requestData.TID, requestData.PC)
		client.Enviar_QueryPath(instruccion, utils.Configs.IpCPU, utils.Configs.PortCPU, "obtener-instruccion", "GET", logger)
		return
	}
}

func Read_Mem(logger *slog.Logger) http.HandlerFunc {
	retardoDePeticion()
	return func(w http.ResponseWriter, r *http.Request) {
		// Crear una estructura para la solicitud que contiene la dirección física
		var requestData struct {
			DireccionFisica uint32 `json:"direccion_fisica"`
			TID             uint32 `json:"tid"`
			PID             uint32 `json:"pid"`
		}

		// Decodificar el cuerpo de la solicitud JSON
		err := json.NewDecoder(r.Body).Decode(&requestData)
		if err != nil {
			// Si hay error al decodificar la solicitud, enviar una respuesta con error
			http.Error(w, fmt.Sprintf("Error al leer la solicitud: %v", err), http.StatusBadRequest)
			return
		}

		// Verificar que la dirección esté dentro de los límites de memoria
		if requestData.DireccionFisica+4 > uint32(len(memUsuario.MemoriaDeUsuario)) {
			// Si hay error al buscar en la memoria, enviar una respuesta con error
			http.Error(w, fmt.Sprintf("Dirección fuera de los límites de memoria. Dirección solicitada: %d", requestData.DireccionFisica), http.StatusBadRequest)
			return
		}

		// Obtener los 4 bytes desde la dirección solicitada
		bytes := memUsuario.MemoriaDeUsuario[requestData.DireccionFisica : requestData.DireccionFisica+4]

		// Convertir los 4 bytes a uint32 (Little Endian)
		valor := binary.LittleEndian.Uint32(bytes)

		// Log obligatorio de la solicitud de lectura
		logger.Info(fmt.Sprintf("## TID: %d - Acción: LEER - Dirección Física: %d", requestData.TID, requestData.DireccionFisica))

		// Agregar log de lectura en espacio de usuario
		logger.Info(fmt.Sprintf("Escritura / lectura en espacio de usuario: “## LEER - (%d:%d) - (%d:%d) - Dir. Física: %d - Tamaño: %d”",
			requestData.PID, requestData.TID, requestData.PID, requestData.TID, requestData.DireccionFisica, 4)) // Tamaño de lectura: 4 bytes

		// Crear la respuesta JSON con el valor leído
		responseData := struct {
			Valor uint32 `json:"valor"`
		}{
			Valor: valor,
		}

		// Serializar la respuesta en JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
		if err := json.NewEncoder(w).Encode(responseData); err != nil {
			// Si hay error al serializar la respuesta, responder con error
			http.Error(w, fmt.Sprintf("Error al enviar la respuesta: %v", err), http.StatusInternalServerError)
		}
	}
}

//recibo 4 bytes y los escribo a partir del byte enviado como direccion fisica
//dentro de la Memoria de Usuario y se responderá como OK.

// falta hacer la logica con el TID recibido
// Función Write_Mem para manejar la escritura en la memoria a partir de la API

// hecha
func Write_Mem(logger *slog.Logger) http.HandlerFunc {
	retardoDePeticion()
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
		for _, particion := range memUsuario.Particiones {
			if requestData.DireccionFisica >= particion.Base && requestData.DireccionFisica < particion.Base+particion.Limite {
				encontrado = true
				break
			}
		}
		if !encontrado {
			logger.Error("Dirección física fuera de rango de particiones")
			http.Error(w, "Dirección física fuera de rango de particiones", http.StatusBadRequest)
			return
		}

		// Verificar que la dirección esté dentro de los límites de memoria
		if int(requestData.DireccionFisica+4) > len(memUsuario.MemoriaDeUsuario) {
			logger.Error("Dirección física fuera de los límites de la memoria", slog.Any("direccion_fisica", requestData.DireccionFisica))
			http.Error(w, "Dirección fuera de los límites de memoria", http.StatusBadRequest)
			return
		}

		// Escribir el valor en little-endian en la memoria
		binary.LittleEndian.PutUint32(memUsuario.MemoriaDeUsuario[requestData.DireccionFisica:], requestData.Valor)

		// Log obligatorio de Escritura en espacio de usuario
		logger.Info(fmt.Sprintf("Escritura / lectura en espacio de usuario: ## Escritura - (PID:TID) - (N/A:%d) - Dir. Física: %d - Tamaño: %d",
			requestData.TID, requestData.DireccionFisica, 4))

		// Confirmar la operación
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))

		// Log de escritura exitosa
		logger.Info(fmt.Sprintf("Escritura en memoria de usuario exitosa: TID %d - Dirección Física: %d - Valor: %d- Tamaño: %d",
			requestData.TID, requestData.DireccionFisica, requestData.Valor, 4))
	}

}

// a partir del tiempo que nos pasa el archivo configs esperamos esa cantidad en milisegundos antes de seguir con la ejecucion del proceso
func retardoDePeticion() {
	time.Sleep(time.Duration((utils.Configs.ResponseDelay * int(time.Millisecond))))
	return
}
