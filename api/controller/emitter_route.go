package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/maogejing/emitter/bootstrap"
)

func NewEmitterRouter(sseServer *bootstrap.SSEServer, group *gin.RouterGroup) {
	group.GET("/listen-system-status", sseServer.SrvHTTP())
}
