package main

import (
	"log/slog"
	"strconv"
	"sync"

	"github.com/sisoputnfrba/tp-golang/cpu/utils" // Se pone esto ya que en el go.mod esta especificado asi
	"github.com/sisoputnfrba/tp-golang/utils/client"
	"github.com/sisoputnfrba/tp-golang/utils/logging"
)

var waitGroup sync.WaitGroup

func main() {

	// Inicio configs
	utils.Configs = utils.Iniciar_configuracion("config.json")

	// Inicio logger
	Logger := logging.IniciarLogger("cpu.log", utils.Configs.LogLevel)

	// Compruebo que funcione el logger
	Logger.Info("Logger iniciado")
	slog.Info("Logger iniciado")

	// Iniciar cpu como server en un hilo para que el programa siga su ejecicion
	go utils.Iniciar_cpu(Logger)
	waitGroup.Add(1)

	// Handshakes a kernel y memoria
	ipKernelParceado, err := strconv.Atoi(utils.Configs.IpKernel)
	if err != nil {
		Logger.Error(err.Error())
	}

	ipMemoryParceado, err := strconv.Atoi(utils.Configs.IpMemory)
	if err != nil {
		Logger.Error(err.Error())
	}

	client.Enviar_handshake(strconv.Itoa(utils.Configs.PortKernel), ipKernelParceado, "Estableciendo handshake con Kernel desde CPU")
	client.Enviar_handshake(strconv.Itoa(utils.Configs.PortMemory), ipMemoryParceado, "Estableciendo handshake con Memoria desde CPU")

	waitGroup.Wait()
}
