package main

import (
	"sync"

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

	waitGroup.Wait()
}
