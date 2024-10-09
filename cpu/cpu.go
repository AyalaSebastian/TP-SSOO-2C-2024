package main

import (

	// "sync"

	"fmt"
	"os"
	"time"

	"github.com/sisoputnfrba/tp-golang/cpu/client"
	"github.com/sisoputnfrba/tp-golang/cpu/server"
	"github.com/sisoputnfrba/tp-golang/cpu/utils"
	"github.com/sisoputnfrba/tp-golang/utils/logging"
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
	Logger.Info("Esperando Contexto de ejecucion a Memoria...")
	err := client.SolicitarContextoEjecucion(utils.Configs.IpMemory, utils.Configs.PortMemory, server.ReceivedPIDTID.PID, server.ReceivedPIDTID.TID, Logger)
	if err != nil {
		Logger.Error("No se pudo obtener el contexto de ejecución: ", err)
		os.Exit(1) // Salir si hay un error
	}

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

type CPU struct {
	Contexto          ContextoEjecucion `json:"contexto"`
	MMU               MMU               `json:"mmu"`
	Memoria           Memoria           `json:"memoria"`
	InstruccionActual string            `json:"instruccion_actual"`
	Logger            *log.Logger
	InterruptFlag     int32 // Flag para interrupciones
	TID               int   // Identificador del TID actual
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
