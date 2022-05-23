package v1

import (
	"github.com/zhangshanwen/transport/code"
	"github.com/zhangshanwen/transport/initialize/db"
	"github.com/zhangshanwen/transport/initialize/service"
	"github.com/zhangshanwen/transport/internal/param"
	"github.com/zhangshanwen/transport/model"
)

func UploadAvatar(c *service.AdminContext) (resp service.Res) {
	p := param.AdminUploadAvatar{}
	if resp.Err = c.Rebind(&p); resp.Err != nil {
		resp.ResCode = code.ParamsError
		return
	}
	g := db.G
	if resp.Err = g.Model(&c.Admin).Updates(&model.Admin{
		Avatar: p.Avatar,
	}).Error; resp.Err != nil {
		return
	}
	resp.Data = c.Admin
	return
}
