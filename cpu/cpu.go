package main

import (

	// "sync"
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

	// Iniciar cpu como server en un hilo para que el programa siga su ejecicion
	server.Inicializar_cpu(logger)
}
