package main

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/sisoputnfrba/tp-golang/kernel/client"
	"github.com/sisoputnfrba/tp-golang/kernel/server"
	"github.com/sisoputnfrba/tp-golang/kernel/utils"
	"github.com/sisoputnfrba/tp-golang/utils/generadores"
	"github.com/sisoputnfrba/tp-golang/utils/logging"
	"github.com/sisoputnfrba/tp-golang/utils/types"
)

// Inicializamos las colas de procesos
var colaNew []types.PCB   //Cola de procesos nuevos (Manejada por FIFO)
var colaReady []types.TCB // Aca tengo dudas de como es, no me queda claro si las colas son distintas para PCB y TCB

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
	crear_proceso(archivoPseudocodigo, tamanioProceso, 0, logger)

	// Iniciamos Kernel como server
	server.Iniciar_kernel(logger)

}

// Acá para mi hay que mandar el path a memoria para que saque las instrucciones del archivo de pseudocódigo y acá mismo armar el PCB con el TCB y todo
func crear_proceso(pseudo string, tamanio int, prioridad int, logger *slog.Logger) {
	pcb := generadores.Generar_PCB()
	logger.Info(fmt.Sprintf("## (%d:0) Se crea el proceso - Estado: NEW", pcb.PID))
	if colaNew == nil {
		// Enviar a memoria el archivo de pseudocódigo y el tamaño del proceso
		success := client.Enviar_parametros_proceso(utils.Configs.IpMemory, utils.Configs.PortMemory, pseudo, tamanio, logger)

		if success {
			tcb := generadores.Generar_TCB(&pcb, prioridad)
			utils.Encolar(&colaReady, tcb)
			logger.Info(fmt.Sprintf("## (%d:%d) Se crea el Hilo - Estado: READY", pcb.PID, tcb.TID))
		} else {
			logger.Error("No se pudo asignar espacio en memoria para el proceso")
			utils.Encolar(&colaNew, pcb)
		}
	}
}

// tcb := generadores.Generar_TCB(&pcb, 0)
