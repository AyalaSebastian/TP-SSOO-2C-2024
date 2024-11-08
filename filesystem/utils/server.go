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
		// sigue el codigo

		// Verificar si se cuenta con el espacio disponible
	}
}

// Retorna slice de los bloques libres
func Bloques_Libres(logger *slog.Logger) []byte {

	// Cargar archivo en un slice de bytes
	file, err := os.ReadFile(Configs.MountDir + "/bitmap.dat")
	if err != nil {
		logger.Error(fmt.Sprintf("Error al leer el archivo bitmap.dat: %s\n", err.Error()))
		return nil
	}

	// Identificar los bloques libres
	var bitesLibres []byte
	for i := 0; i < len(file); i++ {
		for j := 0; j < 8; j++ {
			if (file[i] & (1 << j)) == 0 {
				bitesLibres = append(bitesLibres, byte(i*8+j))
			}
		}
	}
	return bitesLibres
}

// Reservo los bloques en el slice
func Reservar_Bloques_Del_Bitmap(bloques []byte, cantidad int, logger *slog.Logger) {

	// Cargar archivo en un slice de bytes
	file, err := os.ReadFile(Configs.MountDir + "/bitmap.dat")
	if err != nil {
		logger.Error(fmt.Sprintf("Error al leer el archivo bitmap.dat: %s\n", err.Error()))
		return
	}

	// Verificar que hay suficientes bloques libres
	bloquesLibres := Bloques_Libres(logger)
	if len(bloquesLibres) < cantidad {
		logger.Error("No hay suficientes bloques libres para reservar")
		return
	}

	// Reservar los bloques
	for i := 0; i < cantidad; i++ {
		bloque := bloquesLibres[i]
		byteIndex := bloque / 8
		bitIndex := bloque % 8
		file[byteIndex] |= (1 << bitIndex) // Marcar el bit como ocupado
	}

	// Guardar los cambios en el archivo
	err = os.WriteFile(Configs.MountDir+"/bitmap.dat", file, 0644)
	if err != nil {
		logger.Error(fmt.Sprintf("Error al escribir el archivo bitmap.dat: %s\n", err.Error()))
		return
	}

	logger.Info(fmt.Sprintf("Reservados %d bloques en el archivo bitmap.dat", cantidad))
}
