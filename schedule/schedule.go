package schedule

import (
	"fmt"
	"time"

	"github.com/zfd81/magpie/server"

	"github.com/zfd81/magpie/config"

	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

var (
	conf              = config.GetConfig()
	LogExpirationTime = 0 - conf.Log.ExpirationTime
	ClearTaskTime     = conf.Log.ClearTaskTime
	crontab           = cron.New() // 新建一个定时任务对象
)

func ClearTask() {
	date := time.Now().AddDate(0, 0, LogExpirationTime).Format("20060102")
	err := server.Remove(date)
	if err != nil {
		log.Error(err)
	}
}

func StartScheduler() {
	spec := fmt.Sprintf("0 0 %d * * ?", ClearTaskTime) //定时任务
	crontab.AddFunc(spec, ClearTask)                   // 添加定时任务
	crontab.Start()                                    // 启动定时器
}
