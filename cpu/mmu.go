package main

/*
import (
	"errors"
	"log"
)

var (
	particiones map[int]Particion
	contextos   map[int]ContextoEjecucion
)

func traducirDireccion(tid, direccionLogica int) (int, error) {
	//particion, existe := particiones[tid]
	if !existe {
		return 0, errors.New("partición no encontrada")
	}

	if direccionLogica >= particion.Limite {
		actualizarContextoSegmentationFault(tid)
		return 0, errors.New("segmentation fault")
	}

	direccionFisica := particion.Base + direccionLogica
	return direccionFisica, nil
}

func actualizarContextoSegmentationFault(tid int) {
	contexto, existe := contextos[tid]
	if !existe {
		log.Printf("Contexto de ejecución no encontrado para Tid %d", tid)
		return
	}

	// Actualizar el contexto de ejecución en memoria
	// Aquí se puede agregar la lógica necesaria para actualizar el contexto

	// Devolver el Tid al Kernel con motivo de Segmentation Fault
	log.Printf("Segmentation Fault en Tid %d", contexto.Tid)
}
*/
