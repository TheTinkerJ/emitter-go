package controller

type ClientChan chan string

type Event struct {
	Message       ClientChan
	NewClients    chan ClientChan
	ClosedClients chan ClientChan
	TotalClients  map[ClientChan]bool
}

type EmitterController struct {
}

func (ec *EmitterController) StartEmitter() {
}
