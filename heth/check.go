package heth

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/parnurzeal/gorequest"
	"j2pay-server/hcommon"
	"j2pay-server/model"
	"j2pay-server/pkg/util"
	"net/http"
	"time"
)

// CheckDoNotify 检测发送回调
func CheckDoNotify() {
	// 初始化的
	initNotifyRows, err := model.SQLSelectTProductNotifyColByStatusAndTime(
		hcommon.NotifyStatusInit,
		time.Now().Unix(),
	)
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return
	}
	// 失败的
	delayNotifyRows, err := model.SQLSelectTProductNotifyColByStatusAndTime(
		hcommon.NotifyStatusFail,
		time.Now().Add(-time.Minute*10).Unix(),
	)
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return
	}
	initNotifyRows = append(initNotifyRows, delayNotifyRows...)
	for _, initNotifyRow := range initNotifyRows {
		//如果地址不合法，不发送http请求
		if util.IsHttp(initNotifyRow.URL) {
			gresp, body, errs := gorequest.New().Post(initNotifyRow.URL).Timeout(time.Second * 30).Send(initNotifyRow.Msg).End()
			if errs != nil {
				hcommon.Log.Errorf("err: [%T] %s", errs[0], errs[0].Error())
				err = model.SQLUpdateTProductNotifyStatusByID(
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
				err = model.SQLUpdateTProductNotifyStatusByID(
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
				err = model.SQLUpdateTProductNotifyStatusByID(
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
				err = model.SQLUpdateTProductNotifyStatusByID(
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
				err = model.SQLUpdateTProductNotifyStatusByID(
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
}
