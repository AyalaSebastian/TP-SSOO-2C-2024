package planificador

import (
	"fmt"
	"log/slog"

	"github.com/sisoputnfrba/tp-golang/kernel/client"
	"github.com/sisoputnfrba/tp-golang/kernel/utils"
	"github.com/sisoputnfrba/tp-golang/utils/generadores"
	"github.com/sisoputnfrba/tp-golang/utils/types"
)

var ColaNew []types.PCB   //Cola de procesos nuevos (Manejada por FIFO)
var ColaReady []types.TCB // Aca tengo dudas de como es, no me queda claro si las colas son distintas para PCB y TCB
var ColaBlocked []types.TCB
var ColaExit []types.TCB //Cola de procesos finalizados

func Inicializar_colas() {
	ColaNew = []types.PCB{}
	ColaReady = []types.TCB{}
	ColaBlocked = []types.TCB{}
	ColaExit = []types.TCB{}
}

// Acá para mi hay que mandar el path a memoria para que saque las instrucciones del archivo de pseudocódigo y acá mismo armar el PCB con el TCB y todo
func Crear_proceso(pseudo string, tamanio int, prioridad int, logger *slog.Logger) {
	pcb := generadores.Generar_PCB()
	utils.MapaPCB[pcb.PID] = pcb // Guardo el PCB en el mapa de PCBs3
	logger.Info(fmt.Sprintf("## (%d:0) Se crea el proceso - Estado: NEW", pcb.PID))
	if len(ColaNew) == 0 {
		// Enviar a memoria el archivo de pseudocódigo y el tamaño del proceso
		parametros := types.PathTamanio{Path: pseudo, Tamanio: tamanio}
		success := client.Enviar_Body(parametros, utils.Configs.IpMemory, utils.Configs.PortMemory, "crear-proceso", logger)
		var reg types.RegCPU

		if success {
			tcb := generadores.Generar_TCB(&pcb, prioridad, reg)
			pcb.TCBs = append(pcb.TCBs, tcb)
			utils.Encolar(&ColaReady, tcb)
			// tcb.Estado = "READY" // No se si es necesario el poner estado a los TCBs, ya que el estado va a estar dado por la cola en la que se encuentra
			logger.Info(fmt.Sprintf("## (%d:%d) Se crea el Hilo - Estado: READY", pcb.PID, tcb.TID))
		} else {
			logger.Error("No se pudo asignar espacio en memoria para el proceso")
			utils.Encolar(&ColaNew, pcb)
			// Acá para mi va un semáforo que bloquee el proceso hasta que finalice un proceso y cuando finalice se desbloquee para ejecutar nuevamente
			// Crear_Proceso(pseudo, tamanio, prioridad, logger)
		}
	}
}

// Se le pasa el pid del proceso a finalizar
func Finalizar_proceso(pid uint32, logger *slog.Logger) {

	success := client.Enviar_QueryPath(pid, utils.Configs.IpMemory, utils.Configs.PortMemory, "finalizar-proceso", logger)

	if success {
		OK := utils.Enviar_proceso_a_exit(pid, &ColaReady, &ColaBlocked, &ColaExit, logger)
		if OK {

			logger.Info(fmt.Sprintf("## Finaliza el proceso %d", pid))
		} else {
			logger.Error("Algo salió mal en Memoria al querer finalizar el proceso")
		}
	} else {
		logger.Error("Algo salió mal en Memoria al querer finalizar el proceso")
	}
}

//todo Lo que tiene que hacer la funcion
//Para la creación de hilos, el Kernel deberá informar a la Memoria y luego
//ingresarlo directamente a la cola de READY correspondiente, según su nivel de prioridad.

// Recibo de la cpu el archivo de instrucciones y la prioridad
func Crear_hilo(path string, prioridad int, logger *slog.Logger) {

	// Crear TCB
	pcb := utils.Obtener_PCB_por_PID(utils.Execute.PID)
	if pcb == nil {
		panic("No se encontro el PCB")
	}
	var reg types.RegCPU

	tcb := generadores.Generar_TCB(pcb, prioridad, reg)
	slog.Info(fmt.Sprintf("Se genero %d, %d,%d,%d", tcb.TID, tcb.Prioridad, tcb.PID, tcb.Registros.AX))

	//	Informar memoria
	infoMemoria := types.EnviarHiloAMemoria{
		TID:  tcb.TID,
		PID:  pcb.PID,
		Path: path,
	}
	if !client.Enviar_Body(infoMemoria, utils.Configs.IpMemory, utils.Configs.PortMemory, "CREAR_HILO", logger) {
		panic("Error la crear hilo")
	}

	// Ingresar a la cola de READY
	utils.Encolar(&ColaReady, tcb) //! Vamos a tener que modificar esto por el nivel de prioridad

	logger.Info(fmt.Sprintf("## (%d:%d) Se crea el Hilo - Estado: READY", pcb.PID, tcb.TID))
}

// Le envio el archivo de pseudocodigo como me viene a la memoria

//todo Lo que tiene que hacer la funcion
//Al momento de finalizar un hilo, el Kernel deberá informar a la Memoria la finalización del mismo y
//deberá mover al estado READY a todos los hilos que se encontraban bloqueados por ese TID. De esta
//manera, se desbloquean aquellos hilos bloqueados por THREAD_JOIN y por mutex tomados por el
//hilo finalizado (en caso que hubiera).

// No recibo parametros
func Finalizar_hilo() {

	// Informar a la memoria la finalizacion de un hilo

	// Quitar de la lista de los TCBs del PCB al finalizado

	// Mover al estado de ready lo que estaban bloqueados por ese TID (THREAD_JOIN y MUTEX)

}
