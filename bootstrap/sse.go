package bootstrap

import (
	"sync"

	"github.com/gin-gonic/gin"
)

type StringMsgChan chan string

type SSEServer struct {
	SrvMsgChan     StringMsgChan
	NewComeClients chan StringMsgChan
	ClosedClients  chan StringMsgChan
	ClientConnMap  map[StringMsgChan]bool
	rootWg         *sync.WaitGroup // rootwg
	stopChan       chan struct{}
}

func NewSSEServer(rootWg *sync.WaitGroup) *SSEServer {
	sseServer := &SSEServer{
		SrvMsgChan:     make(StringMsgChan),
		NewComeClients: make(chan StringMsgChan),
		ClosedClients:  make(chan StringMsgChan),
		ClientConnMap:  make(map[StringMsgChan]bool),
		rootWg:         rootWg,
		stopChan:       make(chan struct{}, 1),
	}
	return sseServer
}

func (sseServer *SSEServer) Listening() {
	defer sseServer.rootWg.Done()
ListeningLoop:
	for {
		select {
		case _, ok := <-sseServer.stopChan:
			if !ok {
				break ListeningLoop
			}
		case client := <-sseServer.NewComeClients:
			sseServer.ClientConnMap[client] = true
		case client := <-sseServer.ClosedClients:
			delete(sseServer.ClientConnMap, client)
			close(client)
		case srvMsg := <-sseServer.SrvMsgChan:
			for clientMsgChan := range sseServer.ClientConnMap {
				clientMsgChan <- srvMsg
			}
		}
	}
}

func (sseServer *SSEServer) StopListening() {
	sseServer.stopChan <- struct{}{}
}

func (sseServer *SSEServer) SrvHTTP() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Content-Type", "text/event-stream")
		ctx.Writer.Header().Set("Cache-Control", "no-cache")
		ctx.Writer.Header().Set("Connection", "keep-alive")
		ctx.Writer.Header().Set("Transfer-Encoding", "chunked")

		clientChan := make(StringMsgChan)
		sseServer.NewComeClients <- clientChan
		defer func() {
			sseServer.ClosedClients <- clientChan
		}()
		ctx.Set("clientChan", clientChan)
		ctx.Next()
	}
}
