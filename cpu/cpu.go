package main

import (
	"log/slog"

	"github.com/sisoputnfrba/tp-golang/cpu/utils" // Se pone esto ya que en el go.mod esta especificado asi

	"github.com/sisoputnfrba/tp-golang/utils/logging"
)

func main() {

	// Inicio configs
	utils.Configs = utils.Iniciar_configuracion("config.json")

	// Inicio logger
	Logger := logging.IniciarLogger("cpu.log", utils.Configs.LogLevel)

	// Compruebo que funcione el logger
	Logger.Info("Logger iniciado")
	slog.Info("Logger iniciado")

}
