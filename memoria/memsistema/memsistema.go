package memSistema

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/sisoputnfrba/tp-golang/utils/types"
)

// Mapas para almacenar los contextos de ejecución
var ContextosPID = make(map[uint32]types.ContextoEjecucionPID) // Contexto por PID
var ContextosTID = make(map[uint32]types.ContextoEjecucionTID) // Contexto por TID

// Función para inicializar un contexto de ejecución de un proceso (PID)
func CrearContextoPID(pid uint32, base, limite uint32) {
	ContextosPID[pid] = types.ContextoEjecucionPID{
		PID:    pid,
		Base:   base,
		Limite: limite,
	}
	fmt.Printf("Contexto PID %d inicializado con Base = %d, Límite = %d\n", pid, base, limite)
}

// Función para inicializar un contexto de ejecución de un hilo (TID)
func CrearContextoTID(tid uint32) {
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

// Función para eliminar el contexto de ejecución de un proceso (PID)
func EliminarContextoPID(pid uint32) {
	if _, exists := ContextosPID[pid]; exists {
		delete(ContextosPID, pid)
		fmt.Printf("Contexto PID %d eliminado\n", pid)
	} else {
		fmt.Printf("Contexto PID %d no existe\n", pid)
	}
}

// Función para eliminar el contexto de ejecución de un hilo (TID)
func EliminarContextoTID(tid uint32) {
	if _, exists := ContextosTID[tid]; exists {
		delete(ContextosTID, tid)
		fmt.Printf("Contexto TID %d eliminado\n", tid)
	} else {
		fmt.Printf("Contexto TID %d no existe\n", tid)
	}
}

// Funcion para cargar el archivo de pseudocodigo

// Funcion para cargar el archivo de pseudocodigo
func CargarPseudocodigo(pid int, tid int, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("Error al abrir el archivo %s: %v", path, err)
	}
	defer file.Close()

	var contextosEjecucion = make(map[int]map[int]*types.ContextoEjecucionTID)
	//Si no existe para el PID TID, lo creo
	if _, exists := contextosEjecucion[pid][tid]; !exists {
		contextosEjecucion[pid][tid] = &types.ContextoEjecucionTID{
			TID:                uint32(tid),
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
	}
	contexto := contextosEjecucion[pid][tid]
	scanner := bufio.NewScanner(file)
	instruccionNum := 0 // Indice de instrucciones

	//Empiezo a leer y guardo linea x linea
	for scanner.Scan() {
		linea := scanner.Text()
		contexto.LISTAINSTRUCCIONES[strconv.Itoa(instruccionNum)] = linea
		instruccionNum++
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("Error al leer el archivo %s: %v", path, err)
	}
	return nil
}

func BuscarSiguienteInstruccion(tid uint32, pc uint32) string {

	contexto, existeTID := ContextosTID[tid]

	if !existeTID {
		fmt.Errorf("TID no encontrado")
		return ""

	}

	indiceInstruccion := pc + 1

	instruccion, existe := contexto.LISTAINSTRUCCIONES[fmt.Sprintf("instr_%d", indiceInstruccion)]
	if !existe {
		fmt.Errorf("Instrucción no encontrada para PC %d en TID %d", pc, tid)
		return ""
	}
	//Log obligatorio
	fmt.Printf("Obtener instruccion: ## Obtener instrucción - (PID:TID) - (%d:%d) - Instrucción: %s \n", tid, tid, instruccion)

	return instruccion

}
