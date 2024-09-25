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
var ColaExit []types.PCB  //Cola de procesos finalizados

func Inicializar_colas() {
	ColaNew = []types.PCB{}
	ColaReady = []types.TCB{}
	ColaExit = []types.PCB{}
}

// Acá para mi hay que mandar el path a memoria para que saque las instrucciones del archivo de pseudocódigo y acá mismo armar el PCB con el TCB y todo
func Crear_proceso(pseudo string, tamanio int, prioridad int, logger *slog.Logger) {
	pcb := generadores.Generar_PCB()
	logger.Info(fmt.Sprintf("## (%d:0) Se crea el proceso - Estado: NEW", pcb.PID))
	if ColaNew == nil {
		// Enviar a memoria el archivo de pseudocódigo y el tamaño del proceso
		parametros := types.PathTamanio{Path: pseudo, Tamanio: tamanio}
		success := client.Enviar_Body(parametros, utils.Configs.IpMemory, utils.Configs.PortMemory, "crear-proceso", logger)

		if success {
			tcb := generadores.Generar_TCB(&pcb, prioridad)
			utils.Encolar(&ColaReady, tcb)
			logger.Info(fmt.Sprintf("## (%d:%d) Se crea el Hilo - Estado: READY", pcb.PID, tcb.TID))
		} else {
			logger.Error("No se pudo asignar espacio en memoria para el proceso")
			utils.Encolar(&ColaNew, pcb)
		}
	}
}

func Finalizar_proceso(pid uint32, logger *slog.Logger) {
	proceso := pid
	success := client.Enviar_QueryPath(proceso, utils.Configs.IpMemory, utils.Configs.PortMemory, "finalizar-proceso", logger)

	if success {
		// Aca tiene que ir la liberacion del PCB y de los TCBs asociados
		logger.Info(fmt.Sprintf("## Finaliza el proceso %d", proceso))
	} else {
		logger.Error("Algo salió mal al finalizar el proceso")
	}
}

//todo Lo que tiene que hacer la funcion
//Para la creación de hilos, el Kernel deberá informar a la Memoria y luego
//ingresarlo directamente a la cola de READY correspondiente, según su nivel de prioridad.

func Crear_hilo() {

}

//todo Lo que tiene que hacer la funcion
//Al momento de finalizar un hilo, el Kernel deberá informar a la Memoria la finalización del mismo y
//deberá mover al estado READY a todos los hilos que se encontraban bloqueados por ese TID. De esta
//manera, se desbloquean aquellos hilos bloqueados por THREAD_JOIN y por mutex tomados por el
//hilo finalizado (en caso que hubiera).

func Finalizar_hilo() {

}
