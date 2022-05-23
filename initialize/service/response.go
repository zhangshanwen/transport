package service

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/zhangshanwen/transport/code"
	"github.com/zhangshanwen/transport/initialize/node"
)

type Res struct {
	StatusCode int         `json:"-"`
	ResCode    int         `json:"code"`
	ReqId      int64       `json:"req_id"`
	Data       interface{} `json:"data"`
	Err        error       `json:"err,omitempty"`
}

func success(c *gin.Context, r Res) {
	if r.ResCode == 0 {
		r.ResCode = code.BaseSuccessCode
	}
	c.JSON(http.StatusOK, r)
}

func failed(c *gin.Context, r Res) {
	if r.ResCode == 0 {
		r.ResCode = code.BaseFailedCode
	}
	c.JSON(http.StatusBadRequest, r)
}

func Json(c *gin.Context, r Res) {
	r.ReqId = node.N.Generate()
	if r.StatusCode == 0 {
		if r.Err == nil {
			success(c, r)
		} else {
			failed(c, r)
		}
		return
	}
	c.JSON(r.StatusCode, r)
}
