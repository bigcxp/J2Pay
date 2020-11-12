package heth

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/parnurzeal/gorequest"
	"j2pay-server/hcommon"
	"j2pay-server/model"
	"net/http"
	"time"
)

// CheckDoNotify 检测发送回调

func CheckDoNotify() {
	// 初始化的
	notify := model.TUserNotify{}
	initNotifyRows, err := notify.SQLSelectTProductNotifyColByStatusAndTime(
		hcommon.NotifyStatusInit,
		time.Now().Unix(),
	)
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return
	}
	// 错误的
	delayNotifyRows, err := notify.SQLSelectTProductNotifyColByStatusAndTime(
		hcommon.NotifyStatusFail,
		time.Now().Add(-time.Minute*10).Unix(),
	)
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return
	}
	initNotifyRows = append(initNotifyRows, delayNotifyRows...)
	for _, initNotifyRow := range initNotifyRows {
		gresp, body, errs := gorequest.New().Post(initNotifyRow.URL).Timeout(time.Second * 30).Send(initNotifyRow.Msg).End()
		if errs != nil {
			hcommon.Log.Errorf("err: [%T] %s", errs[0], errs[0].Error())
			err = notify.SQLUpdateTProductNotifyStatusByID(
				&model.TUserNotify{
					ID:           initNotifyRow.ID,
					HandleStatus: hcommon.NotifyStatusFail,
					HandleMsg:    errs[0].Error(),
					UpdateTime:   time.Now().Unix(),
				},
			)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			}
			continue
		}
		if gresp.StatusCode != http.StatusOK {
			// 状态错误
			hcommon.Log.Errorf("req status error: %d", gresp.StatusCode)
			err = notify.SQLUpdateTProductNotifyStatusByID(
				&model.TUserNotify{
					ID:           initNotifyRow.ID,
					HandleStatus: hcommon.NotifyStatusFail,
					HandleMsg:    fmt.Sprintf("http status: %d", gresp.StatusCode),
					UpdateTime:   time.Now().Unix(),
				},
			)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			}
			continue
		}
		resp := gin.H{}
		err = json.Unmarshal([]byte(body), &resp)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			err = notify.SQLUpdateTProductNotifyStatusByID(
				&model.TUserNotify{
					ID:           initNotifyRow.ID,
					HandleStatus: hcommon.NotifyStatusFail,
					HandleMsg:    body,
					UpdateTime:   time.Now().Unix(),
				},
			)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			}
			continue
		}
		_, ok := resp["error"]
		if ok {
			// 处理成功
			err = notify.SQLUpdateTProductNotifyStatusByID(
				&model.TUserNotify{
					ID:           initNotifyRow.ID,
					HandleStatus: hcommon.NotifyStatusPass,
					HandleMsg:    body,
					UpdateTime:   time.Now().Unix(),
				},
			)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			}
		} else {
			//hcommon.Log.Errorf("no error in resp")
			err = notify.SQLUpdateTProductNotifyStatusByID(
				&model.TUserNotify{
					ID:           initNotifyRow.ID,
					HandleStatus: hcommon.NotifyStatusFail,
					HandleMsg:    body,
					UpdateTime:   time.Now().Unix(),
				},
			)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			}
			continue
		}
	}
}
