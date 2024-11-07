package utils

import (
	"sync"
)

var MutexPlanificador sync.Mutex
var Planificador = sync.NewCond(&MutexPlanificador)
