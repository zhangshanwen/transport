package v1

import (
	"github.com/jinzhu/copier"

	"github.com/zhangshanwen/transport/code"
	"github.com/zhangshanwen/transport/initialize/db"
	"github.com/zhangshanwen/transport/initialize/service"
	"github.com/zhangshanwen/transport/internal/param"
	"github.com/zhangshanwen/transport/model"
)

func Edit(c *service.AdminContext) (resp service.Res) {
	pId := param.UriId{}
	if resp.Err = c.BindUri(&pId); resp.Err != nil {
		resp.ResCode = code.ParamsError
		return
	}
	p := param.AdminEdit{}
	if resp.Err = c.Rebind(&p); resp.Err != nil {
		resp.ResCode = code.ParamsError
		return
	}
	m := model.Admin{}
	g := db.G
	if resp.Err = g.First(&m, pId.Id).Error; resp.Err != nil {
		return
	}
	if resp.Err = copier.Copy(&m, &p); resp.Err != nil {
		return
	}
	if resp.Err = g.Model(&m).Updates(&m).Error; resp.Err != nil {
		return
	}
	resp.Data = m
	return
}
