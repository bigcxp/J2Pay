package app

import (
	"context"
	"j2pay-server/app/model"
	"j2pay-server/hcommon"
	"j2pay-server/xenv"
	"time"
)

// GetLock 获取运行锁
func GetLock(ctx context.Context, tx hcommon.DbExeAble, k string) (bool, error) {
	genLock := func() error {
		_, err := SQLCreateTAppLockUpdate(
			ctx,
			tx,
			&model.DBTAppLock{
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

	lockRow, err := SQLGetTAppLockColByK(
		ctx,
		tx,
		[]string{
			model.DBColTAppLockCreateTime,
		},
		k,
	)
	if err != nil {
		return false, err
	}
	if lockRow == nil {
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
func ReleaseLock(ctx context.Context, tx hcommon.DbExeAble, k string) error {
	_, err := SQLUpdateTAppLockByK(
		ctx,
		tx,
		&model.DBTAppLock{
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
		context.Background(),
		xenv.DbCon,
		name,
	)
	if err != nil {
		hcommon.Log.Warnf("GetLock err: [%T] %s", err, err.Error())
		return
	}
	if !ok {
		return
	}
	defer func() {
		err := ReleaseLock(
			context.Background(),
			xenv.DbCon,
			name,
		)
		if err != nil {
			hcommon.Log.Warnf("ReleaseLock err: [%T] %s", err, err.Error())
			return
		}
	}()
	f()
}

// SQLGetWithdrawMap 获取提币map
func SQLGetWithdrawMap(ctx context.Context, tx hcommon.DbExeAble, cols []string, ids []int64) (map[int64]*model.DBTWithdraw, error) {
	if !hcommon.IsStringInSlice(cols, model.DBColTWithdrawID) {
		cols = append(cols, model.DBColTWithdrawID)
	}
	itemMap := make(map[int64]*model.DBTWithdraw)
	itemRows, err := model.SQLSelectTWithdrawCol(
		ctx,
		tx,
		cols,
		ids,
	)
	if err != nil {
		return nil, err
	}
	for _, itemRow := range itemRows {
		itemMap[itemRow.ID] = itemRow
	}
	return itemMap, nil
}

// SQLGetProductMap 获取产品map
func SQLGetProductMap(ctx context.Context, tx hcommon.DbExeAble, cols []string, ids []int64) (map[int64]*model.DBTProduct, error) {
	if !hcommon.IsStringInSlice(cols, model.DBColTProductID) {
		cols = append(cols, model.DBColTProductID)
	}
	itemMap := make(map[int64]*model.DBTProduct)
	itemRows, err := model.SQLSelectTProductCol(
		ctx,
		tx,
		cols,
		ids,
	)
	if err != nil {
		return nil, err
	}
	for _, itemRow := range itemRows {
		itemMap[itemRow.ID] = itemRow
	}
	return itemMap, nil
}

// SQLGetAppConfigTokenMap 获取代币map
func SQLGetAppConfigTokenMap(ctx context.Context, tx hcommon.DbExeAble, cols []string, ids []int64) (map[int64]*model.DBTAppConfigToken, error) {
	if !hcommon.IsStringInSlice(cols, model.DBColTAppConfigTokenID) {
		cols = append(cols, model.DBColTAppConfigTokenID)
	}
	itemMap := make(map[int64]*model.DBTAppConfigToken)
	itemRows, err := model.SQLSelectTAppConfigTokenCol(
		ctx,
		tx,
		cols,
		ids,
	)
	if err != nil {
		return nil, err
	}
	for _, itemRow := range itemRows {
		itemMap[itemRow.ID] = itemRow
	}
	return itemMap, nil
}

// SQLGetAddressKeyMap 获取地址map
func SQLGetAddressKeyMap(ctx context.Context, tx hcommon.DbExeAble, cols []string, addresses []string) (map[string]*model.DBTAddressKey, error) {
	if !hcommon.IsStringInSlice(cols, model.DBColTAddressKeyAddress) {
		cols = append(cols, model.DBColTAddressKeyAddress)
	}
	itemMap := make(map[string]*model.DBTAddressKey)
	itemRows, err := SQLSelectTAddressKeyColByAddress(
		ctx,
		tx,
		cols,
		addresses,
	)
	if err != nil {
		return nil, err
	}
	for _, itemRow := range itemRows {
		itemMap[itemRow.Address] = itemRow
	}
	return itemMap, nil
}
