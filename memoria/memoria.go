package main

import (
	"github.com/sisoputnfrba/tp-golang/memoria/utils"
	//"github.com/sisoputnfrba/tp-golang/utils/client"
	"github.com/sisoputnfrba/tp-golang/utils/logging"
)

func main() {
	// Inicializamos la configuracion y el logger
	utils.Configs = utils.Iniciar_configuracion("config.json")
	logger := logging.Iniciar_Logger("memoria.log", utils.Configs.LogLevel)

	// Solo lo puse para probar el logger, después lo sacamos
	logger.Info("Hola! Soy el módulo memoria")

	//client.Enviar_handshake(utils.Configs.IpFilesystem, utils.Configs.PortFilesystem, "Estableciendo handshake con Filesystem desde Memoria")

	// Inicializamos la memoria (Lo levantamos como servidor)
	utils.Iniciar_memoria(logger)

}
