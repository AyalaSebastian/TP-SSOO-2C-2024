package main

import (
	"fmt"
	"log/slog"
    "os"
    "bufio"
	"github.com/sisoputnfrba/tp-golang/memoria/memUsuario"
	"github.com/sisoputnfrba/tp-golang/memoria/server"
	"github.com/sisoputnfrba/tp-golang/memoria/utils"
	"github.com/sisoputnfrba/tp-golang/utils/logging"
	"github.com/sisoputnfrba/tp-golang/utils/types"
)

func main() {
	// Inicializamos la configuracion y el logger
	utils.Configs = utils.Iniciar_configuracion("config.json")
	logger := logging.Iniciar_Logger("memoria.log", utils.Configs.LogLevel)
	//var instrucciones []string

	// Inicializamos la memoria (Lo levantamos como servidor)
	server.Iniciar_memoria(logger)

	memUsuario.Inicializar_Memoria_De_Usuario()

	//leerArchivoPseudoCodigo(tid, utils.Configs.InstructionPath, pc, instrucciones)

}

func Nuevo_Hilo(tid int, base uint32, limite uint32) {

	var registros = make(map[int]*types.RegCPU)

	registros[tid] = &types.RegCPU{
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

func Ver_Contexto(pid int, tid int, logger *slog.Logger) (*types.RegCPU, error) {

	var cpuRegisters = make(map[int]*types.RegCPU)

	regCPU, existe := cpuRegisters[tid]
	if !existe {
		return nil, fmt.Errorf("El hilo con TID %d no existe para el proceso PID %d", tid, pid)
	}
	logger.Info(fmt.Sprintf("## Contexto solicitado (%d : %d)", pid, tid))
	// Devolver el contexto completo
	return regCPU, nil

}















/*
func Update_Contexto(pid uint32, tid uint32, logger *slog.Logger) {

	var cpuRegisters = make(map[int]*types.RegCPU)

	regCPU, existe := cpuRegisters[tid]
	if !existe {
		logger.Error(fmt.Sprintf("El hilo con TID %d no existe para el proceso PID %d", tid, pid))
		return
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
*/

/*
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
*/


/*PREGUNTAS PARA SOPORTE:
de donde vienen el pid y tid que recibe memoria para:
Se deberá almacenar, por cada PID del sistema, la parte del contexto de ejecución común para el proceso, en este caso, es la requerida para poder traducir las direcciones lógicas a físicas: base y límite.
Luego, por cada TID del sistema, se tendrán los registros de la CPU propios de cada hilo: AX, BX, CX, DX, EX, FX, GX, HX y PC. Siendo todos ellos inicializados en 0.
*/
