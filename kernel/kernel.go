package main

import (
	"github.com/sisoputnfrba/tp-golang/utils/logging"
)

func main() {
	// Inicializamos la configuracion y el logger
	config = iniciarConfiguracion("config.json")
	logger := logging.IniciarLogger("kernel.log", config.LogLevel)

	// Solo lo puse para probar el logger, después lo sacamos
	logger.Info("Hola! Soy el módulo kernel")

	iniciar_kernel(logger)

}
