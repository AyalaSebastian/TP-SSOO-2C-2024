package main

import (
	// "sync"

	"github.com/sisoputnfrba/tp-golang/kernel/utils"
	"github.com/sisoputnfrba/tp-golang/utils/client"
	"github.com/sisoputnfrba/tp-golang/utils/logging"
)

// var waitGroup sync.WaitGroup

func main() {
	// Inicializamos la configuracion y el logger
	utils.Configs = utils.Iniciar_Configuracion("config.json")
	logger := logging.Iniciar_Logger("kernel.log", utils.Configs.LogLevel)

	// Solo lo puse para probar el logger, después lo sacamos
	logger.Info("Hola! Soy el módulo kernel")

	// Enviamos handshake a cpu
	client.Enviar_handshake(utils.Configs.IpCPU, utils.Configs.PortCPU, "Estableciendo handshake con CPU desde Kernel")
	client.Enviar_handshake(utils.Configs.IpMemory, utils.Configs.PortMemory, "Estableciendo handshake con Memoria desde Kernel")

	// Iniciamos Kernel como server
	utils.Iniciar_kernel(logger)
	// waitGroup.Add(1)

	// waitGroup.Wait()
}
