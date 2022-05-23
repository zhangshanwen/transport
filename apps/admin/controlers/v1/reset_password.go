package v1

import (
	"github.com/zhangshanwen/transport/apps/admin/conf"
	"github.com/zhangshanwen/transport/code"
	"github.com/zhangshanwen/transport/initialize/db"
	"github.com/zhangshanwen/transport/initialize/service"
	"github.com/zhangshanwen/transport/internal/param"
	"github.com/zhangshanwen/transport/internal/response"
	"github.com/zhangshanwen/transport/model"
)

func ResetPassword(c *service.AdminContext) (resp service.Res) {
	pId := param.UriId{}
	if resp.Err = c.BindUri(&pId); resp.Err != nil {
		resp.ResCode = code.ParamsError
		return
	}
	admin := model.Admin{}
	g := db.G
	if resp.Err = g.First(&admin, pId.Id).Error; resp.Err != nil {
		return
	}
	if resp.Err = admin.SetPassword(conf.C.ResetPassword); resp.Err != nil {
		return
	}
	if resp.Err = g.Model(&admin).Updates(&model.Admin{
		Password: admin.Password,
	}).Error; resp.Err != nil {
		return
	}
	resp.Data = response.PasswordResponse{Password: conf.C.ResetPassword}
	return
}
