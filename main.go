package main

import (
	"fmt"
	"github.com/go-ini/ini"
	"github.com/sirupsen/logrus"
	"logagent/common"
	"logagent/etcd"
	"logagent/kafka"
	"logagent/tailfile"
)

type Config struct {
	KafkaConfig `ini:"kafka"`

	EtcdConfig `ini:"etcd"`
}
type KafkaConfig struct {
	Address  string `ini:"address"`
	ChanSize int64  `ini:"chan_size"`
}

type EtcdConfig struct {
	Address    string `ini:"address"`
	CollectKey string `ini:"collect_key"`
}

func main() {
	ip, err := common.GetOutboundIP()
	if err != nil {
		logrus.Errorf("get ip failed, err:%v", err)
		return
	}

	var configObj = new(Config)
	err = ini.MapTo(configObj, "./conf/config.ini")

	err = kafka.Init([]string{configObj.KafkaConfig.Address}, configObj.KafkaConfig.ChanSize)
	if err != nil {
		logrus.Errorf("init kafka failed, err:%v", err)
		return
	}
	logrus.Info("init kafka success!")

	err = etcd.Init([]string{configObj.EtcdConfig.Address})
	if err != nil {
		logrus.Errorf("init etcd failed, err:%v", err)
		return
	}
	collectKey := fmt.Sprintf(configObj.EtcdConfig.CollectKey, ip)
	logrus.Info(collectKey)
	allConf, err := etcd.GetConf(collectKey)
	if err != nil {
		logrus.Errorf("get conf from etcd failed, err:%v", err)
		return
	}
	go etcd.WatchConf(collectKey)
	err = tailfile.Init(allConf)
	if err != nil {
		logrus.Errorf("init tailfile failed, err:%v", err)
		return
	}
	logrus.Info("init tailfile success!")
	select {}

}
