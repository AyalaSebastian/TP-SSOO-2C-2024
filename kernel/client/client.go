package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/sisoputnfrba/tp-golang/utils/types"
)

// Devuelve true si memoria pudo crear el proceso, false en caso contrario
func Enviar_parametros_proceso(ip string, puerto int, path string, tamanio int, logger *slog.Logger) bool {
	mensaje := types.PathTamanio{Path: path, Tamanio: tamanio}
	body, err := json.Marshal(mensaje)
	if err != nil {
		logger.Error("Se produjo un error codificando el mensaje")
		return false
	}
	// Siento que esto lo podriamos modularizar en una funcion que reciba el ip, puerto, body y el endpoint
	url := fmt.Sprintf("http://%s:%d/crear-proceso", ip, puerto)
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

func Informar_memoria_creacion_hilo() {

}

func Enviar_server[T any](dato T, ip string, puerto int, endpoint string, logger *slog.Logger) bool {

	body, err := json.Marshal(dato)
	if err != nil {
		logger.Error("Se produjo un error codificando el mensaje")
		return false
	}

	url := fmt.Sprintf("http://%s:%d/%s", ip, puerto, endpoint)
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
