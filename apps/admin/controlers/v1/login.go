package v1

import (
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"github.com/zhangshanwen/transport/apps/admin/tools"
	"github.com/zhangshanwen/transport/code"
	"github.com/zhangshanwen/transport/common"
	"github.com/zhangshanwen/transport/initialize/db"
	"github.com/zhangshanwen/transport/initialize/service"
	"github.com/zhangshanwen/transport/internal/param"
	"github.com/zhangshanwen/transport/internal/response"
	"github.com/zhangshanwen/transport/model"
)

func Login(c *service.AdminContext) (resp service.Res) {
	p := param.AdminLogin{}
	if resp.Err = c.Rebind(&p); resp.Err != nil {
		resp.ResCode = code.ParamsError
		return
	}
	admin := model.Admin{Username: p.Username}
	g := db.G
	g = g.Begin()
	defer func() {
		if resp.Err == nil {
			g.Commit()
		} else {
			g.Rollback()
		}
	}()
	if resp.Err = g.Where(&admin).First(&admin).Error; resp.Err != nil {
		return
	}

	if !admin.CheckPassword(p.Password) {
		resp.ResCode = code.ActPWdError
		resp.Err = errors.New("ActPWdError")
		return
	}
	if resp.Err = g.Model(&admin).Updates(&model.Admin{
		LastLoginTime: time.Now().Unix(),
	}).Error; resp.Err != nil {
		return
	}
	r := response.AdminRolePermissionResponse{}
	if resp.Err = copier.Copy(&r, &admin); resp.Err != nil {
		return
	}
	var token string
	token, resp.Err = tools.CreateToken(admin.Id)
	if resp.Err != nil {
		return
	}
	r.Authorization = token
	if resp.Err = rolePermission(g, c, &admin, &r); resp.Err != nil {
		return
	}
	// TODO 向oss 模块获取地址
	//oss := tools.NewOss()
	r.Avatar.Name = admin.Avatar
	//r.Avatar.Url = oss.GetUrl(admin.Avatar)
	resp.Data = r
	return
}

func rolePermission(g *gorm.DB, c *service.AdminContext, admin *model.Admin, r *response.AdminRolePermissionResponse) (err error) {
	m := model.Role{}
	if err = g.Preload("Permissions").Preload("Permissions.Routes").First(&m, admin.RoleId).Error; err != nil {
		return
	}
	routes := map[string]interface{}{}
	permissionMap := make(map[int64]*model.Permission)
	for _, item := range m.Permissions {
		for _, route := range item.Routes {
			routes[fmt.Sprintf("%s||%s", route.Method, route.Path)] = true
		}
		permissionMap[item.Id] = item
		if item.Key != "" {
			r.Keys = append(r.Keys, item.Key)
		}
	}

	if err = db.R.HSet(c.Context, fmt.Sprintf(common.RedisRoutesKey, admin.Id), routes).Err(); err != nil {
		return
	}
	var parents []*model.Permission
	for k, v := range permissionMap {
		if v.ParentId == 0 {
			parents = append(parents, v)
		} else {
			parent := permissionMap[v.ParentId]
			if parent.ParentId == 0 {
				parent.Children = append(permissionMap[v.ParentId].Children, *permissionMap[k])
			}
		}
	}
	if err = copier.Copy(&r.List, &parents); err != nil {
		return
	}
	return
}
