package main

import (

	// "sync"

	checkinterrupt "github.com/sisoputnfrba/tp-golang/cpu/checkInterrupt"
	"github.com/sisoputnfrba/tp-golang/cpu/cicloDeInstruccion"
	"github.com/sisoputnfrba/tp-golang/cpu/server"
	"github.com/sisoputnfrba/tp-golang/cpu/utils"
	"github.com/sisoputnfrba/tp-golang/utils/logging"
)

// var waitGroup sync.WaitGroup

func main() {

	// Inicio configs
	utils.Configs = utils.Iniciar_configuracion("config.json")
	// Inicio logger
	logger := logging.Iniciar_Logger("cpu.log", utils.Configs.LogLevel)

	// Compruebo que funcione el logger
	logger.Info("Logger iniciado")

	// Iniciar cpu como server en un hilo para que el programa siga su ejecicion
	server.Inicializar_cpu(logger)

	// Esperar hasta recibir el TID y PID
	logger.Info("Esperando TID y PID del Kernel...")
	server.Recibir_PIDTID(logger)

	// si se recibio una interrupcion mientras estoy ejecutando un proceso
	server.ReciboInterrupcionTID(logger)
	checkinterrupt.ChequearInterrupcion(cicloDeInstruccion.GlobalPIDTID.TID, logger)

}
