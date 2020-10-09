package conf

import (
	"context"
	"github.com/Gitforxuyang/eva/config"
	"github.com/Gitforxuyang/eva/util/logger"
	"github.com/mitchellh/mapstructure"
)

var (
	dynamic *Dynamic
)

type Dynamic struct {
}

func Registry() {
	config.RegisterNotify(func(c map[string]interface{}) {
		d := &Dynamic{}
		err := mapstructure.Decode(c, &d)
		if err != nil {
			logger.GetLogger().Error(context.TODO(), "获取动态配置错误", logger.Fields{
				"err": err,
			})
		}
		dynamic = d
	})
}

func GetDynamic() *Dynamic {
	return dynamic
}

