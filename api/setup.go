package api

import (
	"github.com/gin-gonic/gin"
	"github.com/maogejing/emitter/api/controller"
)

func Setup(gin *gin.Engine) {
	emitterRouter := gin.Group("emitter")
	controller.NewEmitterRouter(emitterRouter)
}
