package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"github.com/sisoputnfrba/tp-golang/utils/conexiones"
	"github.com/sisoputnfrba/tp-golang/utils/types"
	
)



func Enviar_Body[T any](dato T, ip string, puerto int, endpoint string, logger *slog.Logger) bool {

	body, err := json.Marshal(dato)
	if err != nil {
		logger.Error("Se produjo un error codificando el mensaje")
		return false
	}

	url := fmt.Sprintf("http://%s:%d/%s", ip, puerto, endpoint)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		logger.Error(fmt.Sprintf("Se produjo un error enviando mensaje a ip:%s puerto:%d", ip, puerto))
		return false
	}
	// Aseguramos que el body sea cerrado
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error("La respuesta del servidor no fue OK")
		return false // Indica que la respuesta no fue exitosa
	}

	return true // Indica que la respuesta fue exitosa
}
