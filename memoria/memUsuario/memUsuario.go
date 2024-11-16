package memUsuario

import (
	"fmt"
	"net/http"

	"github.com/sisoputnfrba/tp-golang/memoria/memsistema"
	"github.com/sisoputnfrba/tp-golang/memoria/utils"
	"github.com/sisoputnfrba/tp-golang/utils/types"
)

// var memoria global
var MemoriaDeUsuario []byte
var Particiones []types.Particion
var BitmapParticiones []bool
var PidAParticion map[uint32]int // Mapa para rastrear la asignación de PIDs a particiones

// Funcion para iniciar la memoria y definir las particiones
func Inicializar_Memoria_De_Usuario() {
	// Inicializar el espacio de memoria con 1024 bytes
	MemoriaDeUsuario = make([]byte, utils.Configs.MemorySize)

	// Asignar las particiones fijas en la memoria usando los datos de config
	var base uint32 = 0
	for i, limite := range utils.Configs.Partitions {
		particion := types.Particion{
			Base:   base,
			Limite: uint32(limite),
		}
		Particiones = append(Particiones, particion)
		fmt.Printf("Partición %d inicializada: Base = %d, Límite = %d\n", i+1, particion.Base, particion.Limite)
		base += uint32(limite)
	}
	// Inicializar el bitmap y el mapa de PIDs
	//todas las particiones estan libres = false
	BitmapParticiones = make([]bool, len(Particiones))
	PidAParticion = make(map[uint32]int)
}

// Función para liberar una partición por PID
func LiberarParticionPorPID(pid uint32) error {
	particion, existe := PidAParticion[pid]
	if !existe {
		return fmt.Errorf("no se encontró el proceso %d asignado a ninguna partición", pid)
	}

	// Liberar la partición y actualizar el bitmap
	BitmapParticiones[particion] = false
	delete(PidAParticion, pid) // Eliminar la entrada del mapa
	fmt.Printf("Proceso %d liberado de la partición %d\n", pid, particion+1)
	return nil
}

// first fit
func AsignarPID(pid uint32, tamanio_proceso int, path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		particiones := utils.Configs.Partitions
		for i := 0; i < len(BitmapParticiones); i++ {
			if !BitmapParticiones[i] {
				if tamanio_proceso < particiones[i] {
					PidAParticion[pid] = i
					BitmapParticiones[i] = true
					fmt.Printf("Proceso %d asignado a la partición %d\n", pid, i+1)
					memsistema.CrearContextoPID(pid, uint32(Particiones[i].Base), uint32(Particiones[i].Limite))
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("OK"))
					return
				}
			}
		}
		(http.Error(w, "NO SE PUDO INICIALIZAR EL PROCESO POR FALTA DE HUECOS EN LAS PARTICIONES", http.StatusInternalServerError))
	}
}

/*
modificaciones:
// Función para asignar un proceso a la primera partición libre
func AsignarProcesoAParticion(pid uint32) error {
	for i, ocupada := range BitmapParticiones {
		if !ocupada { // Si la partición está libre
			BitmapParticiones[i] = true   // Marcar como ocupada
			PidAParticion[pid] = i        // Asignar el PID a esta partición
			fmt.Printf("Proceso %d asignado a la partición %d\n", pid, i+1)
			return nil
		}
	}
	return fmt.Errorf("No hay particiones libres para asignar el proceso %d", pid)
}
*/
