package utils

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	IpMemory   string `json:"ip_memory"`
	PortMemory int    `json:"port_memory"`
	IpKernel   string `json:"ip_kernel"`
	PortKernel int    `json:"port_kernel"`
	Port       int    `json:"port"`
	LogLevel   string `json:"log_level"`
}

var Configs Config // Variable global dentro del package

func Iniciar_configuracion(filePath string) Config {

	configFile, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&Configs)

	return Configs
}

type RegCPU struct {
	PC     uint32 `json:"pc"`     // Program Counter (Proxima instruccion a ejecutar)
	AX     uint32 `json:"ax"`     // Registro Numérico de propósito general
	BX     uint32 `json:"bx"`     // Registro Numérico de propósito general
	CX     uint32 `json:"cx"`     // Registro Numérico de propósito general
	DX     uint32 `json:"dx"`     // Registro Numérico de propósito general
	EX     uint32 `json:"ex"`     // Registro Numérico de propósito general
	FX     uint32 `json:"fx"`     // Registro Numérico de propósito general
	GX     uint32 `json:"gx"`     // Registro Numérico de propósito general
	HX     uint32 `json:"hx"`     // Registro Numérico de propósito general
	Base   uint32 `json:"base"`   // Dirección base de la partición del proceso
	Limite uint32 `json:"limite"` // Tamaño de la partición del proceso
}

type ContextoEjecucion struct {
	Registros RegCPU `json:"registros"`
}

type Particion struct {
	Registros RegCPU `json:"registros"`
}

type CPU struct {
	Contexto          ContextoEjecucion `json:"contexto"`
	MMU               MMU               `json:"mmu"`
	Memoria           Memoria           `json:"memoria"`
	InstruccionActual string            `json:"instruccion_actual"`
}
type MMU struct {
}
type Memoria struct {
}
