package memUsuario

import (
	"net/http"

	"github.com/sisoputnfrba/tp-golang/memoria/utils"
)

func CrearProceso(pid uint32) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		particiones := utils.Config.Partitions
		cant_particiones := len(particiones)
		for i := 0; i < cant_particiones; i++ {
			if particiones[i] != (-1) {
				particiones[i] = -1
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
				return
			}
		}
		(http.Error(w, "NO SE PUDO INICIALIZAR EL PROCESO POR FALTA DE ESPACIO LIBRE EN LA MEMORIA", http.StatusInternalServerError))
		return
	}
}
