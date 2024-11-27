package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/sisoputnfrba/tp-golang/cpu/utils"
	"github.com/sisoputnfrba/tp-golang/utils/types"
)

// Variable global para almacenar el contexto de ejecución
var ReceivedContextoEjecucion *types.ContextoEjecucion = nil

var Proceso types.Proceso

// Función que solicita el contexto de ejecución al módulo de memoria
func SolicitarContextoEjecucion(pidTid types.PIDTID, logger *slog.Logger) error {
	url := fmt.Sprintf("http://%s:%d/contexto", utils.Configs.IpMemory, utils.Configs.PortMemory) // URL del módulo de memoria
	logger.Info(fmt.Sprintf("## TID: %d - Solicito Contexto Ejecución", pidTid.TID))
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
	Proceso.ContextoEjecucion = contexto
	logger.Info("Contexto de ejecución recibido con éxito")

	return nil
}

func DevolverTIDAlKernel(tid uint32, logger *slog.Logger, endpoint string, motivo string) bool {
	cliente := &http.Client{}
	url := fmt.Sprintf("http://%s:%d/%s/%v", utils.Configs.IpKernel, utils.Configs.PortKernel, endpoint, tid)
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

func EnviarContextoDeEjecucion[T any](dato T, endpoint string, logger *slog.Logger) bool {

	body, err := json.Marshal(dato)
	if err != nil {
		logger.Error("Se produjo un error codificando el mensaje")
		return false
	}
	//ipMemory y portMemory
	url := fmt.Sprintf("http://%s:%d/%s", utils.Configs.IpMemory, utils.Configs.PortMemory, endpoint)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		logger.Error(fmt.Sprintf("Se produjo un error enviando mensaje a ip:%s puerto:%d", "127.0.0.1", 8002))
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

func CederControlAKernell[T any](dato T, endpoint string, logger *slog.Logger) {

	body, err := json.Marshal(dato)
	if err != nil {
		logger.Error("Se produjo un error codificando el mensaje")
		return
	}

	url := fmt.Sprintf("http://%s:%d/%s", utils.Configs.IpKernel, utils.Configs.PortKernel, endpoint)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		logger.Error(fmt.Sprintf("Se produjo un error enviando mensaje a ip:%s puerto:%d", "127.0.0.1", 8001))
		return
	}
	// Aseguramos que el body sea cerrado
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error("La respuesta del servidor no fue OK")
		return // Indica que la respuesta no fue exitosa
	}
	os.Exit(0)
}

// EnviarDesalojo envia el PID, TID y el motivo del desalojo a la API Kernel utilizando la configuración global de IP y puerto.
func EnviarDesalojo(pid uint32, tid uint32, motivo string, logger *slog.Logger) {

	// Crear el objeto que contiene los datos a enviar
	hiloDesalojado := types.HiloDesalojado{
		PID:    pid,
		TID:    tid,
		Motivo: motivo,
	}

	// Convertir el objeto a JSON
	body, err := json.Marshal(hiloDesalojado)
	if err != nil {
		logger.Error("Error al codificar mensaje de desalojo", slog.String("error", err.Error()))
		return
	}

	// Formar la URL de la API Kernel usando las configuraciones globales
	url := fmt.Sprintf("http://%s:%d/desalojo", utils.Configs.IpKernel, utils.Configs.PortKernel)

	// Enviar la solicitud POST
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		logger.Error(fmt.Sprintf("Error al enviar desalojo a %s:%d", utils.Configs.IpKernel, utils.Configs.PortKernel), slog.String("error", err.Error()))
		return
	}
	defer resp.Body.Close()

	// Verificar que la respuesta sea exitosa
	if resp.StatusCode != http.StatusOK {
		logger.Error("Error al procesar la solicitud de desalojo", slog.Int("status", resp.StatusCode))
		return
	}

	// Log de éxito
	logger.Info("Desalojo enviado correctamente", slog.Int("PID", int(pid)), slog.Int("TID", int(tid)), slog.String("Motivo", motivo))
}
