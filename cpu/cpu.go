package main

import (

	// "sync"

	"github.com/sisoputnfrba/tp-golang/cpu/utils" // Se pone esto ya que en el go.mod esta especificado asi
//	"github.com/sisoputnfrba/tp-golang/utils/client"
	"github.com/sisoputnfrba/tp-golang/utils/logging"
)

// var waitGroup sync.WaitGroup

func main() {
	TID := 0
	// Inicio configs
	utils.Configs = utils.Iniciar_configuracion("config.json")

	// Inicio logger
	Logger := logging.Iniciar_Logger("cpu.log", utils.Configs.LogLevel)

	// Compruebo que funcione el logger
	Logger.Info("Logger iniciado")

	//Obtencion Contexto Ejecucion
	Logger.Info("TID: ",TID," - Solicito Contexto Ejecucion")
	//Actualización de Contexto de Ejecución
	Logger.Info("TID: ",TID," - Actualizo Contexto Ejecucion")
	//Interrupcion Recibida
	Logger.Info("LLega Interrupcion Al Puerto Interrupt")
	//Fetch Instuccion
	//Logger.Info("TID:",TID,"- FETCH","- Program Counter: ",ProgramCounter)
	//instruccion Ejecutada
	//Logger.Info("TID: ",TID,"- Ejecutando:",Instruccion,Parametros)
	//lectura/escritura de memoria
	//.Info("TID",TID,"- Accion:",Accion,"- Direccion Fisica:",DireccionFisica)
	// "Handshake" a memoria
	//Enviar_handshake(utils.Configs.IpMemory, utils.Configs.PortMemory, "Estableciendo handshake con Memoria desde CPU")

	// Iniciar cpu como server en un hilo para que el programa siga su ejecicion
	utils.Iniciar_cpu(Logger)
	// waitGroup.Add(1)

	// waitGroup.Wait()
}
