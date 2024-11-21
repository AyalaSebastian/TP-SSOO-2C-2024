package mmu

import (
	"errors"
	"log"
	"log/slog"

	"github.com/sisoputnfrba/tp-golang/cpu/client"
	"github.com/sisoputnfrba/tp-golang/utils/types"
)

func TraducirDireccion(proceso *types.Proceso, direccionLogica uint32, logger *slog.Logger) (uint32, error) {
	/*
		particion, existe := particiones[tid]
		if !existe {
			return 0, errors.New("partición no encontrada")
		}
	*/direccionFisica := proceso.ContextoEjecucion.Registros.Base + direccionLogica
	if direccionFisica >= proceso.ContextoEjecucion.Registros.Limite {
		client.EnviarContextoDeEjecucion(proceso, "actualizar_contexto", logger)
		// Devolver el Tid al Kernel con motivo de Segmentation Fault
		if client.DevolverTIDAlKernel(proceso.Tid, logger, "THREAD_INTERRUPT", "Segmentation Fault") {
			log.Printf("Segmentation Fault en Tid %d", proceso.Tid)
		}
		return 0, errors.New("segmentation fault")
	}

	return direccionFisica, nil
}
