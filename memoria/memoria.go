package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/sisoputnfrba/tp-golang/memoria/utils"
	"github.com/sisoputnfrba/tp-golang/utils/logging"
)

func main() {
	// Inicializamos la configuracion y el logger
	utils.Configs = utils.Iniciar_configuracion("config.json")
	logger := logging.Iniciar_Logger("memoria.log", utils.Configs.LogLevel)
	var instrucciones = []

	// Inicializamos la memoria (Lo levantamos como servidor)
	utils.Iniciar_memoria(logger)

	memory := Inicializar_Memoria()

	leerArchivoPseudoCodigo(tid,utils.Configs.InstructionPath,pc,instrucciones)

}

func Inicializar_Memoria() []byte {
	return make([]byte, utils.Configs.MemorySize)
}

// Estructura para representar una partición de memoria
type Particion struct {
	Base   uint32
	Limite uint32
}

// Definición de las particiones fijas
var particiones = []Particion{
	{Base: 0, Limite: 512},    // Primera partición: del byte 0 al byte 511
	{Base: 512, Limite: 16},   // Segunda partición: del byte 512 al 527
	{Base: 528, Limite: 32},   // Tercera partición: del byte 528 al 559
	{Base: 560, Limite: 16},   // Cuarta partición: del byte 560 al 575
	{Base: 576, Limite: 256},  // Quinta partición: del byte 576 al 831
	{Base: 832, Limite: 64},   // Sexta partición: del byte 832 al 895
	{Base: 896, Limite: 128},  // Séptima partición: del byte 896 al 1023
}


func Nuevo_Hilo(tid int, base uint32, limite uint32) {
	registros[tid] = &RegCPU{
		PC:     0,
		AX:     0,
		BX:     0,
		CX:     0,
		DX:     0,
		EX:     0,
		FX:     0,
		GX:     0,
		HX:     0,
		Base:   base,
		Limite: limite,
	}
}

func Ver_Contexto(pid int, tid int) (*RegCPU, error) {
	regCPU, existe := cpuRegisters[tid]
	if !existe {
		return nil, logger.Error(fmt.Sprintf("El hilo con TID %d no existe para el proceso PID %d", tid, pid))
	}
	logger.Info(fmt.Sprintf("## Contexto solicitado (%d : %d)", pid, tid))
	// Devolver el contexto completo
	return regCPU, nil

}

func Update_Contexto(pid uint32, tid uint32) {
	regCPU, existe := cpuRegisters[tid]
	if !existe {
		return nil, logger.Error(fmt.Sprintf("El hilo con TID %d no existe para el proceso PID %d", tid, pid))
	}
	regCPU.AX = RegCPU.AX
	regCPU.BX = req.RegCPU.BX
	regCPU.CX = req.RegCPU.CX
	regCPU.DX = req.RegCPU.DX
	regCPU.EX = req.RegCPU.EX
	regCPU.FX = req.RegCPU.FX
	regCPU.GX = req.RegCPU.GX
	regCPU.HX = req.RegCPU.HX
	regCPU.PC = req.RegCPU.PC
	regCPU.Base = req.RegCPU.Base
	regCPU.Limite = req.RegCPU.Limite

	logger.Info(fmt.Sprintf("## Contexto actualizado (%d : %d)", pid, tid))
	return
}

func leerArchivoPseudoCodigo(tid uint32,pc int,archivo string,lista_instrucciones [string]) [string]{
	file, err := os.Open(archivo)
	if err != nil {
		fmt.Println("Error al abrir el archivo:", err)
		return
	}

	// Crear un scanner para leer línea por línea
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		instruccion := scanner.Text() // Obtener la línea actual
		instrucciones = append(lista_instrucciones, instruccion)     // Hacer algo con la línea (imprimirla en este caso)
	}

	// Manejar posibles errores durante la lectura
	if err := scanner.Err(); err != nil {
		fmt.Println("Error al leer el archivo:", err)
	}
	defer file.Close()
	Client.Obtener_Instrucción(tid,lista_instrucciones,pc)
}

//PREGUNTAS PARA SOPORTE:
//COMO REALIZAR EL ARCHIVO DE PSEUDOCODIGO
//ARREGLAR FUNCION INICIAR MEMORIA
// que seria : partitions": [512, 16, 32, 16, 256, 64, 128],
