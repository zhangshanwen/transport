package v1

import (
	"fmt"
	"github.com/jinzhu/copier"

	"github.com/zhangshanwen/transport/code"
	"github.com/zhangshanwen/transport/initialize/db"
	"github.com/zhangshanwen/transport/initialize/service"
	"github.com/zhangshanwen/transport/internal/param"
	"github.com/zhangshanwen/transport/internal/response"
	"github.com/zhangshanwen/transport/model"
)

func Get(c *service.AdminContext) (resp service.Res) {
	p := param.AdminRecords{}
	if resp.Err = c.Rebind(&p); resp.Err != nil {
		resp.ResCode = code.ParamsError
		return
	}
	m := model.Admin{}
	var ms []model.Admin
	g := db.G.Model(&m)
	if p.Username != "" {
		m.Username = fmt.Sprintf("%%%s%%", p.Username)
		g = g.Where(&m)
	}
	r := response.AdminResponse{}
	if resp.Err = db.FindByPagination(g, &p.Pagination, &r.Pagination); resp.Err != nil {
		return
	}
	if resp.Err = g.Preload("Role").Find(&ms).Error; resp.Err != nil {
		return
	}
	if resp.Err = copier.Copy(&r.List, &ms); resp.Err != nil {
		return
	}
	resp.Data = r
	return
}
