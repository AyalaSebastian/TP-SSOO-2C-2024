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
		client.EnviarContextoDeEjecucion(contextoEjecucion, "DUMP_MEMORY", logger)
		client.CederControlAKernell(dumpMemory, "DUMP_MEMORY", logger)
	case "IO":
		//	Informar memoria
		io := estructuraTiempo{}
		client.EnviarContextoDeEjecucion(contextoEjecucion, "IO", logger)
		client.CederControlAKernell(io, "IO", logger)
	case "PROCESS_CREATE":

		//	Informar memoria
		processCreate := estructuraProcessCreate{}
		client.EnviarContextoDeEjecucion(contextoEjecucion, "PROCESS_CREATE", logger)
		client.CederControlAKernell(processCreate, "PROCESS_CREATE", logger)
	case "THREAD_CREATE":

		//	Informar memoria
		threadCreate := estructuraThreadCreate{}
		client.EnviarContextoDeEjecucion(contextoEjecucion, "THREAD_CREATE", logger)
		client.CederControlAKernell(threadCreate, "THREAD_CREATE", logger)
	case "THREAD_JOIN":
		//	Informar memoria
		threadJoin := estructuraTid{}
		client.EnviarContextoDeEjecucion(contextoEjecucion, "THREAD_JOIN", logger)
		client.CederControlAKernell(threadJoin, "THREAD_JOIN", logger)
	case "THREAD_CANCEL":
		//	Informar memoria
		threadCancel := estructuraTid{}
		client.EnviarContextoDeEjecucion(contextoEjecucion, "THREAD_CANCEL", logger)
		client.CederControlAKernell(threadCancel, "THREAD_CANCEL", logger)
	case "MUTEX_CREATE":
		//	Informar memoria
		mutexCreate := estructuraRecurso{}
		client.EnviarContextoDeEjecucion(contextoEjecucion, "MUTEX_CREATE", logger)
		client.CederControlAKernell(mutexCreate, "MUTEX_CREATE", logger)
	case "MUTEX_LOCK":
		//	Informar memoria
		mutexLock := estructuraRecurso{}
		client.EnviarContextoDeEjecucion(contextoEjecucion, "MUTEX_LOCK", logger)
		client.CederControlAKernell(mutexLock, "MUTEX_LOCK", logger)
	case "MUTEX_UNLOCK":
		//	Informar memoria
		mutexUnlock := estructuraRecurso{}
		client.EnviarContextoDeEjecucion(contextoEjecucion, "MUTEX_UNLOCK", logger)
		client.CederControlAKernell(mutexUnlock, "MUTEX_UNLOCK", logger)
	case "THREAD_EXIT":
		//	Informar memoria
		threadExit := estructuraEmpty{}
		client.EnviarContextoDeEjecucion(contextoEjecucion, "THREAD_EXIT", logger)
		client.CederControlAKernell(threadExit, "THREAD_EXIT", logger)
	case "PROCESS_EXIT":
		//	Informar memoria
		processExit := estructuraEmpty{}
		client.EnviarContextoDeEjecucion(contextoEjecucion, "PROCESS_EXIT", logger)
		client.CederControlAKernell(processExit, "PROCESS_EXIT", logger)
	default:
		logger.Error(fmt.Sprintf("Operación desconocida: %s", operacion))

	}
}

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



	http.HandleFunc("POST /helloworld", ComunicacionMemoria)
	http.ListenAndServe(":8002", nil)

	// waitGroup.Add(1)

	// waitGroup.Wait()

	//LOGGERS OBLIGATORIOS
	//Obtencion Contexto Ejecucion
	//Logger.Info("TID: ",TID," - Solicito Contexto Ejecucion")
	//Actualizacion de Contexto de Ejecucion
	//Logger.Info("TID: ",TID," - Actualizo Contexto Ejecucion")
	//Interrupcion Recibida
	//Logger.Info("LLega Interrupcion Al Puerto Interrupt")
	//Fetch Instuccion
	//Logger.Info("TID:",TID,"- FETCH","- Program Counter: ",ProgramCounter)
	//instruccion Ejecutada
	//Logger.Info("TID: ",TID,"- Ejecutando:",Instruccion,Parametros)
	//lectura/escritura de memoria
	//.Info("TID",TID,"- Accion:",Accion,"- Direccion Fisica:",DireccionFisica)

}


// Metodo para simular el ciclo de instruccion
func CicloInstruccion() {
	var cpu *utils.CPU
	cpu.Fetch()
	cpu.Decode()
	cpu.Execute()
	cpu.CheckInterrupt()
}

// Metodo Fetch
func Fetch(cpu *utils.CPU) {
	fmt.Println("Fetch: Obtener la instruccion de la memoria")
	instruccion := cpu.memoria.ObtenerInstruccion(cpu.contexto.ProgramCounter)
	fmt.Printf("Instruccion obtenida: %s\n", instruccion)
	// Actualizar el Program Counter
	cpu.contexto.ProgramCounter++
}

func (cpu *CPU) SET(registro string, valor uint32) {
	//log obligatorio
	cpu.Logger.Printf("## TID: %d - Ejecutando: SET - %s, %d", cpu.Contexto.TID, registro, valor)

	switch registro {
	case "PC":
		cpu.Contexto.Registros.PC = valor
	case "AX":
		cpu.Contexto.Registros.AX = valor
	case "BX":
		cpu.Contexto.Registros.BX = valor
	case "CX":
		cpu.Contexto.Registros.CX = valor
	case "DX":
		cpu.Contexto.Registros.DX = valor
	case "EX":
		cpu.Contexto.Registros.EX = valor
	case "FX":
		cpu.Contexto.Registros.FX = valor
	case "GX":
		cpu.Contexto.Registros.GX = valor
	case "HX":
		cpu.Contexto.Registros.HX = valor
	case "Base":
		cpu.Contexto.Registros.Base = valor
	case "Limite":
		cpu.Contexto.Registros.Limite = valor
	}
}

func (cpu *CPU) READ_MEM(registroDatos, registroDireccion string) {
	// Implementar la logica para leer de memoria

	//log obligatorio
	cpu.Logger.Printf("## TID: %d - Acción: LEER - Dirección Física: %s", cpu.Contexto.TID, registroDireccion)
}

func (cpu *CPU) WRITE_MEM(registroDireccion, registroDatos string) {
	// Implementar la logica para escribir en memoria

	//log obligatorio
	cpu.Logger.Printf("## TID: %d - Acción: ESCRIBIR - Dirección Física: %s", cpu.Contexto.TID, registroDireccion)
}

func (cpu *CPU) SUM(registroDestino, registroOrigen string) {
	//log obligatorio
	cpu.Logger.Printf("## TID: %d - Ejecutando: SUM - %s, %s", cpu.Contexto.TID, registroDestino, registroOrigen)

	switch registroDestino {
	case "PC":
		cpu.Contexto.Registros.PC += cpu.getRegistroValor(registroOrigen)
	case "AX":
		cpu.Contexto.Registros.AX += cpu.getRegistroValor(registroOrigen)
	case "BX":
		cpu.Contexto.Registros.BX += cpu.getRegistroValor(registroOrigen)
	case "CX":
		cpu.Contexto.Registros.CX += cpu.getRegistroValor(registroOrigen)
	case "DX":
		cpu.Contexto.Registros.DX += cpu.getRegistroValor(registroOrigen)
	case "EX":
		cpu.Contexto.Registros.EX += cpu.getRegistroValor(registroOrigen)
	case "FX":
		cpu.Contexto.Registros.FX += cpu.getRegistroValor(registroOrigen)
	case "GX":
		cpu.Contexto.Registros.GX += cpu.getRegistroValor(registroOrigen)
	case "HX":
		cpu.Contexto.Registros.HX += cpu.getRegistroValor(registroOrigen)
	case "Base":
		cpu.Contexto.Registros.Base += cpu.getRegistroValor(registroOrigen)
	case "Limite":
		cpu.Contexto.Registros.Limite += cpu.getRegistroValor(registroOrigen)
	}
}

func (cpu *CPU) SUB(registroDestino, registroOrigen string) {
	//log obligatorio
	cpu.Logger.Printf("## TID: %d - Ejecutando: SUB - %s, %s", cpu.Contexto.TID, registroDestino, registroOrigen)

	switch registroDestino {
	case "PC":
		cpu.Contexto.Registros.PC -= cpu.getRegistroValor(registroOrigen)
	case "AX":
		cpu.Contexto.Registros.AX -= cpu.getRegistroValor(registroOrigen)
	case "BX":
		cpu.Contexto.Registros.BX -= cpu.getRegistroValor(registroOrigen)
	case "CX":
		cpu.Contexto.Registros.CX -= cpu.getRegistroValor(registroOrigen)
	case "DX":
		cpu.Contexto.Registros.DX -= cpu.getRegistroValor(registroOrigen)
	case "EX":
		cpu.Contexto.Registros.EX -= cpu.getRegistroValor(registroOrigen)
	case "FX":
		cpu.Contexto.Registros.FX -= cpu.getRegistroValor(registroOrigen)
	case "GX":
		cpu.Contexto.Registros.GX -= cpu.getRegistroValor(registroOrigen)
	case "HX":
		cpu.Contexto.Registros.HX -= cpu.getRegistroValor(registroOrigen)
	case "Base":
		cpu.Contexto.Registros.Base -= cpu.getRegistroValor(registroOrigen)
	case "Limite":
		cpu.Contexto.Registros.Limite -= cpu.getRegistroValor(registroOrigen)
	}
}

func (cpu *CPU) JNZ(registro string, instruccion uint32) {
	//log obligatorio
	cpu.Logger.Printf("## TID: %d - Ejecutando: JNZ - %s, %d", cpu.Contexto.TID, registro, instruccion)

	if cpu.getRegistroValor(registro) != 0 {
		cpu.Contexto.Registros.PC = instruccion
	}
}

func (cpu *CPU) LOG(registro string) {
	valor := cpu.getRegistroValor(registro)
	cpu.Logger.Printf("Valor de %s: %d", registro, valor)
}

// obtener el valor del registro pedido
func (cpu *CPU) getRegistroValor(registro string) uint32 {
	switch registro {
	case "PC":
		return cpu.Contexto.Registros.PC
	case "AX":
		return cpu.Contexto.Registros.AX
	case "BX":
		return cpu.Contexto.Registros.BX
	case "CX":
		return cpu.Contexto.Registros.CX
	case "DX":
		return cpu.Contexto.Registros.DX
	case "EX":
		return cpu.Contexto.Registros.EX
	case "FX":
		return cpu.Contexto.Registros.FX
	case "GX":
		return cpu.Contexto.Registros.GX
	case "HX":
		return cpu.Contexto.Registros.HX
	case "Base":
		return cpu.Contexto.Registros.Base
	case "Limite":
		return cpu.Contexto.Registros.Limite
	default:
		return 0
	}
}

func (cpu *CPU) execute(instruccion string, args ...interface{}) {
	switch instruccion {
	case "SET":
		cpu.SET(args[0].(string), args[1].(uint32))
	case "READ_MEM":
		cpu.READ_MEM(args[0].(string), args[1].(string))
	case "WRITE_MEM":
		cpu.WRITE_MEM(args[0].(string), args[1].(string))
	case "SUM":
		cpu.SUM(args[0].(string), args[1].(string))
	case "SUB":
		cpu.SUB(args[0].(string), args[1].(string))
	case "JNZ":
		cpu.JNZ(args[0].(string), args[1].(uint32))
	case "LOG":
		cpu.LOG(args[0].(string))
	}
}

func cargarConfiguracion(ruta string) (ContextoEjecucion, error) {
	var contexto ContextoEjecucion
	file, err := ioutil.ReadFile(ruta)
	if err != nil {
		return contexto, err
	}
	err = json.Unmarshal(file, &contexto)
	return contexto, err
}

func (cpu *CPU) checkInterrupt() {
	if atomic.LoadInt32(&cpu.InterruptFlag) != 0 {
		// Actualizar el contexto de ejecucion en la memoria
		cpu.actualizarContextoEnMemoria()

		// Devolver el TID al Kernel con motivo de la interrupcion
		cpu.devolverTIDAlKernel()
	}
}

func (cpu *CPU) actualizarContextoEnMemoria() {
	// Implementar la logica para actualizar el contexto de ejecucion en la memoria
	cpu.Logger.Println("Contexto de ejecución actualizado en la memoria")
}

func (cpu *CPU) devolverTIDAlKernel() {
	// Implementar la logica para devolver el TID al Kernel
	cpu.Logger.Printf("TID %d devuelto al Kernel con motivo de la interrupcion", cpu.TID)
}

//dentro del main

	// Cargar configuracion desde config.json
    contexto, err := cargarConfiguracion("config.json")
    if err != nil {
        log.Fatalf("Error cargando configuracion: %v", err)
    }

    // Crear una instancia de CPU
    cpu := CPU{
        Contexto:      contexto,
        Logger:        logger,
        InterruptFlag: 0,
        TID:           1, // Ejemplo de TID
    }

    // Ejemplo de ejecucion de instrucciones
    cpu.execute("SET", "AX", 10)
    cpu.execute("SUM", "AX", "BX")
    cpu.execute("LOG", "AX")

    // Chequear interrupciones
    cpu.checkInterrupt()

*/
