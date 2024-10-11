package utils

import (
	"fmt"
	"log/slog"

	"github.com/sisoputnfrba/tp-golang/utils/types"
)

// Agrega un elemento a la cola (IMPORTANTE: la cola debe coincidar con el tipo de elemento)
func Encolar[T any](cola *[]T, elemento T) {
	*cola = append(*cola, elemento)
}

// Desencola un elemento de la cola y retorna ese elemento
func Desencolar[T any](cola *[]T) T {
	if len(*cola) == 0 {
		var vacio T
		return vacio // O manejar el caso de cola vacía
	}
	elemento := (*cola)[0]
	*cola = (*cola)[1:] // Elimina el primer elemento
	return elemento
}

func Sacar_TCB_Del_Map(mapaPCBS *map[uint32]types.PCB, pid uint32, tid uint32, logger *slog.Logger) {
	pcb, existe := (*mapaPCBS)[pid]
	if !existe {
		logger.Error(fmt.Sprintf("El PCB con PID %d no existe", pid))
		return
	}

	// Verificamos si el TCB existe dentro del PCB
	_, existeTCB := pcb.TCBs[tid]
	if !existeTCB {
		logger.Error(fmt.Sprintf("El TCB con TID %d no existe en el PCB con PID %d", tid, pid))
		return
	}

	// Eliminamos el TCB del mapa de TCBs
	delete(pcb.TCBs, tid)
	logger.Info(fmt.Sprintf("TCB con TID %d eliminado del PCB con PID %d", tid, pid)) // Esto no hace falta que vaya

	// Actualizamos el PCB en el mapa de PCBs
	(*mapaPCBS)[pid] = pcb
}

// Devuelve dos valores, el primero es la solicitud y el segundo es un booleano que indica si la cola está vacía
func Proxima_solicitud(cola *[]SolicitudIO) (SolicitudIO, bool) {
	if len(*cola) == 0 {
		return SolicitudIO{}, false
	}
	solicitud := (*cola)[0]
	*cola = (*cola)[1:] // Remueve la primera solicitud del slice
	return solicitud, true
}
