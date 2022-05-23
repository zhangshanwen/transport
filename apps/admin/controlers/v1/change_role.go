package v1

import (
	"github.com/zhangshanwen/transport/code"
	"github.com/zhangshanwen/transport/initialize/db"
	"github.com/zhangshanwen/transport/initialize/service"
	"github.com/zhangshanwen/transport/internal/param"
	"github.com/zhangshanwen/transport/model"
)

func ChangeRole(c *service.AdminContext) (resp service.Res) {
	pId := param.UriId{}
	if resp.Err = c.BindUri(&pId); resp.Err != nil {
		resp.ResCode = code.ParamsError
		return
	}
	p := param.AdminChangeRole{}
	if resp.Err = c.Rebind(&p); resp.Err != nil {
		resp.ResCode = code.ParamsError
		return
	}
	g := db.G
	m := model.Admin{}
	if pId.Id == c.Admin.Id {
		m = c.Admin
	} else {
		if resp.Err = g.First(&m, pId).Error; resp.Err != nil {
			return
		}
	}
	if resp.Err = g.Model(&m).Updates(&model.Admin{
		RoleId: p.RoleId,
	}).Error; resp.Err != nil {
		return
	}
	resp.Data = m.Role
	return
}
