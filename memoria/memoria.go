package main

import (
	"github.com/sisoputnfrba/tp-golang/memoria/memUsuario"
	"github.com/sisoputnfrba/tp-golang/memoria/server"
	"github.com/sisoputnfrba/tp-golang/memoria/utils"
	"github.com/sisoputnfrba/tp-golang/utils/logging"
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
