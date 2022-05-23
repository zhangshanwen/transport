package router

import (
	"github.com/zhangshanwen/transport/common"
	"github.com/zhangshanwen/transport/initialize/app"
)

func RegisterRouter() {
	api := app.R.Group(common.BackendPrefix)
	group := api.Group(common.V1)

	InitAdmin(group)
}
