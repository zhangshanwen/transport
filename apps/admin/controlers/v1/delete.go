package v1

import (
	"github.com/zhangshanwen/transport/code"
	"github.com/zhangshanwen/transport/initialize/db"
	"github.com/zhangshanwen/transport/initialize/service"
	"github.com/zhangshanwen/transport/internal/param"
	"github.com/zhangshanwen/transport/model"
)

func Delete(c *service.AdminContext) (resp service.Res) {
	pId := param.UriId{}
	if resp.Err = c.BindUri(&pId); resp.Err != nil {
		resp.ResCode = code.ParamsError
		return
	}
	m := model.Admin{}
	m.Id = pId.Id
	g := db.G
	if resp.Err = g.Model(&m).Delete(&m).Error; resp.Err != nil {
		return
	}
	return
}
