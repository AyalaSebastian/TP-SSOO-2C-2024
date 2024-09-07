package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/sisoputnfrba/tp-golang/utils/types"
)

// le hacemos un clousure a la funcion para que reciba el logger
func Recibir_handshake(logger *slog.Logger) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var mensaje types.HandShake
		err := decoder.Decode(&mensaje)
		if err != nil {
			logger.Info(fmt.Sprintf("Error al decodificar mensaje: %s", err.Error()))

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error al decodificar mensaje / Con el Handshake"))
			return
		}

		logger.Info("Se pudo establecer la conexion, siguiendo con la funcion")
		logger.Info(fmt.Sprintf("%+v", mensaje))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("HANDSHAKE OK"))
	}
}

//func RecibirHandshake(w http.ResponseWriter, r *http.Request) {}
