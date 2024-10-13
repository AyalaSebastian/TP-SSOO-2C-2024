package main

import (
	
	"github.com/sisoputnfrba/tp-golang/memoria/utils"
	"github.com/sisoputnfrba/tp-golang/utils/logging"
	"github.com/sisoputnfrba/tp-golang/utils/conexiones"
	"github.com/sisoputnfrba/tp-golang/utils/types"
)

func Nueva_Memoria() *types.ContextoEjecucion{
	return &types.ContextoEjecucion{
		RegCPU: make(map[int]*types.RegCPU)
	}
}
func Nuevo_Hilo(tid int, base uint32, limite uint32) {
    registros[tid] = &RegCPU{
        PC:     0,
        AX:     0,
        BX:     0,
        CX:     0,
        DX:     0,
        EX:     0,
        FX:     0,
        GX:     0,
        HX:     0,
        Base:   base,
        Limite: limite,

func Ver_Contexto(pid int, tid int) (*RegCPU, error){
	regCPU, existe := cpuRegisters[tid]
    if !existe {
        return nil, logger.Error(fmt.Sprintf("El hilo con TID %d no existe para el proceso PID %d", tid, pid))
    }
	logger.Info(fmt.Sprintf("## Contexto solicitado (%d : %d)", pid, tid))
    // Devolver el contexto completo
    return regCPU, nil


}

func Update_Contexto(pid uint32, tid uint32){
	regCPU, existe := cpuRegisters[tid]
    if !existe {
        return nil, logger.Error(fmt.Sprintf("El hilo con TID %d no existe para el proceso PID %d", tid, pid))
    }
	regCPU.AX = req.RegCPU.AX
	regCPU.BX = req.RegCPU.BX
	regCPU.CX = req.RegCPU.CX
	regCPU.DX = req.RegCPU.DX
	regCPU.EX = req.RegCPU.EX
	regCPU.FX = req.RegCPU.FX
	regCPU.GX = req.RegCPU.GX
	regCPU.HX = req.RegCPU.HX
	regCPU.PC = req.RegCPU.PC
	regCPU.Base = req.RegCPU.Base
	regCPU.Limite = req.RegCPU.Limite

	
	logger.Info(fmt.Sprintf("## Contexto actualizado (%d : %d)", pid, tid))
    return
}




func main() {
	// Inicializamos la configuracion y el logger
	utils.Configs = utils.Iniciar_configuracion("config.json")
	logger := logging.Iniciar_Logger("memoria.log", utils.Configs.LogLevel)

	// Inicializamos la memoria (Lo levantamos como servidor)
	utils.Iniciar_memoria(logger)
	memoria := Nueva_Memoria()
}
