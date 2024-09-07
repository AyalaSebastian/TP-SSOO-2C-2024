package main

import (

	// "sync"

	"github.com/sisoputnfrba/tp-golang/memoria/utils"
	"github.com/sisoputnfrba/tp-golang/utils/client"
	"github.com/sisoputnfrba/tp-golang/utils/logging"
)

// var waitGroup sync.WaitGroup

func main() {
	// Inicializamos la configuracion y el logger
	utils.Configs = utils.Iniciar_configuracion("config.json")
	logger := logging.Iniciar_Logger("memoria.log", utils.Configs.LogLevel)

	// Solo lo puse para probar el logger, después lo sacamos
	logger.Info("Hola! Soy el módulo memoria")

	// ipFilesystemParceado, err := strconv.Atoi(utils.Configs.IpFilesystem)
	// if err != nil {
	// 	logger.Error(err.Error())
	// }
	client.Enviar_handshake(utils.Configs.IpFilesystem, utils.Configs.PortFilesystem, "Estableciendo handshake con Filesystem desde Memoria")

	utils.Iniciar_memoria(logger)
	// waitGroup.Add(1)

	// Handshakes a kernel, cpu y fileSystem
	// ipKernelParceado, err := strconv.Atoi(utils.Configs.IpKernel)
	// if err != nil {
	// 	logger.Error(err.Error())

	// }
	// ipCPUParceado, err := strconv.Atoi(utils.Configs.IpCPU)
	// if err != nil {
	// 	logger.Error(err.Error())

	// }

	// client.Enviar_handshake(strconv.Itoa(config.PortKernel), ipKernelParceado, "Estableciendo handshake con Kernel desde Memoria")
	// client.Enviar_handshake(strconv.Itoa(config.PortCPU), ipCPUParceado, "Estableciendo handshake con Cpu desde Memoria")

	// waitGroup.Wait()
}
