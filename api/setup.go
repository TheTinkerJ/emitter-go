package api

import (
	"github.com/gin-gonic/gin"
	"github.com/maogejing/emitter/api/controller"
	"github.com/maogejing/emitter/bootstrap"
)

func Setup(sseServer *bootstrap.SSEServer, gin *gin.Engine) {
	emitterRouter := gin.Group("emitter")
	controller.NewEmitterRouter(sseServer, emitterRouter)
}
