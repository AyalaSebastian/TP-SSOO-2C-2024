package utils

import (
	"fmt"
	"log/slog"

	"github.com/sisoputnfrba/tp-golang/utils/types"
)

// Uso [T any] para poder encolar y desencolar tanto PCBs como TCBs

func Encolar[T any](cola *[]T, elemento T) {
	*cola = append(*cola, elemento)
}

func Desencolar[T any](cola *[]T) T {
	if len(*cola) == 0 {
		var vacio T
		return vacio // O manejar el caso de cola vacía
	}
	elemento := (*cola)[0]
	*cola = (*cola)[1:] // Elimina el primer elemento
	return elemento
}

func Sacar_TCB_Del_Slice(mapaPCBS *map[uint32]types.PCB, pid uint32, tid uint32, logger *slog.Logger) {
	var nuevaCola []types.TCB

	// Dejamos los TCB con TID distintos a los del parametro
	for _, tcb := range (*mapaPCBS)[pid].TCBs {
		if tcb.TID != tid {
			nuevaCola = append(nuevaCola, tcb)
		} else {
			logger.Info(fmt.Sprintf("TCB del PID: %d, TID: %d eliminado de la cola de TCBs", pid, tid))
		}
	}
	pcb := (*mapaPCBS)[pid]
	pcb.TCBs = nuevaCola
	(*mapaPCBS)[pid] = pcb
}
