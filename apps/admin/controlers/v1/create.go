package v1

import (
	"errors"

	"github.com/jinzhu/copier"

	"github.com/zhangshanwen/transport/code"
	"github.com/zhangshanwen/transport/initialize/db"
	"github.com/zhangshanwen/transport/initialize/service"
	"github.com/zhangshanwen/transport/internal/param"
	"github.com/zhangshanwen/transport/internal/response"
	"github.com/zhangshanwen/transport/model"
)

func Create(c *service.AdminContext) (resp service.Res) {
	p := param.Register{}
	if resp.Err = c.Rebind(&p); resp.Err != nil {
		resp.ResCode = code.ParamsError
		return
	}
	m := model.Admin{Username: p.Username}
	g := db.G
	var count int64
	if resp.Err = g.Model(&m).Where(&m).Count(&count).Error; resp.Err != nil {
		return
	}
	if count > 0 {
		resp.Err = errors.New("username is existed")
		resp.ResCode = code.UsernameIsExisted
		return
	}
	if resp.Err = copier.Copy(&m, &p); resp.Err != nil {
		return
	}
	if resp.Err = m.SetPassword(p.Password); resp.Err != nil {
		return
	}
	if resp.Err = g.Create(&m).Error; resp.Err != nil {
		return
	}
	r := response.AdminInfo{}
	if resp.Err = copier.Copy(&r, &m); resp.Err != nil {
		return
	}
	resp.Data = m
	return
}
