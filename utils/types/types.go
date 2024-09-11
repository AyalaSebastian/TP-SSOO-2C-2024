package types

type HandShake struct {
	Mensaje string `json:"mensaje"`
}

type PCB struct {
	PID    uint32   `json:"pid"`
	TCBs   []TCB    `json:"tcb"`
	Mutexs []string `json:"mutexs"`
}

type TCB struct {
	TID       uint32 `json:"tid"`
	Prioridad string `json:"prioridad"`
	Estado    string `json:"estado"` //Puede ser "NEW", "READY", "EXECUTE", "BLOCKED", "EXIT" (En mayusculas)
	PID       uint32 `json:"pid"`    //PID del proceso al que pertenece

}
