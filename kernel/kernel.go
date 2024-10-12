package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/sisoputnfrba/tp-golang/kernel/planificador"
	"github.com/sisoputnfrba/tp-golang/kernel/server"
	"github.com/sisoputnfrba/tp-golang/kernel/utils"
	"github.com/sisoputnfrba/tp-golang/utils/logging"
)

func main() {

	// Inicializamos la configuracion y el logger
	utils.Configs = utils.Iniciar_Configuracion("config.json")
	logger := logging.Iniciar_Logger("kernel.log", utils.Configs.LogLevel)

	// Inicializamos las colas de procesos
	planificador.Inicializar_colas()

	// Inicializamos el mapa de PCBs
	utils.InicializarPCBMapGlobal()

	// Obtener los parametros del primer proceso a ejecutar
	archivoPseudocodigo := os.Args[1]
	tamanioProceso, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Error: El tamaño del proceso debe ser un número entero.")
		panic(err)
	}

	// Creación del proceso inicial
	planificador.Crear_proceso(archivoPseudocodigo, tamanioProceso, 0, logger)

	// pcb1 := generadores.Generar_PCB()

	// utils.MapaPCB[pcb1.PID] = pcb1

	// tcb1 := generadores.Generar_TCB(&pcb1, 0)

	// utils.Execute.PID = pcb1.PID
	// utils.Execute.TID = tcb1.TID

	utils.Execute = &utils.ExecuteActual{PID: planificador.ColaReady[0].PID, TID: planificador.ColaReady[0].TID}

	// Inicializamos la cola de IO (esto tira un panic error y nosé por que, si se quiere probar otra cosa, comentarlo)
	go planificador.Procesar_cola_IO(&planificador.ColaIO, logger)

	// Solo para probar la funcion de procesar cola io
	solicitud := utils.SolicitudIO{PID: 1, TID: 0, Duracion: 5000, Timestamp: time.Now()}
	utils.Encolar(&planificador.ColaIO, solicitud)
	utils.Encolar(&planificador.ColaBlocked, utils.Bloqueado{PID: utils.Execute.PID, TID: utils.Execute.TID}) // Acá me falta el motivo pero no se como ponerlo
	logger.Info(fmt.Sprintf("## (%d:%d) - Bloqueado por: IO", utils.Execute.PID, utils.Execute.TID))
	utils.Execute = nil

	// Iniciamos Kernel como server
	server.Iniciar_kernel(logger)
}
