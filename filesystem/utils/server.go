package utils

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/sisoputnfrba/tp-golang/utils/conexiones"
	"github.com/sisoputnfrba/tp-golang/utils/types"
)

// Usar tamanio de puntero uint32
// Creo que esta funcion MetaDataFile no va, la memoria nos envia el nombre del archivo en la peticion del DUMP, el cual ya viene con todos estos datos
type MetadataFile struct {
	PID       uint32
	TID       uint32
	Timestamp time.Time
}

func Iniciar_fileSystem(logger *slog.Logger) {
	mux := http.NewServeMux()

	// Endpoints
	mux.HandleFunc("/dump", DUMP(logger))

	conexiones.LevantarServidor(strconv.Itoa(Configs.Port), mux, logger)
}

func DUMP(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var magic types.DumpFile
		err := decoder.Decode(&magic)
		if err != nil {
			logger.Error(fmt.Sprintf("Error al decodificar mensaje: %s\n", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error al decodificar mensaje"))
			return
		}
		// Verificar si se cuenta con el espacio disponible
		bloquesNecesarios, espacioSuficiente := Verificar_Espacio_Disponible(magic.Tamanio, logger)
		if !espacioSuficiente {
			logger.Error("No hay espacio suficiente para el archivo")
			w.WriteHeader(http.StatusInsufficientStorage)
			w.Write([]byte("No hay espacio suficiente para el archivo"))
			return
		}

		// Reservar el bloque de índice y los bloques de datos correspondientes en el bitmap
		Reservar_Bloques_Del_Bitmap(bloquesNecesarios, magic.Tamanio, logger)

		// Crear archivo de metadata
		Crear_Archivo_Metadata(magic.Nombre, int(bloquesNecesarios[0]), magic.Tamanio, logger) //Tengo dudas con el segundo parametro

		// Acceder al archivo de punteros y grabar todos los punteros reservados
		// codigo

		// Acceder bloque a bloque e ir escribiendo el contenido de la memoria.
		// codigo

	}
}

// Retorna slice de los bloques libres
// func Bloques_Libres(logger *slog.Logger) []byte {

// 	// Cargar archivo en un slice de bytes
// 	file, err := os.ReadFile(Configs.MountDir + "/bitmap.dat")
// 	if err != nil {
// 		logger.Error(fmt.Sprintf("Error al leer el archivo bitmap.dat: %s\n", err.Error()))
// 		return nil
// 	}

// 	// Identificar los bloques libres
// 	var bitesLibres []byte
// 	for i := 0; i < len(file); i++ {
// 		for j := 0; j < 8; j++ {
// 			if (file[i] & (1 << j)) == 0 {
// 				bitesLibres = append(bitesLibres, byte(i*8+j))
// 			}
// 		}
// 	}
// 	return bitesLibres
// }

// Devuelve dos valores: un array con los indices a reservar([]byte) y si hay espacio suficiente para el archivo(true/false)
func Verificar_Espacio_Disponible(tamanioArchivo int, logger *slog.Logger) ([]byte, bool) {
	bloquesNecesarios := tamanioArchivo / Configs.BlockSize
	if tamanioArchivo%Configs.BlockSize > 0 {
		bloquesNecesarios++ // Si el tamaño no es multiplo del BlockSize, se necesita un bloque más
	}

	// Identificar los bloques libres
	var bitesLibres []byte
	totalBitesLibres := 0

	for i := 0; i < len(Bitmap); i++ {
		for j := 0; j < 8; j++ {
			if (Bitmap[i] & (1 << j)) == 0 {
				totalBitesLibres++

				if len(bitesLibres) < bloquesNecesarios+1 { // +1 para el bloque de índice
					bitesLibres = append(bitesLibres, byte(i*8+j))
				}
			}
		}
	}
	if bloquesNecesarios > totalBitesLibres {
		return bitesLibres, false
	}

	return bitesLibres, true
}

// Reservo los bloques en el slice
func Reservar_Bloques_Del_Bitmap(bloques []byte, cantidad int, logger *slog.Logger) {

	// Reservar los bloques
	for i := 0; i < cantidad; i++ {
		bloque := bloques[i]
		byteIndex := bloque / 8
		bitIndex := bloque % 8
		Bitmap[byteIndex] |= (1 << bitIndex) // Marcar el bit como ocupado
	}

	// Guardar los cambios en el archivo
	err := os.WriteFile(Configs.MountDir+"/bitmap.dat", Bitmap, 0644)
	if err != nil {
		logger.Error(fmt.Sprintf("Error al escribir el archivo bitmap.dat: %s", err.Error()))
		return
	}

	logger.Info(fmt.Sprintf("Reservados %d bloques en el archivo bitmap.dat", cantidad))
}

func Crear_Archivo_Metadata(nombreArchivo string, indice int, tamanio int, logger *slog.Logger) {
	rutaArchivo := Configs.MountDir + "/files/" + nombreArchivo

	metadata := fmt.Sprintf(`{
        "index_block": %d,
        "size": %d
    }`, indice, tamanio)

	err := os.WriteFile(rutaArchivo, []byte(metadata), 0644)
	if err != nil {
		logger.Error(fmt.Sprintf("error al crear archivo de metadata: %v", err))
		return
	}
}
