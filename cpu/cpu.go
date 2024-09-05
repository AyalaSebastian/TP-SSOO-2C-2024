package main

import (
	"log/slog"

	"github.com/sisoputnfrba/tp-golang/utils/logging"
)

func main() {

	// Inicio configs
	config = Iniciar_configuracion("config.json")

	// Inicio logger
	Logger := logging.IniciarLogger("cpu.log", config.LogLevel)

	// Compruebo que funcione el logger
	Logger.Info("Logger iniciado")
	slog.Info("Logger iniciado")

}
