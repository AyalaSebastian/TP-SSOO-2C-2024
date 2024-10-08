package main

import (
	"strconv"

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
	// archivoPseudocodigo := os.Args[1]
	// tamanioProceso, err := strconv.Atoi(os.Args[2])
	// if err != nil {
	// 	fmt.Println("Error: El tamaño del proceso debe ser un número entero.")
	// 	panic(err)
	// }

	// // Creación del proceso inicial
	// planificador.Crear_proceso(archivoPseudocodigo, tamanioProceso, 0, logger)

	// // Solo para probar la funcion de finalizar proceso, esto no va aca
	// planificador.Finalizar_proceso(1, logger)

	//todo condiciones para probar crear hilo
	// pcb := generadores.Generar_PCB()
	// pcb.PID = 0
	// utils.MapaPCB[pcb.PID] = pcb

	// utils.Execute.PID = 0
	// utils.Execute.TID = 0

	// planificador.Crear_hilo("hola", 0, logger)

	//todo -- Condiciones para probar finalizar hilo
	// tener un proceso con 3 hilos
	// Los dos procesos que no se van a finalizar tienen que estar bloqueados uno por THREAD_JOIN y el otro por MUTEX
	pcb := generadores.Generar_PCB()

	tcb1 := generadores.Generar_TCB(&pcb, 0)
	tcb2 := generadores.Generar_TCB(&pcb, 0)
	tcb3 := generadores.Generar_TCB(&pcb, 0)

	tidPCB1Parseado := strconv.Itoa(int(tcb1.TID))

	pcb.Mutexs["Mutex1"] = tidPCB1Parseado

	bloqueadoTCB2 := utils.Bloqueado{
		PID:      tcb2.PID,
		TID:      tcb2.TID,
		Motivo:   utils.THREAD_JOIN,
		QuienFue: tidPCB1Parseado,
	}

	bloqueadoTCB3 := utils.Bloqueado{
		PID:      tcb3.PID,
		TID:      tcb3.TID,
		Motivo:   utils.Mutex,
		QuienFue: "Mutex1",
	}

	planificador.Finalizar_hilo(tcb1.TID, tcb1.PID, logger)

	logger.Info("para q no rompan los huevos %d, %d", bloqueadoTCB2.PID, bloqueadoTCB3.TID)

	// Iniciamos Kernel como server
	server.Iniciar_kernel(logger)
}
