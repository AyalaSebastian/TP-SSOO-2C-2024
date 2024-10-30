package utils

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/sisoputnfrba/tp-golang/utils/conexiones"
	"github.com/sisoputnfrba/tp-golang/utils/types"
)

func Iniciar_memoria(logger *slog.Logger) {
	mux := http.NewServeMux()

	// Endpoints
	mux.HandleFunc("POST /CREAR-PROCESO", Crear_proceso(logger))
	mux.HandleFunc("PATCH /FINALIZAR-PROCESO/{pid}", Finalizar_proceso(logger))
	mux.HandleFunc("POST /CREAR_HILO", Crear_hilo(logger))
	mux.HandleFunc("POST /FINALIZAR_HILO", Finalizar_hilo(logger))
	mux.HandleFunc("POST /MEMORY-DUMP", Memory_dump(logger))


	// Comunicacion con CPU
	//pasa el contexto de ejecucion a cpu
	mux.HandleFunc("POST /contexto", Obtener_Contexto_De_Ejecucion(logger))

	//envia proxima instr a cpu fase fetch
	mux.HandleFunc("GET /instruccion /{tid}/{pc}", Obtener_Instrucción(logger))

	//recibo msj de cpu para que haga la instruccion read mem
	mux.HandleFunc("POST /read_mem", Read_Mem(logger))

	//recibo msj de cpu para que haga la instruccion write mem
	mux.HandleFunc("POST /write_mem", Write_Mem(logger))

	conexiones.LevantarServidor(strconv.Itoa(Configs.Port), mux, logger)

}

//Coms con KERNEL

func Crear_proceso(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var magic types.PathTamanio
		err := decoder.Decode(&magic)
		if err != nil {
			logger.Error(fmt.Sprintf("Error al decodificar mensaje: %s\n", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error al decodificar mensaje"))
			return
		}
		logger.Info(fmt.Sprintf("Me llegaron los siguientes parametros para crear proceso: %+v", magic))

		// IMPORTANTE: Acá tiene que ir todo para que la memoria CREE el proceso (Está en pagina 20 y 21 del enunciado)

		// Si memoria pudo asignar el espacio necesario para el proceso responde con OK a Kernel
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}

func Finalizar_proceso(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pid := r.PathValue("pid") //Recibimos el pid a finalizar

		logger.Info(fmt.Sprintf("Liberando memoria de Proceso con PID = %+v", pid))

		// IMPORTANTE: Acá tiene que ir todo para que la memoria FINALICE el proceso (Está en pagina 21 del enunciado)

		// Si memoria pudo finalizar el proceso responde con OK a Kernel
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}

}

func Crear_hilo(logger *slog.Logger) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var infoHilo types.EnviarHiloAMemoria
		err := json.NewDecoder(r.Body).Decode(&infoHilo)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Aca va toda la logica para crear el hilo pag(22)

		pidParceado := strconv.Itoa(int(infoHilo.PID))
		logger.Info("## Hilo Creado - (PID:TID) - (%d:%d)", pidParceado, infoHilo.TID)

		// En caso de haberse creado el hilo

		respuesta, err := json.Marshal("OK")
		if err != nil {
			http.Error(w, "Error al codificar los datos como JSON", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(respuesta)

		logger.Info(fmt.Sprintf("## Hilo Creado - (PID:TID) - (%d:%d)", infoHilo.PID, infoHilo.TID))
	}
}

func Finalizar_hilo(logger *slog.Logger) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var infoHilo types.PIDTID
		err := json.NewDecoder(r.Body).Decode(&infoHilo)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		
		// Aca va toda la logica para finalizar el hilo

		// En caso de haberse finalizado el hilo
		respuesta, err := json.Marshal("OK")
		if err != nil {
			http.Error(w, "Error al codificar los datos como JSON", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		//w.Write([]byte("OK"))

		w.Write(respuesta)
		logger.Info(fmt.Sprintf("## Hilo Destruido - (PID:TID) - (%d:%d)", infoHilo.PID, infoHilo.TID))
	}
}

/*
func Memory_dump(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//		decoder := json.NewDecoder(r.Body)
		/*
			if err != nil {
				logger.Error(fmt.Sprintf("Error al decodificar mensaje: %s\n", err.Error()))
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Error al decodificar mensaje"))
				return
			}
		
        var req a
        err := json.NewDecoder(r.Body).Decode(&req)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}


        w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}
*/

// Comunicacion con CPU

func Actualizar_Contexto(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
        var req
        err := json.NewDecoder(r.Body).Decode(&req)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
        Actualizar_Contexto(req.PID, req.TID)
        
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
	
	}
}

//pasa el contexto de ejecucion a cpu
func Enviar_proxima_instruccion(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var tid types.ContextoEjecucionTID
		tid := r.PathValue("tid")
		pc := r.PathValue("pc")
		var req types.RegCPU
        err := json.NewDecoder(r.Body).Decode(&req)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
        regCPU, err := Ver_Contexto(req.PID, req.TID)
        if err != nil {
            http.Error(w, err.Error(), http.StatusNotFound)
            return
        }
        ipCPU := Config.IpCPU //Cambiar por referencia a archivo config
        puertoCPU := Config.PortCPU 


        exito := client.Enviar_Body(regCPU, conf, puertoCPU, endpointCPU, logger)
        if !exito {
            http.Error(w, "Error al enviar el contexto al CPU", http.StatusInternalServerError)
            return
        }
        
        w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))

	}
}

func Obtener_ID(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
        var req
        err := json.NewDecoder(r.Body).Decode(&req)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
        w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))

	}
}

func Read_Mem(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
        var req
        err := json.NewDecoder(r.Body).Decode(&req)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

    
        w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))

	}
}

//recibo 4 bytes y los escribo a partir del byte enviado como direccion fisica 
//dentro de la Memoria de Usuario y se responderá como OK.

// Función Write_Mem para manejar la escritura en la memoria a partir de la API
func Write_Mem(w http.ResponseWriter, r *http.Request, logger *slog.Logger) {
    // Estructura para recibir los datos del JSON
    var requestData struct {
        DireccionFisica uint32 `json:"direccion_fisica"`
        Valor           uint32 `json:"valor"`
        TID             uint32 `json:"tid"`
    }

    // Decodificar el JSON del cuerpo de la solicitud
    err := json.NewDecoder(r.Body).Decode(&requestData)
    if err != nil {
        logger.Error("Error al decodificar JSON en Write_Mem", slog.Any("error", err))
        http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
        return
    }

    // Validar si la dirección física está dentro de los límites de la memoria
    if requestData.DireccionFisica > uint32(len(memoria)-4) { // -4 para no salir del límite al escribir 4 bytes
        logger.Error("Dirección física fuera de rango")
        http.Error(w, "Dirección física fuera de rango", http.StatusBadRequest)
        return
    }

    // Escribir el valor en little-endian en la memoria
    binary.LittleEndian.PutUint32(memoria[requestData.DireccionFisica:], requestData.Valor)

    // Log de escritura en memoria de usuario
    logger.Info(fmt.Sprintf("Escritura / lectura en espacio de usuario: ## Escritura - (PID:TID) - (N/A:%d) - Dir. Física: %d - Tamaño: %d",
        requestData.TID, requestData.DireccionFisica, 4))

    // Confirmar que la operación fue exitosa
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))

    // Log de escritura exitosa
    logger.Info(fmt.Sprintf("Escritura en memoria exitosa: TID %d - Dirección Física: %d - Valor: %d",
        requestData.TID, requestData.DireccionFisica, requestData.Valor))
}
