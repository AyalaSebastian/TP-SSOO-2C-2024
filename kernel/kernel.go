package main

import (
	"github.com/sisoputnfrba/tp-golang/kernel/planificador"
	"github.com/sisoputnfrba/tp-golang/kernel/server"
	"github.com/sisoputnfrba/tp-golang/kernel/utils"
	"github.com/sisoputnfrba/tp-golang/utils/generadores"
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

	//! DESCOMENTAR
	// archivoPseudocodigo := os.Args[1]
	// tamanioProceso, err := strconv.Atoi(os.Args[2])
	// if err != nil {
	// 	fmt.Println("Error: El tamaño del proceso debe ser un número entero.")
	// 	panic(err)
	// }

	// // Creación del proceso inicial
	// planificador.Crear_proceso(archivoPseudocodigo, tamanioProceso, 0, logger)

	// Solo para probar la funcion de finalizar proceso, esto no va aca
	// planificador.Finalizar_proceso(1, logger)
	//! DESCOMENTAR

	pcb1 := generadores.Generar_PCB()

	utils.MapaPCB[pcb1.PID] = pcb1

	tcb1 := generadores.Generar_TCB(&pcb1, 0)

	utils.Execute.PID = pcb1.PID
	utils.Execute.TID = tcb1.TID

	// Iniciamos Kernel como server
	server.Iniciar_kernel(logger)
}
