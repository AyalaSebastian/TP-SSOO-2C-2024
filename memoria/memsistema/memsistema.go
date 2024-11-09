package memSistema

import (
	"fmt"

	"github.com/sisoputnfrba/tp-golang/utils/types"
)

// Mapas para almacenar los contextos de ejecución
var ContextosPID = make(map[int]types.ContextoEjecucionPID) // Contexto por PID
var ContextosTID = make(map[int]types.ContextoEjecucionTID) // Contexto por TID

// Función para inicializar un contexto de ejecución de un proceso (PID)
func crearContextoPID(pid int, base, limite uint32) {
	ContextosPID[pid] = types.ContextoEjecucionPID{
		PID:    pid,
		Base:   base,
		Limite: limite,
	}
	fmt.Printf("Contexto PID %d inicializado con Base = %d, Límite = %d\n", pid, base, limite)
}

// Función para inicializar un contexto de ejecución de un hilo (TID)
func crearContextoTID(tid int) {
	ContextosTID[tid] = types.ContextoEjecucionTID{
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
		LISTAINSTRUCCIONES: make(map[string]string), // pseudocodigo
	}
	fmt.Printf("Contexto TID %d inicializado con registros en 0\n", tid)
}

func BuscarSiguienteInstruccion(tid uint32, pc uint32) {

}
