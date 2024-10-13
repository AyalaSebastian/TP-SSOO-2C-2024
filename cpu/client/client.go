package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/sisoputnfrba/tp-golang/utils/types"
)

// Variable global para almacenar el contexto de ejecución
var ReceivedContextoEjecucion *types.ContextoEjecucion = nil

// Función para obtener el contexto de ejecución
//func GetContextoEjecucion() *types.ContextoEjecucion {
//	return contextosEjecucion
//}

// Función que solicita el contexto de ejecución al módulo de memoria
func SolicitarContextoEjecucion(ipMemory string, portMemory int, pid uint32, tid uint32, logger *slog.Logger) error {
	url := fmt.Sprintf("http://%s:%d/contexto", ipMemory, portMemory) // URL del módulo de memoria

	// Crear la estructura PIDTID con los valores recibidos
	pidTid := struct {
		TID uint32 `json:"tid"`
		PID uint32 `json:"pid"`
	}{TID: tid, PID: pid}

	// Codificarla en JSON
	jsonData, err := json.Marshal(pidTid)
	if err != nil {
		logger.Error("Error al codificar TID y PID a JSON: ", slog.Any("error", err))
		return err
	}

	// Crear la solicitud POST
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Error("Error al crear la solicitud: ", slog.Any("error", err))
		return err
	}

	// Establecer los encabezados
	req.Header.Set("Content-Type", "application/json")

	// Enviar la solicitud
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("Error al enviar la solicitud al módulo de memoria: ", slog.Any("error", err))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error(fmt.Sprintf("Error en la respuesta del módulo de memoria: Código de estado %d", resp.StatusCode))
		return fmt.Errorf("error en la respuesta del módulo de memoria: Código de estado %d", resp.StatusCode)
	}

	// Decodificar la respuesta
	var contexto types.ContextoEjecucion
	err = json.NewDecoder(resp.Body).Decode(&contexto)
	if err != nil {
		logger.Error("Error al decodificar el contexto de ejecución: ", slog.Any("error", err))
		return err
	}

	// Asignar el contexto recibido a la variable global
	ReceivedContextoEjecucion = &contexto
	logger.Info("Contexto de ejecución recibido con éxito")

	return nil
}

func DevolverTIDAlKernel(tid uint32, logger slog.Logger,endpoint string , motivo string) bool {
	cliente := &http.Client{}
	url := fmt.Sprintf("http://%s:%d/%s/%v", "127.0.0.1", 8001, endpoint,(tid,motivo))
	req, err := http.NewRequest("PATCH", url, nil)
	if err != nil {
		return false
	}
	// Establecer el encabezado Content-Type
	req.Header.Set("Content-Type", "application/json")

	// Enviar la solicitud
	resp, err := cliente.Do(req)
	if err != nil {
		logger.Error(fmt.Sprintf("Error al enviar la solicitud: %v", err))
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error("La respuesta del servidor no fue OK")
		return false // Indica que la respuesta no fue exitosa
	}

	return true
}

func ActualizarContextoDeEjecucion(tid uint32, Logger slog.Logger) {

}

/*
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
*/
