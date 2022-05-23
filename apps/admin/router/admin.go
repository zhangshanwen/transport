package router

import (
	"github.com/gin-gonic/gin"

	admin "github.com/zhangshanwen/transport/apps/admin/controlers/v1"
	"github.com/zhangshanwen/transport/apps/admin/middleware"
	"github.com/zhangshanwen/transport/common"
)

func InitAdmin(Router *gin.RouterGroup) {
	r := Router.Group(common.Admins)
	{
		r.POST(common.UriLogin, middleware.H(admin.Login))        // 登录
		r.GET(common.UriEmpty, middleware.J(admin.Get))           // 获取所有管理员
		r.POST(common.UriEmpty, middleware.J(admin.Create))       // 创建管理员
		r.PUT(common.UriId, middleware.J(admin.Edit))             // 修改管理员信息
		r.DELETE(common.UriId, middleware.J(admin.Delete))        // 删除管理员
		r.PUT(common.UriAvatar, middleware.J(admin.UploadAvatar)) // 上传头像

		role := r.Group(common.Roles)
		{
			change := role.Group(common.Change)
			{
				change.PUT(common.UriId, middleware.J(admin.ChangeRole)) // 修改角色
			}
		}

		password := r.Group(common.Password)

		{
			password.PUT(common.Change, middleware.J(admin.ChangePassword)) // 修改密码
			reset := password.Group(common.Reset)
			{
				reset.GET(common.UriId, middleware.J(admin.ResetPassword)) // 重置密码
			}
		}
	}
}
