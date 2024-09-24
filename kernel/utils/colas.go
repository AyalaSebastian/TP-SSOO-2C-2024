package utils

import (
	"github.com/sisoputnfrba/tp-golang/utils/types"
)

func Encolar(cola *[]types.PCB, pcb types.PCB) {
	*cola = append(*cola, pcb)
}

func Desencolar(cola *[]types.PCB) types.PCB {
	if len(*cola) == 0 {
		return types.PCB{} // O manejar el caso de cola vacía
	}
	pcb := (*cola)[0]
	*cola = (*cola)[1:] // Elimina el primer elemento
	return pcb
}
