package mmu

import (
	"errors"
	"log"
	"log/slog"

	"github.com/sisoputnfrba/tp-golang/cpu/client"
)

func TraducirDireccion(tid uint32, direccionLogica uint32, base uint32, limite uint32, Logger slog.Logger) (uint32, error) {
	/*
		particion, existe := particiones[tid]
		if !existe {
			return 0, errors.New("partición no encontrada")
		}
	*/
	if direccionLogica <= limite {
		actualizarContextoSegmentationFault(tid, Logger)
		return 0, errors.New("segmentation fault")
	}

	direccionFisica := base + direccionLogica
	return direccionFisica, nil
}

func actualizarContextoSegmentationFault(tid uint32, Logger *slog.Logger) {
	/*	contexto, existe := contextos[tid]
		if !existe {
			log.Printf("Contexto de ejecución no encontrado para Tid %d", tid)
			return
		}*/

	// Actualizar el contexto de ejecución en memoria
	// Aquí se puede agregar la lógica necesaria para actualizar el contexto

	// Devolver el Tid al Kernel con motivo de Segmentation Fault
	if client.DevolverTIDAlKernel(tid, Logger, "THREAD_INTERRUPT", "Segmentation Fault") {
		log.Printf("Segmentation Fault en Tid %d", tid)
	}

}
