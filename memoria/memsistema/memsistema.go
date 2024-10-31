package memsistema

import (
	"fmt"

	"github.com/sisoputnfrba/tp-golang/memoria/utils"
	"github.com/sisoputnfrba/tp-golang/utils/types"
)

// var memoria global
var Memoria []byte

// Función para iniciar la memoria y definir las particiones
func Inicializar_Memoria() {

	// Inicializar el espacio de memoria con 1024 bytes
	Memoria = make([]byte, utils.Configs.MemorySize)

	// Asignar las particiones fijas en la memoria
	for i, particion := range Particiones {
		fmt.Printf("Partición %d inicializada: Base = %d, Límite = %d\n", i+1, particion.Base, particion.Limite)
	}
}

// Definición de las particiones fijas
var Particiones = []types.Particion{
	{Base: 0, Limite: 512},   // Primera partición: del byte 0 al byte 511
	{Base: 512, Limite: 16},  // Segunda partición: del byte 512 al 527
	{Base: 528, Limite: 32},  // Tercera partición: del byte 528 al 559
	{Base: 560, Limite: 16},  // Cuarta partición: del byte 560 al 575
	{Base: 576, Limite: 256}, // Quinta partición: del byte 576 al 831
	{Base: 832, Limite: 64},  // Sexta partición: del byte 832 al 895
	{Base: 896, Limite: 128}, // Séptima partición: del byte 896 al 1023
}

// Mapas para almacenar los contextos de ejecución
var contextosPID = make(map[int]types.ContextoEjecucionPID) // Contexto por PID
var contextosTID = make(map[int]types.ContextoEjecucionTID) // Contexto por TID

// Función para inicializar un contexto de ejecución de un proceso (PID)
func inicializarContextoPID(pid int, base, limite uint32) {
	contextosPID[pid] = types.ContextoEjecucionPID{
		PID:    pid,
		Base:   base,
		Limite: limite,
	}
	fmt.Printf("Contexto PID %d inicializado con Base = %d, Límite = %d\n", pid, base, limite)
}

// Función para inicializar un contexto de ejecución de un hilo (TID)
func inicializarContextoTID(tid int) {
	contextosTID[tid] = types.ContextoEjecucionTID{
		TID:                tid,
		PC:                 0,
		AX:                 0,
		BX:                 0,
		CX:                 0,
		DX:                 0,
		EX:                 0,
		FX:                 0,
		GX:                 0,
		HX:                 0,
		LISTAINSTRUCCIONES: make(map[string]string),
	}
	fmt.Printf("Contexto TID %d inicializado con registros en 0\n", tid)
}
