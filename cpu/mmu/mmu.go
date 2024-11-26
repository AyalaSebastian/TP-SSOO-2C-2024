package mmu

import (
	"errors"
	"fmt"
	"log"
	"log/slog"

	"github.com/sisoputnfrba/tp-golang/cpu/client"
	"github.com/sisoputnfrba/tp-golang/utils/types"
)

func TraducirDireccion(proceso *types.Proceso, direccionLogica uint32, logger *slog.Logger) (uint32, error) {

	direccionFisica := proceso.ContextoEjecucion.Registros.Base + direccionLogica
	if direccionFisica >= proceso.ContextoEjecucion.Registros.Limite {
		client.EnviarContextoDeEjecucion(proceso, "actualizar_contexto", logger)
		logger.Info(fmt.Sprintf("## TID: %d - Actualizo Contexto Ejecución", proceso.Tid))
		// Devolver el Tid al Kernel con motivo de Segmentation Fault
		if client.DevolverTIDAlKernel(proceso.Tid, logger, "THREAD_INTERRUPT", "Segmentation Fault") {
			log.Printf("Segmentation Fault en Tid %d", proceso.Tid)
		}
		return 0, errors.New("segmentation fault")
	}

	return direccionFisica, nil
}
