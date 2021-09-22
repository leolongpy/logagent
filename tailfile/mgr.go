package tailfile

import (
	"github.com/sirupsen/logrus"
	"logagent/common"
)

type tailTaskMgr struct {
	tailTaskMap      map[string]*tailTask       //所有的任务
	collectEntryList []common.CollectEntry      //所有的配置项
	confChan         chan []common.CollectEntry // 等待新配置的通道
}

var ttMgr *tailTaskMgr

func Init(allconf []common.CollectEntry) (err error) {
	ttMgr = &tailTaskMgr{
		tailTaskMap:      make(map[string]*tailTask, 20),
		collectEntryList: allconf,
		confChan:         make(chan []common.CollectEntry),
	}
	for _, conf := range allconf {
		tt := newTailTask(conf.Path, conf.Topic)
		err = tt.Init()
		if err != nil {
			logrus.Errorf("create tailObj for path:%s failed, err:%v", conf.Path, err)
			continue
		}
		logrus.Infof("create a tail task for path:%s success", conf.Path)
		ttMgr.tailTaskMap[tt.path] = tt
		go tt.run()
	}
	go ttMgr.watch()
	return
}
func (t *tailTaskMgr) watch() {
	for {
		newConf := <-t.confChan
		logrus.Infof("get new conf from etcd, conf:%v, start manage tailTask...", newConf)
		for _, conf := range newConf {
			if t.isExits(conf) {
				continue
			}
			tt := newTailTask(conf.Path, conf.Topic)
			err := tt.Init()
			if err != nil {
				logrus.Errorf("create tailObj for path:%s failed, err:%v", conf.Path, err)
				continue
			}
			logrus.Infof("create a tail task for path:%s success", conf.Path)
			t.tailTaskMap[tt.path] = tt
			go tt.run()
			//删除不存在的
			for key, task := range t.tailTaskMap {
				var found bool
				for _, conf := range newConf {
					if key == conf.Path {
						found = true
						break
					}
				}
				if !found {
					// 这个tailTask要停掉了
					logrus.Infof("the task collect path:%s need to stop.", task.path)
					delete(t.tailTaskMap, key) // 从管理类中删掉
					task.cancel()
				}
			}
		}
	}
}
func (t *tailTaskMgr) isExits(conf common.CollectEntry) bool {
	_, ok := t.tailTaskMap[conf.Path]
	return ok
}

func SendNewConf(newConf []common.CollectEntry) {
	ttMgr.confChan <- newConf
}
