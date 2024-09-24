package utils

// Uso [T any] para poder encolar y desencolar tanto PCBs como TCBs

func Encolar[T any](cola *[]T, elemento T) {
	*cola = append(*cola, elemento)
}

func Desencolar[T any](cola *[]T) T {
	if len(*cola) == 0 {
		var vacio T
		return vacio // O manejar el caso de cola vacía
	}
	elemento := (*cola)[0]
	*cola = (*cola)[1:] // Elimina el primer elemento
	return elemento
}
