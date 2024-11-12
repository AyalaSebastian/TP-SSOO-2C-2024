package memUsuario

import (
	"fmt"
	"net/http"

	"github.com/sisoputnfrba/tp-golang/memoria/memSistema"
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
		return fmt.Errorf("No se encontró el proceso %d asignado a ninguna partición", pid)
	}

	// Liberar la partición y actualizar el bitmap
	BitmapParticiones[particion] = false
	delete(PidAParticion, pid) // Eliminar la entrada del mapa
	fmt.Printf("Proceso %d liberado de la partición %d\n", pid, particion+1)
	return nil
}

func AsignarPID(pid uint32, tamanio_proceso int, path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		particiones := utils.Configs.Partitions
		cant_particiones := len(particiones)
		for i := 0; i < cant_particiones; i++ {
			if particiones[i] != (-1) {
				if tamanio_proceso < len(MemoriaDeUsuario) {
					particiones[i] = -1
					memSistema.CrearContextoPID(pid, uint32(tamanio_proceso), uint32(len(MemoriaDeUsuario)))
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("OK"))
					return
				}
			}
		}
		(http.Error(w, "NO SE PUDO INICIALIZAR EL PROCESO POR FALTA DE HUECOS EN LAS PARTICIONES", http.StatusInternalServerError))
		return
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
