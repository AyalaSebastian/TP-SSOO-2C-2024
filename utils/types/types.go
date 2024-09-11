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
	PC uint32 `json:"pc"` // Program Counter (Proxima instruccion a ejecutar)
	AX uint32 `json:"ax"`
	BX uint32 `json:"bx"`
	CX uint32 `json:"cx"`
	DX uint32 `json:"dx"`
	EX uint32 `json:"ex"`
	FX uint32 `json:"fx"`
	GX uint32 `json:"gx"`
	HX uint32 `json:"hx"`
}

type ContextoEjecucion struct {
	Base      int    `json:"base"`
	Limite    int    `json:"limite"`
	Registros RegCPU `json:"registros"`
}
