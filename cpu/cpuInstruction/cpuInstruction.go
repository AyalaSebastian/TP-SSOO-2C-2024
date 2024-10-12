package cpuInstruction

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/sisoputnfrba/tp-golang/cpu/client"
	"github.com/sisoputnfrba/tp-golang/cpu/server"
)

// Función para asignar el valor a un registro
func AsignarValorRegistro(registro string, valor uint32, logger *slog.Logger) {
	// Obtener una referencia a los registros
	registros := &client.ReceivedContextoEjecucion.Registros

	// Asignar el valor al registro correspondiente
	switch registro {
	case "PC":
		registros.PC = valor
	case "AX":
		registros.AX = valor
	case "BX":
		registros.BX = valor
	case "CX":
		registros.CX = valor
	case "DX":
		registros.DX = valor
	case "EX":
		registros.EX = valor
	case "FX":
		registros.FX = valor
	case "GX":
		registros.GX = valor
	case "HX":
		registros.HX = valor
	case "Base":
		registros.Base = valor
	case "Limite":
		registros.Limite = valor
	default:
		logger.Error(fmt.Sprintf("Registro desconocido: %s", registro))
		return
	}

	// Log de la instrucción ejecutada
	logger.Info(fmt.Sprintf("Instrucción Ejecutada: “## TID: %d - Ejecutando: SET - Registro: %s, Valor: %d”", server.ReceivedPIDTID.TID, registro, valor))
}

// Función para sumar el valor de dos registros
func SumarRegistros(registroDestino, registroOrigen string, logger *slog.Logger) {

	// Obtener los valores de los registros
	valorDestino := obtenerValorRegistro(registroDestino, logger)
	valorOrigen := obtenerValorRegistro(registroOrigen, logger)

	// Sumar los valores
	nuevoValor := valorDestino + valorOrigen

	// Asignar el nuevo valor al registro destino
	AsignarValorRegistro(registroDestino, nuevoValor, logger)

	// Log de la instrucción ejecutada
	logger.Info(fmt.Sprintf("Instrucción Ejecutada: “## TID: %d - Ejecutando: SUM - Registro Destino: %s, Registro Origen: %s”", server.ReceivedPIDTID.TID, registroDestino, registroOrigen))
}

// Función para restar el valor de dos registros
func RestarRegistros(registroDestino, registroOrigen string, logger *slog.Logger) {

	// Obtener los valores de los registros
	valorDestino := obtenerValorRegistro(registroDestino, logger)
	valorOrigen := obtenerValorRegistro(registroOrigen, logger)

	// Restar los valores
	nuevoValor := valorDestino - valorOrigen

	// Asignar el nuevo valor al registro destino
	AsignarValorRegistro(registroDestino, nuevoValor, logger)

	// Log de la instrucción ejecutada
	logger.Info(fmt.Sprintf("Instrucción Ejecutada: “## TID: %d - Ejecutando: SUB - Registro Destino: %s, Registro Origen: %s”", server.ReceivedPIDTID.TID, registroDestino, registroOrigen))
}

// Función para realizar el salto condicional JNZ
func SaltarSiNoCero(registro, instruccion string, logger *slog.Logger) {

	// Obtener el valor del registro
	valorRegistro := obtenerValorRegistro(registro, logger)

	// Si el valor del registro es distinto de cero, actualizar el Program Counter (PC)
	if valorRegistro != 0 {
		// Convertir la instrucción a un valor numérico
		instruccionNueva, err := strconv.ParseUint(instruccion, 10, 32)
		if err != nil {
			logger.Error(fmt.Sprintf("Error al convertir instrucción para JNZ: %s", instruccion))
			return
		}

		// Asignar el nuevo valor del PC
		AsignarValorRegistro("PC", uint32(instruccionNueva), logger)

		// Log de la instrucción ejecutada
		logger.Info(fmt.Sprintf("Instrucción Ejecutada: “## TID: %d - Ejecutando: JNZ - Registro: %s, Nueva Instrucción: %s”", server.ReceivedPIDTID.TID, registro, instruccion))
	}
}

// Función para escribir en el log el valor de un registro
func LogRegistro(registro string, logger *slog.Logger) {
	// Obtener una referencia a los registros
	valor := obtenerValorRegistro(registro, logger)

	// Log de la instrucción ejecutada
	logger.Info(fmt.Sprintf("Instrucción Ejecutada: “## TID: %d - Ejecutando: LOG - Registro: %s, Valor: %d”", server.ReceivedPIDTID.TID, registro, valor))
}

// Función auxiliar para obtener el valor de un registro
func obtenerValorRegistro(registro string, logger *slog.Logger) uint32 {
	registros := &client.ReceivedContextoEjecucion.Registros

	switch registro {
	case "PC":
		return registros.PC
	case "AX":
		return registros.AX
	case "BX":
		return registros.BX
	case "CX":
		return registros.CX
	case "DX":
		return registros.DX
	case "EX":
		return registros.EX
	case "FX":
		return registros.FX
	case "GX":
		return registros.GX
	case "HX":
		return registros.HX
	case "Base":
		return registros.Base
	case "Limite":
		return registros.Limite
	default:
		logger.Error(fmt.Sprintf("Registro desconocido: %s", registro))
		return 0
	}
}
