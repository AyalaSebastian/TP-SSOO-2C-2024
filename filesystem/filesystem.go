package main

import (
	"log/slog"

	"github.com/sisoputnfrba/tp-golang/cpu/utils"
	"github.com/sisoputnfrba/tp-golang/utils/logging"
)

func main() {

	// Inicio configs
	utils.Configs = utils.Iniciar_configuracion("config.json")

	// Inicio log
	log := logging.IniciarLogger("filesystem.log", utils.Configs.LogLevel)

	log.Info("Logger iniciado")
	slog.Info("Logger iniciado")
}
