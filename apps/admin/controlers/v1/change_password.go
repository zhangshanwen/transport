package v1

import (
	"github.com/zhangshanwen/transport/code"
	"github.com/zhangshanwen/transport/initialize/db"
	"github.com/zhangshanwen/transport/initialize/service"
	"github.com/zhangshanwen/transport/internal/param"
	"github.com/zhangshanwen/transport/internal/response"
	"github.com/zhangshanwen/transport/model"
)

func ChangePassword(c *service.AdminContext) (resp service.Res) {
	p := param.PasswordParam{}
	if resp.Err = c.Rebind(&p); resp.Err != nil {
		resp.ResCode = code.ParamsError
		return
	}
	g := db.G
	if resp.Err = c.Admin.SetPassword(p.Password); resp.Err != nil {
		return
	}
	if resp.Err = g.Model(&c.Admin).Updates(&model.Admin{
		Password: c.Admin.Password,
	}).Error; resp.Err != nil {
		return
	}
	resp.Data = response.PasswordResponse{Password: p.Password}
	return
}
