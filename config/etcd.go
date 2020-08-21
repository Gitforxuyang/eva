package config

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Gitforxuyang/eva/util/logger"
	"github.com/Gitforxuyang/eva/util/trace"
	"github.com/coreos/etcd/clientv3"
)

const (
	//服务注册前缀
	ETCD_CONFIG_PREFIX = "/eva/config/"
)

func watch(client *clientv3.Client) {
	go func() {
		for true {
		LOOP:
			w := client.Watch(context.TODO(), fmt.Sprintf("%s%s", ETCD_CONFIG_PREFIX, config.name))
			for wresp := range w {
				for _, ev := range wresp.Events {
					err := config.v.MergeConfig(bytes.NewBuffer(ev.Kv.Value))
					if err != nil {
						logger.GetLogger().Error(context.TODO(), "动态更新配置出错", logger.Fields{"err": err})
						goto LOOP
					}
					err = config.v.UnmarshalKey("trace", &config.trace)
					if err != nil {
						logger.GetLogger().Error(context.TODO(), "动态更新配置出错", logger.Fields{"err": err})
						goto LOOP
					}
					err = config.v.UnmarshalKey("log", &config.log)
					if err != nil {
						logger.GetLogger().Error(context.TODO(), "动态更新配置出错", logger.Fields{"err": err})
						goto LOOP
					}
					err = config.v.UnmarshalKey("dynamic", &config.dynamic)
					if err != nil {
						logger.GetLogger().Error(context.TODO(), "动态更新配置出错", logger.Fields{"err": err})
						goto LOOP
					}
					//动态修改采集比例
					trace.SetRatio(config.trace.Ratio)
					//主动通知动态配置发生变化
					config.changeNotify(config.dynamic)
					logger.GetLogger().Info(context.TODO(), "动态更新配置成功", logger.Fields{
						"dynamic": config.dynamic,
						"log":     config.log,
						"trace":   config.trace,
					})

				}
			}

		}
	}()
}
