package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/zhangshanwen/transport/apps/admin/tools"
	"github.com/zhangshanwen/transport/code"
	"github.com/zhangshanwen/transport/common"
	"github.com/zhangshanwen/transport/initialize/db"
	"github.com/zhangshanwen/transport/initialize/service"
	"github.com/zhangshanwen/transport/internal/header"
)

func verifyJwt(c *gin.Context) (res service.Res, claims *tools.Claims) {
	h := header.Authorization{}
	if res.Err = c.ShouldBindHeader(&h); res.Err != nil {
		res.StatusCode = http.StatusUnauthorized
		res.ResCode = code.AuthFailed
		return
	}
	claims, res.Err = tools.VerifyToken(h.Authorization)
	if res.Err != nil {
		return
	}
	return
}

func AdminJwtHandel(fun func(ctx *service.AdminContext) service.Res) gin.HandlerFunc {
	return func(c *gin.Context) {
		res, claims := verifyJwt(c)
		if res.Err != nil {
			res.StatusCode = http.StatusUnauthorized
			res.ResCode = code.AuthFailed
			service.Json(c, res)
			return
		}
		sC := &service.AdminContext{Context: c}
		if res.Err = db.G.First(&sC.Admin, claims.Payload.Uid).Error; res.Err != nil {
			res.StatusCode = http.StatusUnauthorized
			res.ResCode = code.AuthFailed
			service.Json(c, res)
			return
		}
		//if res.Err = verifyPermission(sC); res.Err != nil {
		//	res.StatusCode = http.StatusForbidden
		//	res.ResCode = code.NoPermission
		//	service.Json(c, res)
		//	return
		//}

		service.Json(c, fun(sC))
	}
}

func verifyPermission(c *service.AdminContext) (err error) {
	// verify last one  is number
	path := c.Request.URL.Path
	urlSplit := strings.Split(c.Request.URL.Path, common.Backlash)
	_, err = strconv.Atoi(urlSplit[len(urlSplit)-1])
	if err == nil {
		urlSplit[len(urlSplit)-1] = common.UriId
		path = strings.Join(urlSplit, common.Backlash)
	}
	key := fmt.Sprintf(common.RedisRoutesKey, c.Admin.Id)
	field := fmt.Sprintf("%s||%s", c.Request.Method, path)
	var val bool
	val, err = db.R.HGet(c, key, field).Bool()
	if err != nil {
		logrus.Error(err)
		return errors.New("No_Permission")
	}
	if !val {
		return errors.New("No_Permission")
	}
	return
}
