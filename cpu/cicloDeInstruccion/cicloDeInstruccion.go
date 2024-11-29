package cicloDeInstruccion

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/sisoputnfrba/tp-golang/cpu/client"
	"github.com/sisoputnfrba/tp-golang/cpu/cpuInstruction"
	"github.com/sisoputnfrba/tp-golang/cpu/utils"
	"github.com/sisoputnfrba/tp-golang/utils/types"
)

//?                       VARIABLES GLOBALES                    //

// * Variable global para almacenar PID y TID
var GlobalPIDTID types.PIDTID

var AnteriorPIDTID types.PIDTID

// * Variable global para almacenar la instrucción obtenida
var Instruccion string

// * Función global que representa el estado de los registros de la CPU
var ContextoEjecucion types.ContextoEjecucion

// * Variable global para almacenar la información de interrupción
var InterrupcionRecibida *types.InterruptionInfo

/////////////////////////////////////////////////////////////////////

func Comenzar_cpu(logger *slog.Logger) {

	logger.Info(fmt.Sprintf("Obtención de Contexto de Ejecución: ## TID: %d - Solicito Contexto Ejecución", GlobalPIDTID.TID))
	if client.SolicitarContextoEjecucion(GlobalPIDTID, logger) == nil {

		for {
			// Obtener el valor actual del PC antes de Fetch
			pcActual := client.ReceivedContextoEjecucion.Registros.PC

			if GlobalPIDTID != AnteriorPIDTID {

				// 1. Fetch: obtener la próxima instrucción desde Memoria basada en el PC (Program Counter)
				err := Fetch(GlobalPIDTID.TID, GlobalPIDTID.PID, logger)
				if err != nil {
					logger.Error("Error en Fetch: ", slog.Any("error", err))
					break // Salimos del ciclo si hay error en Fetch
				}

				// Si no hay más instrucciones, salir del ciclo
				if Instruccion == "" {
					logger.Info("No hay más instrucciones. Ciclo de ejecución terminado.")
					break
				}

				// 2. Decode: interpretar la instrucción obtenida
				Decode(Instruccion, logger)

				// 3. Execute: ejecutar la instrucción decodificada (esta dentro de Decode)

			}

			// 4. Chequear interrupciones
			CheckInterrupt(GlobalPIDTID.TID, logger)

			// Si el PC no fue modificado por alguna instrucción, lo incrementamos en 1
			if client.ReceivedContextoEjecucion.Registros.PC == pcActual {
				client.ReceivedContextoEjecucion.Registros.PC++
				logger.Info(fmt.Sprintf("PC no modificado por instrucción. Actualizado PC a: %d", client.ReceivedContextoEjecucion.Registros.PC))
			} else {
				logger.Info(fmt.Sprintf("PC modificado por instrucción a: %d", client.ReceivedContextoEjecucion.Registros.PC))
			}

		}
		logger.Info("Fin de la ejecución del CPU.")
	}
}

//! /////////////////////////////////////////////////////////////////////////////
//////////////////!               FETCH                /////////////////////////
//! //////////////////////////////////////////////////////////////////////////////

// Función Fetch para obtener la próxima instrucción
func Fetch(tid uint32, pid uint32, logger *slog.Logger) error {
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
	url := fmt.Sprintf("http://%s:%d/instruccion", utils.Configs.IpMemory, utils.Configs.PortMemory)

	// Crear la solicitud POST
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonData))
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

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error al leer el cuerpo de la respuesta:", err)
		return err
	}

	// Convertir a string
	bodyString := string(bodyBytes)

	// Guardar la instrucción en la variable global
	Instruccion = bodyString

	// Log de Fetch exitoso
	logger.Info(fmt.Sprintf("Fetch Instrucción: ## TID: %d - FETCH - Program Counter: %d", tid, pc))

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

// type estructuraThreadCreate struct {
// 	archivoInstrucciones string
// 	Prioridad            int
// }

// Función Execute para ejecutar la instrucción decodificada
func Execute(operacion string, args []string, logger *slog.Logger) {
	var proceso types.Proceso
	proceso.ContextoEjecucion = *client.ReceivedContextoEjecucion
	proceso.Pid = GlobalPIDTID.PID
	proceso.Tid = GlobalPIDTID.TID

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
		cpuInstruction.AsignarValorRegistro(registro, uint32(valor), GlobalPIDTID.TID, logger)

	case "READ_MEM":
		if len(args) != 2 {
			logger.Error("Error en argumentos de READ_MEM: se esperaban 2 argumentos")
			return
		}
		registroDatos := args[0]
		registroDireccion := args[1]
		cpuInstruction.LeerMemoria(registroDatos, registroDireccion, GlobalPIDTID, logger)

	case "WRITE_MEM":
		if len(args) != 2 {
			logger.Error("Error en argumentos de WRITE_MEM: se esperaban 2 argumentos")
			return
		}
		registroDireccion := args[0]
		registroDatos := args[1]
		cpuInstruction.EscribirMemoria(registroDireccion, registroDatos, GlobalPIDTID.TID, logger)

	case "SUM":
		if len(args) != 2 {
			logger.Error("Error en argumentos de SUM: se esperaban 2 argumentos")
			return
		}
		registroDestino := args[0]
		registroOrigen := args[1]
		cpuInstruction.SumarRegistros(registroDestino, registroOrigen, GlobalPIDTID.TID, logger)

	case "SUB":
		if len(args) != 2 {
			logger.Error("Error en argumentos de SUB: se esperaban 2 argumentos")
			return
		}
		registroDestino := args[0]
		registroOrigen := args[1]
		cpuInstruction.RestarRegistros(registroDestino, registroOrigen, GlobalPIDTID.TID, logger)

	case "JNZ":
		if len(args) != 2 {
			logger.Error("Error en argumentos de JNZ: se esperaban 2 argumentos")
			return
		}
		registro := args[0]
		instruccion := args[1]
		cpuInstruction.SaltarSiNoCero(registro, instruccion, GlobalPIDTID.TID, logger)

	case "LOG":
		if len(args) != 1 {
			logger.Error("Error en argumentos de LOG: se esperaba 1 argumento")
			return
		}
		registro := args[0]
		cpuInstruction.LogRegistro(registro, GlobalPIDTID, logger)

	case "DUMP_MEMORY":

		//	Informar memoria
		dumpMemory := estructuraEmpty{}
		client.EnviarContextoDeEjecucion(proceso, "actualizar_contexto", logger)
		logger.Info(fmt.Sprintf("## TID: %d - Actualizo Contexto Ejecución", GlobalPIDTID.TID))
		AnteriorPIDTID = GlobalPIDTID
		client.CederControlAKernell(dumpMemory, "DUMP_MEMORY", logger)

	case "IO":
		//	Informar memoria
		io := estructuraTiempo{}
		client.EnviarContextoDeEjecucion(proceso, "actualizar_contexto", logger)
		logger.Info(fmt.Sprintf("## TID: %d - Actualizo Contexto Ejecución", GlobalPIDTID.TID))
		AnteriorPIDTID = GlobalPIDTID
		client.CederControlAKernell(io, "IO", logger)

	case "PROCESS_CREATE":

		//	Informar memoria
		processCreate := estructuraProcessCreate{}
		client.EnviarContextoDeEjecucion(proceso, "actualizar_contexto", logger)
		logger.Info(fmt.Sprintf("## TID: %d - Actualizo Contexto Ejecución", GlobalPIDTID.TID))
		AnteriorPIDTID = GlobalPIDTID
		client.CederControlAKernell(processCreate, "PROCESS_CREATE", logger)

	case "THREAD_CREATE":
		// Parsear la prioridad a entero
		prio, err := strconv.Atoi(args[1])
		if err != nil {
			logger.Error("Error al convertir la prioridad para THREAD_CREATE")
			return
		}

		//	Informar memoria
		threadCreate := types.ThreadCreateParams{
			Path:      args[0],
			Prioridad: prio,
		}
		client.EnviarContextoDeEjecucion(proceso, "actualizar_contexto", logger)
		logger.Info(fmt.Sprintf("## TID: %d - Actualizo Contexto Ejecución", GlobalPIDTID.TID))
		AnteriorPIDTID = GlobalPIDTID
		client.CederControlAKernell(threadCreate, "THREAD_CREATE", logger)

	case "THREAD_JOIN":
		//	Informar memoria
		threadJoin := estructuraTid{}
		client.EnviarContextoDeEjecucion(proceso, "actualizar_contexto", logger)
		logger.Info(fmt.Sprintf("## TID: %d - Actualizo Contexto Ejecución", GlobalPIDTID.TID))
		AnteriorPIDTID = GlobalPIDTID
		client.CederControlAKernell(threadJoin, "THREAD_JOIN", logger)

	case "THREAD_CANCEL":
		//	Informar memoria
		threadCancel := estructuraTid{}
		client.EnviarContextoDeEjecucion(proceso, "actualizar_contexto", logger)
		logger.Info(fmt.Sprintf("## TID: %d - Actualizo Contexto Ejecución", GlobalPIDTID.TID))
		AnteriorPIDTID = GlobalPIDTID
		client.CederControlAKernell(threadCancel, "THREAD_CANCEL", logger)

	case "MUTEX_CREATE":
		//	Informar memoria
		mutexCreate := estructuraRecurso{}
		client.EnviarContextoDeEjecucion(proceso, "actualizar_contexto", logger)
		logger.Info(fmt.Sprintf("## TID: %d - Actualizo Contexto Ejecución", GlobalPIDTID.TID))
		AnteriorPIDTID = GlobalPIDTID
		client.CederControlAKernell(mutexCreate, "MUTEX_CREATE", logger)

	case "MUTEX_LOCK":
		//	Informar memoria
		mutexLock := estructuraRecurso{}
		client.EnviarContextoDeEjecucion(proceso, "actualizar_contexto", logger)
		logger.Info(fmt.Sprintf("## TID: %d - Actualizo Contexto Ejecución", GlobalPIDTID.TID))
		AnteriorPIDTID = GlobalPIDTID
		client.CederControlAKernell(mutexLock, "MUTEX_LOCK", logger)

	case "MUTEX_UNLOCK":
		//	Informar memoria
		mutexUnlock := estructuraRecurso{}
		client.EnviarContextoDeEjecucion(proceso, "actualizar_contexto", logger)
		logger.Info(fmt.Sprintf("## TID: %d - Actualizo Contexto Ejecución", GlobalPIDTID.TID))
		AnteriorPIDTID = GlobalPIDTID
		client.CederControlAKernell(mutexUnlock, "MUTEX_UNLOCK", logger)

	case "THREAD_EXIT":
		//	Informar memoria
		threadExit := estructuraEmpty{}
		client.EnviarContextoDeEjecucion(proceso, "actualizar_contexto", logger)
		logger.Info(fmt.Sprintf("## TID: %d - Actualizo Contexto Ejecución", GlobalPIDTID.TID))
		AnteriorPIDTID = GlobalPIDTID
		client.CederControlAKernell(threadExit, "THREAD_EXIT", logger)

	case "PROCESS_EXIT":
		//	Informar memoria
		processExit := estructuraEmpty{}
		client.EnviarContextoDeEjecucion(proceso, "actualizar_contexto", logger)
		logger.Info(fmt.Sprintf("## TID: %d - Actualizo Contexto Ejecución", GlobalPIDTID.TID))
		AnteriorPIDTID = GlobalPIDTID
		client.CederControlAKernell(processExit, "PROCESS_EXIT", logger)

	default:
		logger.Error(fmt.Sprintf("Operación desconocida: %s", operacion))

	}
}

// PARA MUCHAS INTERRUPCIONES
/*
func ChequearInterrupciones(tidActual uint32, logger *slog.Logger) {
	// Crear una nueva lista para almacenar interrupciones pendientes
	var interrupcionesPendientes []types.InterruptionInfo

	var proceso types.Proceso
	proceso.ContextoEjecucion = *client.ReceivedContextoEjecucion
	proceso.Pid = GlobalPIDTID.PID
	proceso.Tid = GlobalPIDTID.TID

	// Revisar cada interrupción en la cola global
	for _, interrupcion := range InterrupcionRecibida {
		if interrupcion.TID == tidActual {

			// Enviar el contexto de ejecución al módulo de memoria
			client.EnviarContextoDeEjecucion(proceso, "actualizar_contexto", logger)

			// Log de la interrupción recibida
			logger.Info("Interrupción Recibida: ## Llega interrupcion al puerto Interrupt", slog.Any("TID", tidActual))

			client.EnviarDesalojo(proceso.Pid, proceso.Tid, interrupcion.NombreInterrupcion, logger)

		} else {
			// Si el TID no coincide, se mantiene en la lista de interrupciones pendientes
			interrupcionesPendientes = append(interrupcionesPendientes, interrupcion)
		}
	}

	// Actualizar la cola global con las interrupciones pendientes
	InterrupcionRecibida = interrupcionesPendientes

	// Log si no hay interrupciones pendientes para el TID actual
	if len(interrupcionesPendientes) == len(InterrupcionRecibida) {
		logger.Info("No hay interrupciones pendientes para el TID actual", slog.Any("TID", tidActual))
	}
}
*/
//PARA UNA SOLA

func CheckInterrupt(tidActual uint32, logger *slog.Logger) {

	var proceso types.Proceso
	proceso.ContextoEjecucion = *client.ReceivedContextoEjecucion
	proceso.Pid = GlobalPIDTID.PID
	proceso.Tid = GlobalPIDTID.TID

	// Verificar si hay una interrupción pendiente
	if InterrupcionRecibida != nil {
		if InterrupcionRecibida.TID == tidActual {
			// Log de la interrupción recibida
			logger.Info("Interrupción Recibida: ## Llega interrupcion al puerto Interrupt", slog.Any("TID", tidActual))

			client.EnviarContextoDeEjecucion(proceso, "actualizar_contexto", logger)

			client.EnviarDesalojo(proceso.Pid, proceso.Tid, InterrupcionRecibida.NombreInterrupcion, logger)

			// Eliminar la interrupción después de procesarla
			InterrupcionRecibida = nil
		} else {
			// Si el TID no coincide, descartar la interrupción
			logger.Info("Interrupción descartada debido a TID no coincidente", slog.Any("Interrupción TID", InterrupcionRecibida.TID), slog.Any("TID actual", tidActual))

			// Descartar la interrupción al no coincidir el TID
			InterrupcionRecibida = nil
		}
	} else {
		// Log si no hay ninguna interrupción activa
		logger.Info("No hay interrupciones pendientes para el TID actual", slog.Any("TID", tidActual))
	}
}
