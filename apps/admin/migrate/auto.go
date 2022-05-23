package migrate

import (
	"github.com/sirupsen/logrus"
	"github.com/zhangshanwen/transport/initialize/db"
	"github.com/zhangshanwen/transport/model"
)

func AutoMigrate() {
	logrus.Info("--------mysql_auto_migrate_start---------")
	_ = db.G.AutoMigrate(
		&model.Admin{},
		&model.Route{},
		&model.Permission{},
		&model.Role{},
	)
	logrus.Info("--------mysql_auto_migrate_end---------")
}
