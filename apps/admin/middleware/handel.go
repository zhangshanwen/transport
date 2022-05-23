package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/zhangshanwen/transport/initialize/service"
)

var (
	H = AdminHandel
	J = AdminJwtHandel
)

func AdminHandel(fun func(c *service.AdminContext) service.Res) gin.HandlerFunc {
	return func(c *gin.Context) {
		service.Json(c, fun(&service.AdminContext{Context: c}))
	}
}
