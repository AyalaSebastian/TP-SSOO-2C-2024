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

//////////////////////////////////////////////////////////////////////////////////////////////
///////////////////               FETCH INSTRUCCIONES               /////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////

// Variable global para almacenar la instrucción obtenida
var Instruccion string

// Función Fetch para obtener la próxima instrucción
func Fetch(ipMemory string, portMemory int, tid uint32, logger *slog.Logger) error {
	if ReceivedContextoEjecucion == nil {
		logger.Error("No se ha recibido el contexto de ejecución. Imposible realizar Fetch.")
		return fmt.Errorf("contexto de ejecución no disponible")
	}

	// Obtener el valor del PC (Program Counter) de la variable global
	pc := ReceivedContextoEjecucion.Registros.PC

	// Crear la estructura de solicitud
	requestData := struct {
		PC  uint32 `json:"pc"`
		TID uint32 `json:"tid"`
	}{PC: pc, TID: tid}

	// Serializar los datos en JSON
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		logger.Error("Error al codificar PC y TID a JSON: ", slog.Any("error", err))
		return err
	}

	// Crear la URL del módulo de Memoria
	url := fmt.Sprintf("http://%s:%d/instruccion", ipMemory, portMemory)

	// Crear la solicitud POST
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Error("Error al crear la solicitud: ", slog.Any("error", err))
		return err
	}

	// Establecer el encabezado de la solicitud
	req.Header.Set("Content-Type", "application/json")

	// Enviar la solicitud
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("Error al enviar la solicitud de Fetch: ", slog.Any("error", err))
		return err
	}
	defer resp.Body.Close()

	// Verificar si la respuesta fue exitosa
	if resp.StatusCode != http.StatusOK {
		logger.Error(fmt.Sprintf("Error en la respuesta de Fetch: Código de estado %d", resp.StatusCode))
		return fmt.Errorf("error en la respuesta de Fetch: Código de estado %d", resp.StatusCode)
	}

	// Decodificar la respuesta para obtener la instrucción
	var fetchedInstruction struct {
		Instruccion string `json:"instruccion"`
	}
	err = json.NewDecoder(resp.Body).Decode(&fetchedInstruction)
	if err != nil {
		logger.Error("Error al decodificar la instrucción recibida: ", slog.Any("error", err))
		return err
	}

	// Guardar la instrucción en la variable global
	Instruccion = fetchedInstruction.Instruccion

	// Log de Fetch exitoso
	logger.Info(fmt.Sprintf("Fetch Instrucción: “## TID: %d - FETCH - Program Counter: %d”", tid, pc))

	return nil
}
func DevolverTIDAlKernel(tid uint32, logger slog.Logger) bool {
	cliente := &http.Client{}
	url := fmt.Sprintf("http://%s:%d/%s/%v", "127.0.0.1", 8001, "THREAD_JOIN", tid)
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
