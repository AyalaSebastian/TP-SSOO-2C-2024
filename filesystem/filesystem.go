package main

import (
	"log/slog"
	"sync"

	"github.com/sisoputnfrba/tp-golang/cpu/utils"
	"github.com/sisoputnfrba/tp-golang/utils/logging"
)

var waitGroup sync.WaitGroup

func main() {

	// Inicio configs
	utils.Configs = utils.Iniciar_configuracion("config.json")

	// Inicio log
	log := logging.IniciarLogger("filesystem.log", utils.Configs.LogLevel)

	log.Info("Logger iniciado")
	slog.Info("Logger iniciado") //! Despues sacar este log

	// Iniciar filesystem como server
	go utils.Iniciar_cpu(log)
	waitGroup.Add(1)

	waitGroup.Wait()
}
