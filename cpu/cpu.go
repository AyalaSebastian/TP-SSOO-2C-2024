package main

import (

	// "sync"
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sisoputnfrba/tp-golang/cpu/client"
	"github.com/sisoputnfrba/tp-golang/cpu/cpuInstruction"
	"github.com/sisoputnfrba/tp-golang/cpu/server"
	"github.com/sisoputnfrba/tp-golang/cpu/utils"
	"github.com/sisoputnfrba/tp-golang/utils/logging"
	"github.com/sisoputnfrba/tp-golang/utils/types"
)

// var waitGroup sync.WaitGroup

func main() {

	// Inicio configs
	utils.Configs = utils.Iniciar_configuracion("config.json")
	// Inicio logger
	Logger := logging.Iniciar_Logger("cpu.log", utils.Configs.LogLevel)

	// Compruebo que funcione el logger
	Logger.Info("Logger iniciado")

	// Iniciar cpu como server en un hilo para que el programa siga su ejecicion
	server.Iniciar_cpu(Logger)

	// Esperar hasta recibir el TID y PID
	Logger.Info("Esperando TID y PID del Kernel...")

	/*   OTRA OPCION
	// Comprobar que ReceivedPIDTID no sea nil
	if server.ReceivedPIDTID == nil {
		Logger.Error("No se recibió TID y PID del módulo kernel")
		os.Exit(1) // Salir si no hay TID y PID
	}
	*/

	// Bucle de espera hasta que receivedTCB esté asignado
	for server.ReceivedPIDTID == nil { // Revisamos el estado de la variable desde el paquete server
		time.Sleep(100 * time.Millisecond) // Pausa pequeña para evitar que el bucle sea intensivo
	}

	// Una vez recibido el TID y PID, continuamos la ejecución
	pidtid := server.ReceivedPIDTID
	Logger.Info(fmt.Sprintf("Recibido TID: %d, PID: %d", pidtid.TID, pidtid.PID))

	// Pido el contexto de ejecucion a memoria
	// Llamar a la función SolicitarContextoEjecucion con el log obligatorio
	Logger.Info(fmt.Sprintf("Obtención de Contexto de Ejecución: “## TID: %d - Solicito Contexto Ejecución”", server.ReceivedPIDTID.TID))
	err := client.SolicitarContextoEjecucion(utils.Configs.IpMemory, utils.Configs.PortMemory, server.ReceivedPIDTID.PID, server.ReceivedPIDTID.TID, Logger)
	if err != nil {
		Logger.Error(fmt.Sprintf("No se pudo obtener el contexto de ejecución: %v", err))
		os.Exit(1) // Salir si hay un error
	}

	// Ciclo principal de ejecución de la CPU
	for {

		// Obtener el valor actual del PC antes de Fetch
		pcAnterior := client.ReceivedContextoEjecucion.Registros.PC

		// 1. Fetch: obtener la próxima instrucción desde Memoria basada en el PC (Program Counter)
		err := Fetch(utils.Configs.IpMemory, utils.Configs.PortMemory, server.ReceivedPIDTID.TID, Logger)
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
			Logger.Info(fmt.Sprintf("PC no modificado por instrucción. Actualizado PC a: %d", client.ReceivedContextoEjecucion.Registros.PC))
		} else {
			Logger.Info(fmt.Sprintf("PC modificado por instrucción a: %d", client.ReceivedContextoEjecucion.Registros.PC))
		}

	}

	Logger.Info("Fin de la ejecución del CPU.")

}

////////////////////////////////////////////////////////////////////////////////
///////////////////               FETCH                /////////////////////////
////////////////////////////////////////////////////////////////////////////////

// Variable global para almacenar la instrucción obtenida
var Instruccion string

// Función Fetch para obtener la próxima instrucción
func Fetch(ipMemory string, portMemory int, tid uint32, logger *slog.Logger) error {
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
	}{PC: pc, TID: tid}

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

/////////////////////////////////////////////////////////////////////////////////
///////////////////               DECODE                /////////////////////////
/////////////////////////////////////////////////////////////////////////////////

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

//////////////////////////////////////////////////////////////////////////////////
///////////////////               EXECUTE                /////////////////////////
//////////////////////////////////////////////////////////////////////////////////

// Función global que representa el estado de los registros de la CPU
var contextoEjecucion types.ContextoEjecucion

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

	default:
		logger.Error(fmt.Sprintf("Operación desconocida: %s", operacion))
	}
}

func checkInterrupt(tid uint32, Logger slog.Logger) {

	server.ReciboInterrupcionTID(&Logger)
	if server.ReceivedInterrupt == tid {
		//actualizo contexto en memoria
		client.ActualizarContextoDeEjecucion(tid, Logger)
		//devuelvo tid a kernell con motivo de interrupcion
		Logger.Info("llega interrupcion al puerto Interrupt")
		client.DevolverTIDAlKernel(tid, Logger)
		return
	}
	return
}

/*

	//inicio cpu

	//esperando a recibir TIP y PID de kernel

	//solicitarle el contexto de ejecucion a memoria para poder inciar la ejecucion

	//recibir el contexto de ejecucion de memoria

	//iniciar la ejecucion

	//para instrucciones que interactuan tendra que traducir
	//las direcciones lógicas (propias del proceso)
	//a direcciones físicas (propias de la memoria).
	//Para ello simulará la existencia de una MMU.

	//Durante el transcurso de la ejecución de un HILO,
	//se irá actualizando su Contexto de Ejecución
	//donde se informará a la Memoria bajo los siguientes escenarios:
	//finalización del mismo (PROCESS_EXIT o THREAD_EXIT),
	//ejecutar una llamada al Kernel (syscall),
	//deber ser desalojado (interrupción) o por la ocurrencia de un error Segmentation Fault.

	//como hacer una Syscall

	//donde hay Lectura/Escritura de Memoria

//fetch: le pido a memoria la instruccion que sigue que esta en el program counter(registro PC)

	// en READ_MEM le dice a memoria leeme esto o solo lo pasa a registro datos
	//y para que lo lee?

	//Al finalizar el ciclo, el PC deberá ser actualizado sumándole 1
	//en caso de que éste no haya sido modificado por la instrucción.
	//habla solo de JNZ ?

	//Check Interrupt: si kernel nos envia el TID se
	//actualiza el contexto de ejecucion en memoria
	//y	devolver el TID al Kernel con ¿motivo de la interrupcion?

*/
