package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/sisoputnfrba/tp-golang/kernel/utils"
	"github.com/sisoputnfrba/tp-golang/utils/client"
	"github.com/sisoputnfrba/tp-golang/utils/generadores"
	"github.com/sisoputnfrba/tp-golang/utils/logging"
)

func main() {
	// Inicializamos la configuracion y el logger
	utils.Configs = utils.Iniciar_Configuracion("config.json")
	logger := logging.Iniciar_Logger("kernel.log", utils.Configs.LogLevel)

	// Solo lo puse para probar el logger, después lo sacamos
	logger.Info("Hola! Soy el módulo kernel")

	// Enviamos handshake a cpu
	client.Enviar_handshake(utils.Configs.IpCPU, utils.Configs.PortCPU, "Estableciendo handshake con CPU desde Kernel")
	client.Enviar_handshake(utils.Configs.IpMemory, utils.Configs.PortMemory, "Estableciendo handshake con Memoria desde Kernel")

	// Obtener los parametros del primer proceso a ejecutar
	archivoPseudocodigo := os.Args[1]
	tamanioProceso, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Error: El tamaño del proceso debe ser un número entero.")
		panic(err)
	}

	// Creación del proceso inicial
	planificar_proceso_inicial(archivoPseudocodigo, tamanioProceso)

	// Iniciamos Kernel como server
	utils.Iniciar_kernel(logger)

}

// Acá para mi hay que mandar el path a memoria para que saque las instrucciones del archivo de pseudocódigo y acá mismo armar el PCB con el TCB y todo
func planificar_proceso_inicial(pseudo string, tamanio int) {

	pcb := generadores.Generar_PCB()
	tcb := generadores.Generar_TCB(&pcb, 0)

	// Falta enviar el pseudo y el tamaño a memoria para que inicialice el proceso

}
