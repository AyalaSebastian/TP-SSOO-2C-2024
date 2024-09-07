package main

import (
	// "strconv"
	// "sync"

	"github.com/sisoputnfrba/tp-golang/cpu/utils"
	// "github.com/sisoputnfrba/tp-golang/utils/client"
	"github.com/sisoputnfrba/tp-golang/utils/logging"
)

// var waitGroup sync.WaitGroup

func main() {

	// Inicio configs
	utils.Configs = utils.Iniciar_configuracion("config.json")

	// Inicio log
	logger := logging.Iniciar_Logger("filesystem.log", utils.Configs.LogLevel)

	logger.Info("Logger iniciado")

	// Iniciar filesystem como server
	utils.Iniciar_cpu(logger)
	// waitGroup.Add(1)

	// Handshakes a memoria
	// 	ipMemoryParceado, err := strconv.Atoi(utils.Configs.IpMemory)
	// 	if err != nil {
	// 		logger.Error(err.Error())

	// 	}

	// 	client.Enviar_handshake(strconv.Itoa(utils.Configs.PortMemory), ipMemoryParceado, "Estableciendo handshake con Memoria desde Filesystem")

	// 	// waitGroup.Wait()
}
