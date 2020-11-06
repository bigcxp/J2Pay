package service

import (
	"golang.org/x/crypto/bcrypt"
	"j2pay-server/hcommon"
	"j2pay-server/model"
	"j2pay-server/model/request"
	"j2pay-server/model/response"
	"j2pay-server/myerr"
	"j2pay-server/pkg/casbin"
	"j2pay-server/pkg/util"
	"j2pay-server/validate"
	"time"
)

// 登录逻辑
func Login(user *request.LoginUser) (string, error) {
	//获取用户
	account, _ := model.GetAccountByWhere("user_name = ?", user.Username)
	//验证google Code
	if account.IsOpen != 0 {
		code, err := validate.NewGoogleAuth().VerifyCode(account.Secret, user.GoogleCode)
		if err != nil {
			return "", err
		}
		if !code {
			return "", myerr.NewNormalValidateError("验证码错误")
		}
	}
	if account.ID == 0 {
		return "", myerr.NewNormalValidateError("用户不存在")
	}
	if bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(user.Password)) != nil {
		return "", myerr.NewNormalValidateError("用户密码错误")
	}
	if account.Status != 1 {
		return "", myerr.NewNormalValidateError("用户状态错误")
	}
	return util.MakeToken(account)
}

// 账户列表
func AccountList(userName string, page, pageSize int) (res response.AccountPage, err error) {
	account := model.Account{}
	if userName == "" {
		res, err = account.GetAll(page, pageSize)
	} else {
		res, err = account.GetAll(page, pageSize, "user_name like ?", "%"+userName+"%")
	}
	if err != nil {
		return
	}
	for i, v := range res.Data {
		res.Data[i].Roles = response.CasRole{}
		res.Data[i].User = response.User{}
		//查询用户对应角色
		if v.RID != 0 {
			role, err1 := model.GetRoleByWhere("id = ?", v.RID)
			if err1 != nil {
				return
			}
			res.Data[i].Roles.ID = role.ID
			res.Data[i].Roles.Name = role.Name
		}
		//查询用户对应组织
		if v.UID != 0 {
			user, err1 := model.GetUserByWhere("id = ?", v.UID)
			if err1 != nil {
				return
			}
			res.Data[i].User.ID = user.ID
			res.Data[i].User.RealName = user.RealName
		}
	}
	return
}

// 用户详情
func AccountDetail(id int64) (res response.AccountList, err error) {
	account := model.Account{ID: id}
	res, err = account.AccountDetail(id)
	if err != nil {
		return
	}
	//查询用户对应角色
	if res.RID != 0 {
		role, err1 := model.GetRoleByWhere("id = ?", res.RID)
		if err1 != nil {
			return
		}
		res.Roles.ID = role.ID
		res.Roles.Name = role.Name
	}
	//查询用户对应组织
	if res.UID != 0 {
		user, err1 := model.GetUserByWhere("id = ?", res.UID)
		if err1 != nil {
			return
		}
		res.User.ID = user.ID
		res.User.RealName = user.RealName
	}
	return
}

// 创建用户
func AccountAdd(user request.AccountAdd) error {
	defer casbin.ClearEnforcer()
	u := model.Account{
		UID:           user.UID,
		RID:           user.RID,
		UserName:      user.UserName,
		Password:      user.Password,
		Secret:        validate.NewGoogleAuth().GetSecret(), //生成google唯一密钥
		QrcodeUrl:     "",
		Status:        hcommon.Open,
		IsOpen:        hcommon.No,
		IsMain:        hcommon.No,
		Token:         "",
		CreateTime:    time.Now().Unix(),
		UpdateTime:    time.Now().Unix(),
		LastLoginTime: time.Now().Unix(),
	}
	// 1.判断账户名是否存在
	if hasName, _ := model.GetAccountByWhere("user_name = ?", user.UserName); hasName.ID > 0 {
		return myerr.NewDbValidateError("用户名已存在")
	}
	// 2.判断密码是否一致
	if user.Password != user.RePassword {
		return myerr.NewNormalValidateError("密码与确认密码不一致")
	}
	// 3.判断角色是否存在
	hasRoles, err := model.GetRoleByWhere("id = ?)", user.RID)
	if err != nil {
		return err
	}
	if hasRoles.ID != user.RID {
		return myerr.NewDbValidateError("选择的角色不存在")
	}
	// 4.密码加密处理
	bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(bcryptPassword)
	return u.Create(user.RID)
}

// 编辑账户
func AccountEdit(account request.AccountEdit) error {
	defer casbin.ClearEnforcer()
	u := model.Account{
		ID:         account.ID,
		Status:     account.Status,
		UpdateTime: time.Now().Unix(),
	}
	return u.Edit(account.RID)
}

//修改密码
func UpdatePassword(id int64) (password response.Password, err error) {
	defer casbin.ClearEnforcer()
	account := model.Account{ID: id}
	//随机获取密码
	password1 := util.RandString(10)
	userPassword := response.Password{Password: password1}
	//加密
	bcryptPassword, _ := bcrypt.GenerateFromPassword([]byte(password1), bcrypt.DefaultCost)
	Password2 := string(bcryptPassword)
	err = account.UpdatePassword(id, Password2)
	return userPassword, err
}

//开启google验证
func OpenGoogle(google request.Google) (err error) {
	defer casbin.ClearEnforcer()
	user, _ := model.GetAccountByWhere("id = ?", google.ID)
	code, err2 := validate.NewGoogleAuth().VerifyCode(user.Secret, google.GoogleCode)
	if err2 != nil {
		return err2
	}
	if !code {
		return myerr.NewDbValidateError("动态验证码错误")
	}
	err = user.Google(google)
	return err
}

// 删除用户
func UserDel(id int64) error {
	defer casbin.ClearEnforcer()
	u := model.Account{
		ID: id,
	}
	return u.Del()
}

// 数据库保存Token 及更新最后登录时间
func EditToken(username, token string) error {
	defer casbin.ClearEnforcer()
	u := model.Account{
		UserName: username,
		Token:    token,
	}
	return u.EditToken(username, token)
}
