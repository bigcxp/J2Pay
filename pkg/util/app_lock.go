package util

import (
	"j2pay-server/model"
	"log"
	"time"
)

// GetLock 获取运行锁

func GetLock(k string) (bool, error) {
	genLock := func() error {
		_, err := model.SQLCreateAppLockUpdate(
			&model.AppLock{
				K:          k,
				V:          1,
				CreateTime: time.Now().Unix(),
			},
		)
		if err != nil {
			return err
		}
		return nil
	}

	lockRow, err := model.SQLGetAppLockColByK(k)
	if err != nil {
		err = genLock()
		if err != nil {
			return false, err
		}
		return true, nil
	}
	if time.Now().Unix()-lockRow.CreateTime > 60*30 {
		err = genLock()
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

// ReleaseLock 释放运行锁
func ReleaseLock(k string) error {
	_, err := model.SQLUpdateTAppLockByK(
		&model.AppLock{
			K:          k,
			V:          0,
			CreateTime: time.Now().Unix(),
		},
	)
	if err != nil {
		return err
	}
	return nil
}

// LockWrap 包装被lock的函数
func LockWrap(name string, f func()) {
	ok, err := GetLock(
		name,
	)
	if err != nil {
		log.Panicf("GetLock err: [%T] %s", err, err.Error())
		return
	}
	if !ok {
		return
	}
	defer func() {
		err := ReleaseLock(
			name,
		)
		if err != nil {
			log.Panicf("ReleaseLock err: [%T] %s", err, err.Error())
			return
		}
	}()
	f()
}
