package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/sisoputnfrba/tp-golang/utils/types"
)

func Enviar_parametros_contexto(ip string, puerto int, path string, tamanio int, logger *slog.Logger) bool {
	mensaje := types.PathTamanio{Path: path, Tamanio: tamanio}
	body, err := json.Marshal(mensaje)
	if err != nil {
		logger.Error("Se produjo un error codificando el mensaje")
		return false
	}
	// Siento que esto lo podriamos modularizar en una funcion que reciba el ip, puerto, body y el endpoint
	url := fmt.Sprintf("http://%s:%d/crear-contexto", ip, puerto)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		logger.Info("Se produjo un error enviando mensaje a ip:%s puerto:%d", ip, puerto)
		return false
	}
	// Aseguramos que el body sea cerrado
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Info("La respuesta del servidor no fue OK")
		return false // Indica que la respuesta no fue exitosa
	}

	logger.Info(fmt.Sprintf("Respuesta del servidor: %s", resp.Status))
	return true // Indica que la respuesta fue exitosa
}
