package cicloDeInstruccion

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/sisoputnfrba/tp-golang/cpu/client"
	"github.com/sisoputnfrba/tp-golang/cpu/cpuInstruction"
	"github.com/sisoputnfrba/tp-golang/cpu/server"
	"github.com/sisoputnfrba/tp-golang/cpu/utils"
	"github.com/sisoputnfrba/tp-golang/utils/types"
)

func Comenzar_cpu(logger *slog.Logger) {

	logger.Info(fmt.Sprintf("Obtención de Contexto de Ejecución: “## TID: %d - Solicito Contexto Ejecución”", GlobalPIDTID.TID))
	client.SolicitarContextoEjecucion(GlobalPIDTID, logger)

	for {

		// Obtener el valor actual del PC antes de Fetch
		pcAnterior := client.ReceivedContextoEjecucion.Registros.PC

		fetch()
		decode()
		execute()
		checkInterrupt()

		// 1. Fetch: obtener la próxima instrucción desde Memoria basada en el PC (Program Counter)
		err := Fetch(utils.Configs.IpMemory, utils.Configs.PortMemory, server.ReceivedPIDTID.TID, server.ReceivedPIDTID.PID, Logger)
		if err != nil {
			Logger.Error("Error en Fetch: ", slog.Any("error", err))
			break // Salimos del ciclo si hay error en Fetch
		}

		// Si no hay más instrucciones, salir del ciclo
		if Instruccion == "" {
			Logger.Info("No hay más instrucciones. Ciclo de ejecución terminado.")
			break
		}

		// 2. Decode: interpretar la instrucción obtenida
		Decode(Instruccion, Logger)

		// 3. Execute: ejecutar la instrucción decodificada (esta dentro de Decode)

		// Si el PC no fue modificado por alguna instrucción, lo incrementamos en 1
		if client.ReceivedContextoEjecucion.Registros.PC == pcAnterior {
			client.ReceivedContextoEjecucion.Registros.PC++
			logger.Info(fmt.Sprintf("PC no modificado por instrucción. Actualizado PC a: %d", client.ReceivedContextoEjecucion.Registros.PC))
		} else {
			logger.Info(fmt.Sprintf("PC modificado por instrucción a: %d", client.ReceivedContextoEjecucion.Registros.PC))
		}

	}
	logger.Info("Fin de la ejecución del CPU.")
}

//?                       VARIABLES GLOBALES                    //

// * Variable global para almacenar PID y TID
var GlobalPIDTID types.PIDTID

var AnteriorPIDTID types.PIDTID

// * Variable global para almacenar la instrucción obtenida
var Instruccion string

// * Función global que representa el estado de los registros de la CPU
var ContextoEjecucion types.ContextoEjecucion

//! /////////////////////////////////////////////////////////////////////////////
//////////////////!               FETCH                /////////////////////////
//! //////////////////////////////////////////////////////////////////////////////

// Función Fetch para obtener la próxima instrucción
func Fetch(ipMemory string, portMemory int, tid uint32, pid uint32, logger *slog.Logger) error {
	if client.ReceivedContextoEjecucion == nil {
		logger.Error("No se ha recibido el contexto de ejecución. Imposible realizar Fetch.")
		return fmt.Errorf("contexto de ejecución no disponible")
	}

	// Obtener el valor del PC (Program Counter) de la variable global
	pc := client.ReceivedContextoEjecucion.Registros.PC

	// Crear la estructura de solicitud
	requestData := struct {
		PC  uint32 `json:"pc"`
		TID uint32 `json:"tid"`
		PID uint32 `json:"pid"`
	}{PC: pc, TID: tid, PID: pid}

	// Serializar los datos en JSON
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		logger.Error("Error al codificar PC y TID a JSON: ", slog.Any("error", err))
		return err
	}

	// Crear la URL del módulo de Memoria
	url := fmt.Sprintf("http://%s:%d/instruccion", ipMemory, portMemory)

	// Crear la solicitud POST
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Error("Error al crear la solicitud: ", slog.Any("error", err))
		return err
	}

	// Establecer el encabezado de la solicitud
	req.Header.Set("Content-Type", "application/json")

	// Enviar la solicitud
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("Error al enviar la solicitud de Fetch: ", slog.Any("error", err))
		return err
	}
	defer resp.Body.Close()

	// Verificar si la respuesta fue exitosa
	if resp.StatusCode != http.StatusOK {
		logger.Error(fmt.Sprintf("Error en la respuesta de Fetch: Código de estado %d", resp.StatusCode))
		return fmt.Errorf("error en la respuesta de Fetch: Código de estado %d", resp.StatusCode)
	}

	// Decodificar la respuesta para obtener la instrucción
	var fetchedInstruction struct {
		Instruccion string `json:"instruccion"`
	}
	err = json.NewDecoder(resp.Body).Decode(&fetchedInstruction)
	if err != nil {
		logger.Error("Error al decodificar la instrucción recibida: ", slog.Any("error", err))
		return err
	}

	// Guardar la instrucción en la variable global
	Instruccion = fetchedInstruction.Instruccion

	// Log de Fetch exitoso
	logger.Info(fmt.Sprintf("Fetch Instrucción: “## TID: %d - FETCH - Program Counter: %d”", tid, pc))

	return nil
}

//! ///////////////////////////////////////////////////////////////////////////////
//! /////////////////               DECODE                /////////////////////////
//! ///////////////////////////////////////////////////////////////////////////////

func Decode(instruccion string, logger *slog.Logger) {
	logger.Info(fmt.Sprintf("Decodificando la instrucción: %s", instruccion))

	// Separar la instrucción en partes, suponiendo que esté en formato "INSTRUCCION ARGUMENTOS" ej: SET AX 5
	partes := strings.Fields(instruccion)
	if len(partes) == 0 {
		logger.Error("Instrucción vacía")
		return
	}

	operacion := partes[0] // Tipo de operación (SET, READ_MEM, etc.)
	args := partes[1:]     // Argumentos de la operación

	// Llamar a Execute para ejecutar la instrucción decodificada
	Execute(operacion, args, logger)
}

type estructuraEmpty struct {
}
type estructuraTid struct {
	tid uint32
}
type estructuraTiempo struct {
	tiempo float32
}
type estructuraRecurso struct {
	recurso string
}
type estructuraProcessCreate struct {
	archivoInstrucciones string
	tamanio              int
	PrioridadTID0        int
}
type estructuraThreadCreate struct {
	archivoInstrucciones string
	Prioridad            int
}

// Función Execute para ejecutar la instrucción decodificada
func Execute(operacion string, args []string, logger *slog.Logger) {
	switch operacion {
	case "SET":
		if len(args) != 2 {
			logger.Error("Error en argumentos de SET: se esperaban 2 argumentos")
			return
		}
		registro := args[0]
		valor, err := strconv.ParseUint(args[1], 10, 32)
		if err != nil {
			logger.Error("Error al convertir el valor para SET")
			return
		}
		// Asignar el valor al registro
		cpuInstruction.AsignarValorRegistro(registro, uint32(valor), logger)

	case "READ_MEM":
		if len(args) != 2 {
			logger.Error("Error en argumentos de READ_MEM: se esperaban 2 argumentos")
			return
		}
		registroDatos := args[0]
		registroDireccion := args[1]
		cpuInstruction.LeerMemoria(registroDatos, registroDireccion, logger)

	case "WRITE_MEM":
		if len(args) != 2 {
			logger.Error("Error en argumentos de WRITE_MEM: se esperaban 2 argumentos")
			return
		}
		registroDireccion := args[0]
		registroDatos := args[1]
		cpuInstruction.EscribirMemoria(registroDireccion, registroDatos, logger)

	case "SUM":
		if len(args) != 2 {
			logger.Error("Error en argumentos de SUM: se esperaban 2 argumentos")
			return
		}
		registroDestino := args[0]
		registroOrigen := args[1]
		cpuInstruction.SumarRegistros(registroDestino, registroOrigen, logger)

	case "SUB":
		if len(args) != 2 {
			logger.Error("Error en argumentos de SUB: se esperaban 2 argumentos")
			return
		}
		registroDestino := args[0]
		registroOrigen := args[1]
		cpuInstruction.RestarRegistros(registroDestino, registroOrigen, logger)

	case "JNZ":
		if len(args) != 2 {
			logger.Error("Error en argumentos de JNZ: se esperaban 2 argumentos")
			return
		}
		registro := args[0]
		instruccion := args[1]
		cpuInstruction.SaltarSiNoCero(registro, instruccion, logger)

	case "LOG":
		if len(args) != 1 {
			logger.Error("Error en argumentos de LOG: se esperaba 1 argumento")
			return
		}
		registro := args[0]
		cpuInstruction.LogRegistro(registro, logger)

	case "DUMP_MEMORY":

		//	Informar memoria
		dumpMemory := estructuraEmpty{}
		client.EnviarContextoDeEjecucion(ContextoEjecucion, "DUMP_MEMORY", logger)
		client.CederControlAKernell(dumpMemory, "DUMP_MEMORY", logger)
	case "IO":
		//	Informar memoria
		io := estructuraTiempo{}
		client.EnviarContextoDeEjecucion(ContextoEjecucion, "IO", logger)
		client.CederControlAKernell(io, "IO", logger)
	case "PROCESS_CREATE":

		//	Informar memoria
		processCreate := estructuraProcessCreate{}
		client.EnviarContextoDeEjecucion(ContextoEjecucion, "PROCESS_CREATE", logger)
		client.CederControlAKernell(processCreate, "PROCESS_CREATE", logger)
	case "THREAD_CREATE":

		//	Informar memoria
		threadCreate := estructuraThreadCreate{}
		client.EnviarContextoDeEjecucion(ContextoEjecucion, "THREAD_CREATE", logger)
		client.CederControlAKernell(threadCreate, "THREAD_CREATE", logger)
	case "THREAD_JOIN":
		//	Informar memoria
		threadJoin := estructuraTid{}
		client.EnviarContextoDeEjecucion(ContextoEjecucion, "THREAD_JOIN", logger)
		client.CederControlAKernell(threadJoin, "THREAD_JOIN", logger)
	case "THREAD_CANCEL":
		//	Informar memoria
		threadCancel := estructuraTid{}
		client.EnviarContextoDeEjecucion(ContextoEjecucion, "THREAD_CANCEL", logger)
		client.CederControlAKernell(threadCancel, "THREAD_CANCEL", logger)
	case "MUTEX_CREATE":
		//	Informar memoria
		mutexCreate := estructuraRecurso{}
		client.EnviarContextoDeEjecucion(ContextoEjecucion, "MUTEX_CREATE", logger)
		client.CederControlAKernell(mutexCreate, "MUTEX_CREATE", logger)
	case "MUTEX_LOCK":
		//	Informar memoria
		mutexLock := estructuraRecurso{}
		client.EnviarContextoDeEjecucion(ContextoEjecucion, "MUTEX_LOCK", logger)
		client.CederControlAKernell(mutexLock, "MUTEX_LOCK", logger)
	case "MUTEX_UNLOCK":
		//	Informar memoria
		mutexUnlock := estructuraRecurso{}
		client.EnviarContextoDeEjecucion(ContextoEjecucion, "MUTEX_UNLOCK", logger)
		client.CederControlAKernell(mutexUnlock, "MUTEX_UNLOCK", logger)
	case "THREAD_EXIT":
		//	Informar memoria
		threadExit := estructuraEmpty{}
		client.EnviarContextoDeEjecucion(ContextoEjecucion, "THREAD_EXIT", logger)
		client.CederControlAKernell(threadExit, "THREAD_EXIT", logger)
	case "PROCESS_EXIT":
		//	Informar memoria
		processExit := estructuraEmpty{}
		client.EnviarContextoDeEjecucion(ContextoEjecucion, "PROCESS_EXIT", logger)
		client.CederControlAKernell(processExit, "PROCESS_EXIT", logger)
	default:
		logger.Error(fmt.Sprintf("Operación desconocida: %s", operacion))

	}
}

//! ////////////////////////////////////////////////////////////////////////////////
//! //////////////             CHECK INTERRUPT                //////////////////////
//! ////////////////////////////////////////////////////////////////////////////////

func checkInterrupt(tid uint32, Logger *slog.Logger) {

	server.ReciboInterrupcionTID(Logger)
	if server.ReceivedInterrupt == tid {
		//actualizo contexto en memoria
		client.EnviarContextoDeEjecucion(tid, "THREAD_UPLOAD", Logger)
		//devuelvo tid a kernell con motivo de interrupcion
		Logger.Info("llega interrupcion al puerto Interrupt")
		client.DevolverTIDAlKernel(tid, Logger, "THREAD_INTERRUPT", "Interrupcion ")
		return
	}
	return
}
