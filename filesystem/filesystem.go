package main

import (
	"log/slog"
	"strconv"
	"sync"

	"github.com/sisoputnfrba/tp-golang/cpu/utils"
	"github.com/sisoputnfrba/tp-golang/utils/client"
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

	// Handshakes a memoria
	ipMemoryParceado, err := strconv.Atoi(utils.Configs.IpMemory)
	if err != nil {
		log.Error(err.Error())

	}

	client.Enviar_handshake(strconv.Itoa(utils.Configs.PortMemory), ipMemoryParceado, "Estableciendo handshake con Memoria desde Filesystem")

	waitGroup.Wait()
}
