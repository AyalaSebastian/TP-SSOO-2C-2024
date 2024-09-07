package main

import (
	"strconv"
	"sync"

	"github.com/sisoputnfrba/tp-golang/utils/client"
	"github.com/sisoputnfrba/tp-golang/utils/logging"
)

var waitGroup sync.WaitGroup

func main() {
	// Inicializamos la configuracion y el logger
	config = iniciarConfiguracion("config.json")
	logger := logging.IniciarLogger("memoria.log", config.LogLevel)

	// Solo lo puse para probar el logger, después lo sacamos
	logger.Info("Hola! Soy el módulo memoria")

	go iniciarMemoria(logger)
	waitGroup.Add(1)

	// Handshakes a kernel, cpu y fileSystem
	ipKernelParceado, err := strconv.Atoi(config.IpKernel)
	if err != nil {
		logger.Error(err.Error())

	}
	ipCPUParceado, err := strconv.Atoi(config.IpCPU)
	if err != nil {
		logger.Error(err.Error())

	}
	ipFilesystemParceado, err := strconv.Atoi(config.IpFilesystem)
	if err != nil {
		logger.Error(err.Error())
	}

	client.Enviar_handshake(strconv.Itoa(config.PortKernel), ipKernelParceado, "Estableciendo handshake con Kernel desde Memoria")
	client.Enviar_handshake(strconv.Itoa(config.PortCPU), ipCPUParceado, "Estableciendo handshake con Cpu desde Memoria")
	client.Enviar_handshake(strconv.Itoa(config.PortFilesystem), ipFilesystemParceado, "Estableciendo handshake con Filesystem desde Memoria")

	waitGroup.Wait()
}
