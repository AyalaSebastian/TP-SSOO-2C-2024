package main

import (
	"log/slog"

	"github.com/sisoputnfrba/tp-golang/utils/logging"
)

func main() {

	// Inicio configs
	config = Iniciar_configuracion("config.json")

	// Inicio log
	log := logging.IniciarLogger("filesystem.log", config.LogLevel)

	log.Info("Logger iniciado")
	slog.Info("Logger iniciado")
}
