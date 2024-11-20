package memsistema

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/sisoputnfrba/tp-golang/memoria/utils"
	"github.com/sisoputnfrba/tp-golang/utils/types"
)

// Mapa para almacenar los contextos de ejecución de los procesos y sus hilos asociados
var ContextosPID = make(map[uint32]types.ContextoEjecucionPID) // Contexto por PID

// Función para inicializar un contexto de ejecución de un proceso (PID)
func CrearContextoPID(pid uint32, base, limite uint32) {
	ContextosPID[pid] = types.ContextoEjecucionPID{
		PID:    pid,
		Base:   base,
		Limite: limite,
		TIDs:   make(map[uint32]types.ContextoEjecucionTID),
	}
	fmt.Printf("Contexto PID %d inicializado con Base = %d, Límite = %d\n", pid, base, limite)
}

// Función para inicializar un contexto de ejecución de un hilo (TID) asociado a un proceso (PID)
func CrearContextoTID(pid, tid uint32, archivoPseudocodigo string) {
	listaInstrucciones := CargarPseudocodigo(int(pid), int(tid), archivoPseudocodigo)
	if proceso, exists := ContextosPID[pid]; exists {
		proceso.TIDs[tid] = types.ContextoEjecucionTID{
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
			LISTAINSTRUCCIONES: listaInstrucciones, // pseudocodigo
		}
		ContextosPID[pid] = proceso // Actualizar el contexto en el mapa
		fmt.Printf("## Contexto Actualizado - (PID:TID) - (%d:%d)\n", pid, tid)
		fmt.Printf("Contexto TID %d inicializado con registros en 0\n", tid)
	} else {
		fmt.Printf("Error: El PID %d no existe\n", pid)
	}
}

// Función para eliminar el contexto de ejecución de un proceso (PID)
func EliminarContextoPID(pid uint32) {
	if _, exists := ContextosPID[pid]; exists {
		delete(ContextosPID, pid)
		fmt.Printf("Contexto PID:, %d eliminado\n", pid)
	} else {
		fmt.Printf("Contexto PID %d no existe\n", pid)
	}
}

// Función para eliminar el contexto de ejecución de un hilo (TID) asociado a un proceso (PID)
func EliminarContextoTID(pid, tid uint32) {
	if proceso, exists := ContextosPID[pid]; exists {
		if _, tidExists := proceso.TIDs[tid]; tidExists {
			delete(proceso.TIDs, tid)
			ContextosPID[pid] = proceso // Actualizar el contexto en el mapa
			fmt.Printf("Contexto TID %d del PID %d eliminado\n", tid, pid)
		} else {
			fmt.Printf("TID %d no existe en el PID %d\n", tid, pid)
		}
	} else {
		fmt.Printf("PID %d no existe\n", pid)
	}
}

func Actualizar_TID(pid uint32, tid uint32, contexto types.ContextoEjecucionTID) {
	if proceso, exists := ContextosPID[pid]; exists {
		if _, tidExists := proceso.TIDs[tid]; tidExists {
			proceso.TIDs[tid] = contexto // Actualizar el contexto en el mapa
			ContextosPID[pid] = proceso  // Actualizar el contexto en el mapa
			fmt.Printf("Contexto TID %d del PID %d actualizado\n", tid, pid)
		} else {
			fmt.Printf("TID %d no existe en el PID %d\n", tid, pid)
		}
	} else {
		fmt.Printf("PID %d no existe\n", pid)
	}
}

// Funcion para cargar el archivo de pseudocodigo

// Funcion para cargar el archivo de pseudocodigo
func CargarPseudocodigo(pid int, tid int, path string) map[string]string {
	file, err := os.Open(utils.Configs.InstructionPath + path)
	if err != nil {
		fmt.Printf("error al abrir el archivo %s: %v", path, err)
	}

	var contextosEjecucion = make(map[int]map[int]*types.ContextoEjecucionTID)
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
		fmt.Printf("error al leer el archivo %s: %v", path, err)
	}
	defer file.Close()
	return contexto.LISTAINSTRUCCIONES
}

func BuscarSiguienteInstruccion(pid, tid uint32, pc uint32) string {

	if proceso, exists := ContextosPID[pid]; exists {
		if hilo, tidExists := proceso.TIDs[tid]; tidExists {
			indiceInstruccion := pc + 1
			instruccion, existe := hilo.LISTAINSTRUCCIONES[fmt.Sprintf("instr_%d", indiceInstruccion)]
			if !existe {
				fmt.Printf("Instrucción no encontrada para PC %d en TID %d", pc, tid)
				return ""
			}

			return instruccion
		} else {
			fmt.Printf("TID %d no existe en el PID %d\n", tid, pid)
			return ""
		}
	} else {
		fmt.Printf("PID %d no existe\n", pid)
		return ""
	}
}
