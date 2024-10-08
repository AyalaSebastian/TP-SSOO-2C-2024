package utils

import (
	"fmt"
	"log/slog"

	"github.com/sisoputnfrba/tp-golang/utils/types"
)

// Hilo ejecutando actualmente
type ExecuteActual struct {
	PID uint32 `json:"pid"`
	TID uint32 `json:"tid"`
}

var Execute ExecuteActual

// Mapa para almacenar los PCB con su PID como clave
var MapaPCB map[uint32]types.PCB

// Inicializa el mapa de PCBs
func InicializarPCBMapGlobal() {
	MapaPCB = make(map[uint32]types.PCB)
}

// Función para obtener el PCB a partir de un PID
func Obtener_PCB_por_PID(pid uint32) *types.PCB {
	pcb, existe := MapaPCB[pid]
	if !existe {
		return nil
	}
	return &pcb
}

func Eliminar_TCBs_de_cola(pcb *types.PCB, cola *[]types.TCB, logger *slog.Logger) {
	var nuevaCola []types.TCB
	// Itera la cola buscando los TCBs que pertenecen al PCB actual
	for _, tcb := range *cola {
		if tcb.PID != pcb.PID {
			nuevaCola = append(nuevaCola, tcb) // Mantiene los TCBs que no pertenecen al PCB actual
		} else {
			logger.Info(fmt.Sprintf("TCB con TID %d y PID %d eliminado de la cola", tcb.TID, tcb.PID))
		}
	}
	*cola = nuevaCola
}

// ! ACA
func Eliminar_TCBs_de_cola_Block(pcb *types.PCB, cola *[]Bloqueado, logger *slog.Logger) {
	var nuevaCola []Bloqueado
	// Itera la cola buscando los TCBs que pertenecen al PCB actual
	for _, tcb := range *cola {
		if tcb.PID != pcb.PID {
			nuevaCola = append(nuevaCola, tcb) // Mantiene los TCBs que no pertenecen al PCB actual
		} else {
			logger.Info(fmt.Sprintf("TCB con TID %d y PID %d eliminado de la cola", tcb.TID, tcb.PID))
		}
	}
	*cola = nuevaCola
}

// Busca los TCBs del PCB en las colas de Ready y Blocked y los mueve a la cola de Exit
// ! ACA
func Enviar_proceso_a_exit(pid uint32, colaReady *[]types.TCB, colaBlocked *[]Bloqueado, colaExit *[]types.TCB, logger *slog.Logger) bool {

	pcb := Obtener_PCB_por_PID(pid)
	if pcb == nil {
		logger.Error(fmt.Sprintf("No existe el proceso con PID: %d", pid))
		return false
	}

	// Elimina TCBs de la cola de ready y blocked si es que hubiera
	Eliminar_TCBs_de_cola(pcb, colaReady, logger)
	Eliminar_TCBs_de_cola_Block(pcb, colaBlocked, logger) //! ACA
	// Eliminar_TCBs_de_cola(pcb, colaBlocked, logger)

	// Mueve todos los TCBs del PCB a la cola de exit
	for _, tcb := range pcb.TCBs {
		*colaExit = append(*colaExit, tcb)
		logger.Info(fmt.Sprintf("TCB con TID %d movido a la cola de Exit", tcb.TID))
	}

	// Limpiar los TCBs del PCB
	pcb.TCBs = nil
	delete(MapaPCB, pid)
	logger.Info(fmt.Sprintf("Todos los TCBs del PCB con PID %d han sido liberados", pcb.PID))
	return true
}

// Estructuras para manejar los bloqueados
type Motivo int

const ( // Esto funciona mas o menos como el enum de c
	THREAD_JOIN Motivo = iota // Vale 0
	Mutex                     // Vale 1
)

// Como no se puede hacer un slice con un struc generico, hago que el QuienFue sea un string
// Y cuando necesite que sea un uint32 lo parseo
// ACLARACIONES: EL QUIENFUE SE PASA SIEMPRE COMO STRING
type Bloqueado struct {
	PID      uint32 `json:"pid"`
	TID      uint32 `json:"tid"`
	Motivo   Motivo `json:"motivo"`
	QuienFue string `json:"quien_fue"` // si es THREAD_JOIN es un uint32, si es Mutex es un string
}
