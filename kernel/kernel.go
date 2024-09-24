package main

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/sisoputnfrba/tp-golang/kernel/client"
	"github.com/sisoputnfrba/tp-golang/kernel/utils"

	// "github.com/sisoputnfrba/tp-golang/utils/client"
	"github.com/sisoputnfrba/tp-golang/utils/generadores"
	"github.com/sisoputnfrba/tp-golang/utils/logging"
	"github.com/sisoputnfrba/tp-golang/utils/types"
)

// Inicializamos las colas de procesos
var colaNew []types.PCB //Cola de procesos nuevos (Manejada por FIFO)

func main() {

	// Inicializamos la configuracion y el logger
	utils.Configs = utils.Iniciar_Configuracion("config.json")
	logger := logging.Iniciar_Logger("kernel.log", utils.Configs.LogLevel)

	// Obtener los parametros del primer proceso a ejecutar
	archivoPseudocodigo := os.Args[1]
	tamanioProceso, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Error: El tamaño del proceso debe ser un número entero.")
		panic(err)
	}

	// Creación del proceso inicial
	planificar_proceso_inicial(archivoPseudocodigo, tamanioProceso, logger)

	// Iniciamos Kernel como server
	utils.Iniciar_kernel(logger)

}

// Acá para mi hay que mandar el path a memoria para que saque las instrucciones del archivo de pseudocódigo y acá mismo armar el PCB con el TCB y todo
func planificar_proceso_inicial(pseudo string, tamanio int, logger *slog.Logger) {
	pcb := generadores.Generar_PCB()
	logger.Info(fmt.Sprintf("## (<PID>:%d) Se crea el proceso - Estado: NEW", pcb.PID))
	if colaNew == nil {
		// Enviar a memoria el archivo de pseudocódigo y el tamaño del proceso
		client.Enviar_path_tamanio(utils.Configs.IpMemory, utils.Configs.PortMemory, pseudo, tamanio, logger)
	}
	utils.Encolar(&colaNew, pcb)
}

// tcb := generadores.Generar_TCB(&pcb, 0)
