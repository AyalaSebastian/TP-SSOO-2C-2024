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

// Abro el archivo y vuelco el contenido en un slice de bytes
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

// 1. que es lo que recebimos en el dump?
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
