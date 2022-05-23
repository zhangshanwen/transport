package service

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/zhangshanwen/transport/model"
)

type AdminContext struct {
	*gin.Context
	Admin model.Admin
}

func (c *AdminContext) Rebind(obj interface{}) (err error) {
	if err = c.Bind(obj); err != nil {
		return
	}
	logrus.WithField("mod", "params").Info(obj)
	return
}
