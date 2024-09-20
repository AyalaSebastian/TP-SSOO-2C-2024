package types

type HandShake struct {
	Mensaje string `json:"mensaje"`
}

type PCB struct {
	PID    uint32   `json:"pid"`
	TCBs   []TCB    `json:"tcb"`
	Mutexs []string `json:"mutexs"` // los podemos representar con 0 y 1 (como el bitarray)
}

type TCB struct {
	TID       uint32 `json:"tid"`
	Prioridad string `json:"prioridad"`
	Estado    string `json:"estado"`    //Puede ser "NEW", "READY", "EXECUTE", "BLOCKED", "EXIT" (En mayusculas)
	PID       uint32 `json:"pid"`       //PID del proceso al que pertenece
	Registros RegCPU `json:"registros"` //! (Verificar si esta bien)
}

type RegCPU struct {
	PC     uint32 `json:"pc"`     // Program Counter (Proxima instruccion a ejecutar)
	AX     uint32 `json:"ax"`     // Registro Numerico de proposito general
	BX     uint32 `json:"bx"`     // Registro Numerico de proposito general
	CX     uint32 `json:"cx"`     // Registro Numerico de proposito general
	DX     uint32 `json:"dx"`     // Registro Numerico de proposito general
	EX     uint32 `json:"ex"`     // Registro Numerico de proposito general
	FX     uint32 `json:"fx"`     // Registro Numerico de proposito general
	GX     uint32 `json:"gx"`     // Registro Numerico de proposito general
	HX     uint32 `json:"hx"`     // Registro Numerico de proposito general
	Base   uint32 `json:"base"`   // Direccion base de la particion del proceso
	Limite uint32 `json:"limite"` // Tamanio de la particion del proceso
}

type ContextoEjecucion struct {
	Registros RegCPU `json:"registros"`
}

type Particion struct {
	Registros RegCPU `json:"registros"`
}

type CPU struct {
	Contexto          ContextoEjecucion `json:"contexto"`
	MMU               MMU               `json:"mmu"`
	Memoria           Memoria           `json:"memoria"`
	InstruccionActual string            `json:"instruccion_actual"`
}
type MMU struct {
}
type Memoria struct {
}
