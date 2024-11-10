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

// Funcion para iniciar la memoria y definir las particiones
func Inicializar_Memoria_De_Usuario() {
	// Inicializar el espacio de memoria con 1024 bytes
	MemoriaDeUsuario = make([]byte, utils.Configs.MemorySize)
	// Asignar las particiones fijas en la memoria
	for i, particion := range Particiones {
		fmt.Printf("Partición %d inicializada: Base = %d, Límite = %d\n", i+1, particion.Base, particion.Limite)
	}
}

// Definicion de las particiones fijas
var Particiones = []types.Particion{
	{Base: 0, Limite: 512},   // Primera particion: del byte 0 al byte 511
	{Base: 512, Limite: 16},  // Segunda particion: del byte 512 al 527
	{Base: 528, Limite: 32},  // Tercera particion: del byte 528 al 559
	{Base: 560, Limite: 16},  // Cuarta particion: del byte 560 al 575
	{Base: 576, Limite: 256}, // Quinta particion: del byte 576 al 831
	{Base: 832, Limite: 64},  // Sexta particion: del byte 832 al 895
	{Base: 896, Limite: 128}, // Septima particion: del byte 896 al 1023
}

// Bitmap para las particiones (true = ocupada, false = libre)
var BitmapParticiones = make([]bool, len(Particiones))

//true es ocupada y false es libre

// Función para marcar una partición como ocupada o libre
func MarcarParticion(indiceDeParticion int, ocupada bool) error {
	if indiceDeParticion < 0 || indiceDeParticion >= len(BitmapParticiones) {
		return fmt.Errorf("índice de partición fuera de rango")
	}
	BitmapParticiones[indiceDeParticion] = ocupada
	return nil
}

// indiceDeParticion= 0 para particion 1 y 7 para la particion 8
// Función para verificar si una partición está ocupada
func EstaParticionOcupada(indiceDeParticion int) (bool, error) {
	if indiceDeParticion < 0 || indiceDeParticion >= len(BitmapParticiones) {
		return false, fmt.Errorf("índice de partición fuera de rango")
	}
	return BitmapParticiones[indiceDeParticion], nil
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
		(http.Error(w, "NO SE PUDO INICIALIZAR EL PROCESO POR FALTA DE ESPACIO LIBRE EN LA MEMORIA", http.StatusInternalServerError))
		return
	}
}
