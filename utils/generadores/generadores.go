package generadores

import (
	"github.com/sisoputnfrba/tp-golang/utils/types"
)

var PidCounter uint32 = 0

// MapaParaTcbs := make(map[string]int)

// Genera un PID único (el tipo de dato uint32 es para que no tome valore negativos).
func Generar_PID() uint32 {
	PidCounter++
	return PidCounter
}

// Genera un PCB con un PID único y con las listas de TCBs y Mutexs vacías.
func Generar_PCB() types.PCB {
	mutex := make(map[string]string)
	// tcbs := make(map[uint32]types.TCB)
	// countTCB := 0
	return types.PCB{
		PID:    Generar_PID(),
		TCBs:   []types.TCB{},
		Mutexs: mutex,
	}
}

// Genera un nuevo TCB y lo añade al PCB recibido por parámetro (pasar el pcb con &).
func Generar_TCB(pcb *types.PCB, prioridad int) types.TCB {

	tid := len(pcb.TCBs) // Usamos la longitud actual de TCBs para generar el próximo TID.
	tidUint32 := uint32(tid)

	tcb := types.TCB{
		TID:       tidUint32,
		Prioridad: prioridad,
		PID:       pcb.PID,
	}

	pcb.TCBs = append(pcb.TCBs, tcb) // Aniadimos el nuevo TCB al PCB.
	return tcb
}
