package main

import (
	"sync"

	"github.com/sisoputnfrba/tp-golang/utils/client"
	"github.com/sisoputnfrba/tp-golang/utils/logging"
)

var waitGroup sync.WaitGroup

func main() {
	// Inicializamos la configuracion y el logger
	config = iniciarConfiguracion("config.json")
	logger := logging.IniciarLogger("kernel.log", config.LogLevel)

	// Solo lo puse para probar el logger, después lo sacamos
	logger.Info("Hola! Soy el módulo kernel")

	// Iniciar cpu como server en un hilo para que el programa siga su ejecicion
	go iniciar_kernel(logger)
	waitGroup.Add(1)

	// Enviamos handshake a cpu
	client.Enviar_handshake(config.IpCPU, config.PortCPU, "Estableciendo handshake con CPU desde kernel")

	waitGroup.Wait()
}
