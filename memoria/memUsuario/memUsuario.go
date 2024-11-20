package memUsuario

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/sisoputnfrba/tp-golang/memoria/memsistema"
	"github.com/sisoputnfrba/tp-golang/memoria/utils"
	"github.com/sisoputnfrba/tp-golang/utils/types"
)

// var memoria global
var MemoriaDeUsuario []byte
var Particiones []types.Particion
var ParticionesDinamicas []int
var BitmapParticiones []bool
var PidAParticion map[uint32]int // Mapa para rastrear la asignación de PIDs a particiones

// Funcion para iniciar la memoria y definir las particiones
func Inicializar_Memoria_De_Usuario(logger *slog.Logger) {
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
		logger.Info("Partición fija", "Número", i+1, "iniciada con Base", particion.Base, "y Límite", particion.Limite)
		base += uint32(limite)
	}
	// Inicializar el bitmap y el mapa de PIDs
	//todas las particiones estan libres = false
	BitmapParticiones = make([]bool, len(Particiones))
	PidAParticion = make(map[uint32]int)
}

// memoria para particiones dinamicas
func Inicializar_Memoria_Dinamica(logger *slog.Logger) {
	MemoriaDeUsuario = make([]byte, utils.Configs.MemorySize)
	particion := types.Particion{
		Base:   0,
		Limite: uint32(utils.Configs.MemorySize),
	}
	Particiones = []types.Particion{particion}
	BitmapParticiones = []bool{false} // Un solo valor para la memoria dinámica
	ParticionesDinamicas = append(ParticionesDinamicas, 1024)
	logger.Info(fmt.Sprintf("Memoria dinámica inicializada: ## Base = %d, Límite = %d", particion.Base, particion.Limite))
	// Inicializar el mapa de PIDs
	PidAParticion = make(map[uint32]int)
}

// Función para liberar una partición por PID
func LiberarParticionPorPID(pid uint32, logger *slog.Logger) error {
	particion, existe := PidAParticion[pid]
	if !existe {
		return fmt.Errorf("no se encontró el proceso %d asignado a ninguna partición", pid)
	}

	// Liberar la partición y actualizar el bitmap
	BitmapParticiones[particion] = false
	delete(PidAParticion, pid) // Eliminar la entrada del mapa
	fmt.Printf("Proceso %d liberado de la partición %d\n", pid, particion+1)

	// Comprobar si el esquema es dinámico
	if utils.Configs.Scheme == "DINAMICA" {
		var particionIndex int
		// Verificar si hay particiones libres adyacentes
		adyacenteIzquierdaLibre := particionIndex > 0 && !BitmapParticiones[particionIndex-1]
		adyacenteDerechaLibre := particionIndex < len(Particiones)-1 && !BitmapParticiones[particionIndex+1]

		// Llamar a combinarParticionesLibres si alguna adyacente está libre
		if adyacenteIzquierdaLibre || adyacenteDerechaLibre {
			combinarParticionesLibres(particionIndex, logger)
		}
	}
	return nil
}

// combinar particiones dinamicas libres en una sola
func combinarParticionesLibres(index int, logger *slog.Logger) {
	// Inicializar variables para los límites de la nueva partición combinada
	base := Particiones[index].Base
	limite := Particiones[index].Limite

	// Verificar partición anterior, si existe y está libre
	if index > 0 && !BitmapParticiones[index-1] {
		base = Particiones[index-1].Base
		limite += Particiones[index-1].Limite
		// Eliminar partición anterior ya que se combina
		Particiones = append(Particiones[:index-1], Particiones[index:]...)
		BitmapParticiones = append(BitmapParticiones[:index-1], BitmapParticiones[index:]...)
		index-- // Ajustar el índice ya que hemos eliminado la partición anterior
	}

	// Verificar partición siguiente, si existe y está libre
	if index < len(Particiones)-1 && !BitmapParticiones[index+1] {
		limite += Particiones[index+1].Limite
		// Eliminar partición siguiente ya que se combina
		Particiones = append(Particiones[:index+1], Particiones[index+2:]...)
		BitmapParticiones = append(BitmapParticiones[:index+1], BitmapParticiones[index+2:]...)
	}

	// Actualizar la partición combinada en el índice actual
	Particiones[index] = types.Particion{Base: base, Limite: limite}
	BitmapParticiones[index] = false // Marcar como libre

	logger.Info("Particiones combinadas", "Nueva Base", base, "Nuevo Límite", limite)
}

func AsignarPID(pid uint32, tamanio_proceso int, path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var asigno = false
		algoritmo := utils.Configs.SearchAlgorithm
		esquema := utils.Configs.Scheme
		if esquema == "FIJAS" {
			switch algoritmo {
			case "FIRST":
				asigno = FirstFitFijo(pid, tamanio_proceso, path)

			case "BEST":
				asigno = BestFitFijo(pid, tamanio_proceso, path)
			case "WORST":
				asigno = WorstFitFijo(pid, tamanio_proceso, path)
			}
			if !asigno {
				(http.Error(w, "NO SE PUDO INICIALIZAR EL PROCESO POR FALTA DE HUECOS EN LAS PARTICIONES", http.StatusInternalServerError))
				return
			} else {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
				return
			}
		} else if esquema == "DINAMICAS" {
			switch algoritmo {
			case "FIRST":
				asigno = FirstFitDinamico(pid, tamanio_proceso, path)
			case "BEST":
				asigno = BestFitDinamico(pid, tamanio_proceso, path)
			case "WORST":
				asigno = WorstFitDinamico(pid, tamanio_proceso, path)
			}
			if asigno {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
				return
			} else {
				compactar := SePuedeCompactar(tamanio_proceso)
				if compactar {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("COMPACTAR"))
				} else {
					(http.Error(w, "NO SE PUDO INICIALIZAR EL PROCESO POR FALTA DE HUECOS EN LAS PARTICIONES", http.StatusInternalServerError))
				}
			}
		}
	}
}

// first fit para particiones fijas
func FirstFitFijo(pid uint32, tamanio_proceso int, path string) bool {
	particion := utils.Configs.Partitions
	for i := 0; i < len(BitmapParticiones); i++ {
		if !BitmapParticiones[i] {
			if tamanio_proceso < particion[i] {
				PidAParticion[pid] = i
				BitmapParticiones[i] = true
				fmt.Printf("Proceso %d asignado a la partición %d\n", pid, i+1)
				memsistema.CrearContextoPID(pid, uint32(Particiones[i].Base), uint32(Particiones[i].Limite))
				return true
			}
		}
	}
	return false
}

// best fit para particiones fijas
func BestFitFijo(pid uint32, tamanio_proceso int, path string) bool {

	particiones := utils.Configs.Partitions
	var menor = 1024
	var pos_menor = -1
	for i := 0; i < len(BitmapParticiones); i++ {
		if !BitmapParticiones[i] {
			if tamanio_proceso < particiones[i] {
				if particiones[i] < menor {
					menor = particiones[i]
					pos_menor = i
				}
			}
		}
	}
	if pos_menor == -1 {
		return false
	} else {
		PidAParticion[pid] = pos_menor
		BitmapParticiones[pos_menor] = true
		fmt.Printf("Proceso %d asignado a la partición %d\n", pid, pos_menor+1)
		memsistema.CrearContextoPID(pid, uint32(Particiones[pos_menor].Base), uint32(Particiones[pos_menor].Limite))
		return true
	}
}

// worst fit para particiones fijas
func WorstFitFijo(pid uint32, tamanio_proceso int, path string) bool {
	particiones := utils.Configs.Partitions
	var mayor = 0
	var pos_mayor = -1
	for i := 0; i < len(BitmapParticiones); i++ {
		if !BitmapParticiones[i] {
			if tamanio_proceso < particiones[i] {
				if particiones[i] > mayor {
					mayor = particiones[i]
					pos_mayor = i
				}
			}
		}
	}
	if pos_mayor == -1 {
		return false
	} else {
		PidAParticion[pid] = pos_mayor
		BitmapParticiones[pos_mayor] = true
		fmt.Printf("Proceso %d asignado a la partición %d\n", pid, pos_mayor+1)
		memsistema.CrearContextoPID(pid, uint32(Particiones[pos_mayor].Base), uint32(Particiones[pos_mayor].Limite))
		return true
	}
}

// empiezo con un solo espacio de memoria de 1024 bytes, si no esta reservado lo hago con el pid entrante, sino no hay espacio
func FirstFitDinamico(pid uint32, tamanio_proceso int, path string) bool {
	for i := 0; i < len(BitmapParticiones); i++ {
		if !BitmapParticiones[i] {
			if tamanio_proceso < ParticionesDinamicas[i] {
				AsignarParticion(pid, i, tamanio_proceso)
				return true
			}
		}
	}
	return false
}
func BestFitDinamico(pid uint32, tamanio_proceso int, path string) bool {
	var pos_menor = -1
	var menor = 1024
	for i := 0; i < len(BitmapParticiones); i++ {
		if !BitmapParticiones[i] {
			if tamanio_proceso < ParticionesDinamicas[i] {
				if ParticionesDinamicas[i] <= menor {
					menor = ParticionesDinamicas[i]
					pos_menor = i
				}
			}
		}
	}
	if pos_menor == -1 {
		return false // no hay huecos hay que compactar o tirar interrupcion
	} else {
		AsignarParticion(pid, pos_menor, tamanio_proceso)
		return true
	}
}

// empiezo con un solo espacio de memoria de 1024 bytes, si no esta reservado lo hago con el pid entrante, sino no hay espacio
func WorstFitDinamico(pid uint32, tamanio_proceso int, path string) bool {
	var pos_mayor = -1
	var mayor = 0
	for i := 0; i < len(BitmapParticiones); i++ {
		if !BitmapParticiones[i] {
			if tamanio_proceso < ParticionesDinamicas[i] {
				if ParticionesDinamicas[i] >= mayor {
					mayor = ParticionesDinamicas[i]
					pos_mayor = i
				}
			}
		}
	}
	if pos_mayor == -1 {
		return false
	} else {
		AsignarParticion(pid, pos_mayor, tamanio_proceso)
		return true
	}
}

// Hay que calcular bien las bases y limites que estan mal
func baseDinamica(posicion int) uint32 {

	if posicion <= 0 {
		return uint32(0)
	} else {
		return uint32(Particiones[posicion-1].Base + Particiones[posicion-1].Limite)
	}
}

func AsignarParticion(pid uint32, posicion, tamanio_proceso int) {
	nuevaParticion := ParticionesDinamicas[posicion] - tamanio_proceso
	ParticionesDinamicas[posicion] = tamanio_proceso
	BitmapParticiones[posicion] = true
	PidAParticion[pid] = posicion
	BitmapParticiones = append(BitmapParticiones, false)
	ParticionesDinamicas = append(ParticionesDinamicas, nuevaParticion)
	fmt.Printf("Proceso %d asignado a la partición %d\n", pid, posicion+1)
	base := baseDinamica(posicion)
	memsistema.CrearContextoPID(pid, base, uint32(tamanio_proceso))
}

func SePuedeCompactar(tamanio_proceso int) bool {
	var espacioLibre = 0
	particiones := utils.Configs.Partitions
	for i := 0; i < len(BitmapParticiones); i++ {
		if !BitmapParticiones[i] {
			espacioLibre += particiones[i]
		}
	}
	if espacioLibre >= tamanio_proceso {
		return true
	}
	return false
}
