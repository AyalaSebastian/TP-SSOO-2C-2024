package utils

import (
	"log/slog"
	"math"
	"os"
)

func Inicializar_Estructura_Filesystem(logger *slog.Logger) {

	// MNT_DIR y MOUNT_DIR y FILES

	// Si no existe MOUNT_DIR lo creo y creo el bloques.dat
	if !Verificar_Si_Existe(Configs.MountDir) {
		err := os.Mkdir(Configs.MountDir, 0755) // Creo el dir
		if err != nil {
			panic("Error al crear el directorio MOUNT_DIR")
		}

		filee, errrrr := os.Create(Configs.MountDir + "/bitmap.dat") // Creo el bloques.dat
		if errrrr != nil {
			panic("Error al crear el archivo de bloques")
		}
		defer filee.Close()

		num := float64(Configs.BlockCount) / 8
		numRedondeado := math.Ceil(num)

		errrrrr := os.Truncate(Configs.MountDir+"/bitmap.dat", int64(numRedondeado)) // Cambio tamanio
		if errrrrr != nil {
			panic("Error al truncar el archivo de bloques")
		}

		file, err := os.Create(Configs.MountDir + "/bloques.dat") // Creo el bloques.dat
		if err != nil {
			panic("Error al crear el archivo de bloques")
		}
		defer file.Close()

		errr := os.Truncate(Configs.MountDir+"/bloques.dat", int64(Configs.BlockSize*Configs.BlockCount)) // Cambio tamanio
		if errr != nil {
			panic("Error al truncar el archivo de bloques")
		}

		errrr := os.Mkdir(Configs.MountDir+"/files", 0755) // Creo el dir files
		if errrr != nil {
			panic("Error al crear el directorio FILES")
		}
		return
	}

	// Si existe verifico el bloques.dat
	if Verificar_Si_Existe(Configs.MountDir) {
		if !Verificar_Si_Existe(Configs.MountDir + "/bloques.dat") {
			file, err := os.Create(Configs.MountDir + "/bloques.dat") // Creo el bloques.dat
			if err != nil {
				panic("Error al crear el archivo de bloques")
			}
			defer file.Close()

			errr := os.Truncate(Configs.MountDir+"/bloques.dat", int64(Configs.BlockSize*Configs.BlockCount)) // Cambio tamanio
			if errr != nil {
				panic("Error al truncar el archivo de bloques")
			}
		}

		if !Verificar_Si_Existe(Configs.MountDir + "/files") {
			errrr := os.Mkdir(Configs.MountDir+"/files", 0755) // Creo el dir files
			if errrr != nil {
				panic("Error al crear el directorio FILES")
			}
		}

		if !Verificar_Si_Existe(Configs.MountDir + "/bitmap.dat") {
			filee, err := os.Create(Configs.MountDir + "/bitmap.dat") // Creo el bloques.dat
			if err != nil {
				panic("Error al crear el archivo de bloques")
			}
			defer filee.Close()

			num := float64(Configs.BlockCount) / 8
			numRedondeado := math.Ceil(num)

			errr := os.Truncate(Configs.MountDir+"/bitmap.dat", int64(numRedondeado)) // Cambio tamanio
			if errr != nil {
				panic("Error al truncar el archivo de bloques")
			}
		}
		return
	}

	logger.Info("Filesystem inicializado")
}

func Verificar_Si_Existe(Path string) bool {

	_, err := os.Stat(Path)

	return !os.IsNotExist(err)
}
