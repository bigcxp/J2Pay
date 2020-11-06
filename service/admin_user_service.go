package service

import (
	"golang.org/x/crypto/bcrypt"
	"j2pay-server/hcommon"
	"j2pay-server/model"
	"j2pay-server/model/request"
	"j2pay-server/model/response"
	"j2pay-server/myerr"
	"j2pay-server/pkg/casbin"
	"j2pay-server/validate"
	"time"
)

// 组织列表
func UserList(page, pageSize int) (res response.AdminUserPage, err error) {
	adminUser := model.AdminUser{}
	res, err = adminUser.GetAll(page, pageSize)
	if err != nil {
		return
	}
	for i, v := range res.Data {
		//查询组织的主账户
		if v.UserName != "" {
			user, err1 := model.GetAccountByWhere("user_name = ?", v.UserName)
			if err1 != nil {
				return
			}
			res.Data[i].Account.Token = user.Token
		}
	}
	return
}

// 组织详情
func UserDetail(id int64) (res response.AdminUserList, err error) {
	adminUser := model.AdminUser{ID: id}
	res, err = adminUser.Detail(id)
	if err != nil {
		return
	}
	if res.UserName != "" {
		user, err1 := model.GetAccountByWhere("user_name = ?", res.UserName)
		if err1 != nil {
			return
		}
		res.Account.Token = user.Token
	}
	return
}

// 创建组织
func UserAdd(user request.UserAdd) error {
	defer casbin.ClearEnforcer()
	time := time.Now().Unix()
	//组织
	u := model.AdminUser{
		UserName:      user.UserName,
		RealName:      user.RealName,
		Address:       user.Address,
		Balance:       user.Balance,
		ReturnUrl:     user.ReturnUrl,
		DaiUrl:        user.DaiUrl,
		Remark:        user.Remark,
		IsCollection:  user.IsCollection,
		IsCreation:    user.IsCreation,
		More:          user.More,
		OrderType:     user.OrderType,
		OrderCharge:   user.OrderCharge,
		ReturnType:    user.ReturnType,
		ReturnCharge:  user.ReturnCharge,
		IsDai:         user.IsDai,
		DaiType:       user.DaiType,
		DaiCharge:     user.DaiCharge,
		PickType:      user.PickType,
		PickCharge:    user.PickCharge,
		IsGas:         user.IsGas,
		Examine:       user.Examine,
		DayTotalCount: user.DayTotalCount,
		MaxOrderCount: user.MaxOrderCount,
		MinOrderCount: user.MinOrderCount,
		Limit:         user.Limit,
		UserLessTime:  user.UserLessTime,
		CreateTime:    time,
		UpdateTime:    time,
		WhitelistIP:   user.WhitelistIP,
	}
	// 1.判断组织是否存在
	if hasName, _ := model.GetUserByWhere("real_name = ?", user.RealName); hasName.ID > 0 {
		return myerr.NewDbValidateError("组织已存在")
	}
	// 2.判断密码是否一致
	if user.Password != user.RePassword {
		return myerr.NewNormalValidateError("密码与确认密码不一致")
	}
	// 3.判断账户是否存在
	if hasName, _ := model.GetAccountByWhere("user_name = ?", user.UserName); hasName.ID > 0 {
		return myerr.NewDbValidateError("账户名已存在")
	}
	Uid, err := u.Create()
	if err != nil {
		return err
	}
	//创建账户
	common := request.CommonAccount{UserName: user.UserName, UID: Uid, RID: 2}
	account := request.AccountAdd{
		Password:      user.Password,
		RePassword:    user.RePassword,
		CommonAccount: common,
	}
	// 密码加密处理
	bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	account.Password = string(bcryptPassword)
	a := model.Account{
		UID: account.UID,
		RID: 0,
		UserName:      account.UserName,
		Password:      account.Password,
		Secret:        validate.NewGoogleAuth().GetSecret(), //生成google唯一密钥
		QrcodeUrl:     "",
		Status:        hcommon.Open,
		IsOpen:        hcommon.No,
		IsMain:        hcommon.Yes,
		Token:         "",
		CreateTime:    time,
		UpdateTime:    time,
		LastLoginTime: time,
	}
	return a.Create(0)
}

// 编辑组织
func UserEdit(user request.UserEdit) error {
	defer casbin.ClearEnforcer()
	time := time.Now().Unix()
	u := model.AdminUser{
		ID:            user.ID,
		RealName:      user.RealName,
		Address:       user.Address,
		ReturnUrl:     user.ReturnUrl,
		DaiUrl:        user.DaiUrl,
		Remark:        user.Remark,
		IsCollection:  user.IsCollection,
		IsCreation:    user.IsCreation,
		More:          user.More,
		OrderType:     user.OrderType,
		OrderCharge:   user.OrderCharge,
		ReturnType:    user.ReturnType,
		ReturnCharge:  user.ReturnCharge,
		IsDai:         user.IsDai,
		DaiType:       user.DaiType,
		DaiCharge:     user.DaiCharge,
		IsGas:         user.IsGas,
		Examine:       user.Examine,
		DayTotalCount: user.DayTotalCount,
		MaxOrderCount: user.MaxOrderCount,
		MinOrderCount: user.MinOrderCount,
		Limit:         user.Limit,
		UserLessTime:  user.UserLessTime,
		UpdateTime:    time,
		CreateTime:    time,
	}
	// 1.判断组织名和手机号是否存在
	//if hasName, _ := model.GetUserByWhere("real_name = ? and id <> ?", user.RealName, user.ID); hasName.ID > 0 {
	//	return myerr.NewDbValidateError("组织名已存在")
	//}
	return u.Edit()
}
