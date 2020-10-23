package service

import (
	"golang.org/x/crypto/bcrypt"
	"j2pay-server/model"
	"j2pay-server/model/request"
	"j2pay-server/model/response"
	"j2pay-server/myerr"
	"j2pay-server/pkg/casbin"
	"j2pay-server/pkg/logger"
	"j2pay-server/pkg/util"
	"j2pay-server/validate"
	"time"
)

// 登录逻辑
func Login(user *request.LoginUser, id string) (string, error) {
	//bcryptPassword, _:= bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	//Password := string(bcryptPassword)
	//fmt.Println(Password)
	//if ! base64Captcha.DefaultMemStore.Verify(id, user.VerifyCode, true) {
	//	return "", myerr.NewNormalValidateError("验证码错误")
	//}
	adminUser := model.GetUserByWhere("user_name = ?", user.Username)
	//验证google Code
	if adminUser.IsOpen != 0 {
		code, err := validate.NewGoogleAuth().VerifyCode(adminUser.Secret, user.GoogleCode)
		if err != nil {
			return "", err
		}
		if !code {
			return "", myerr.NewNormalValidateError("验证码错误")
		}

	}
	if adminUser.ID== 0 {
		return "", myerr.NewNormalValidateError("用户不存在")
	}
	if bcrypt.CompareHashAndPassword([]byte(adminUser.Password), []byte(user.Password)) != nil {
		return "", myerr.NewNormalValidateError("用户密码错误")
	}
	if adminUser.Status != 1 {
		return "", myerr.NewNormalValidateError("用户状态错误")
	}
	return util.MakeToken(adminUser)
}

// 用户列表
func UserList(pid int, page, pageSize int) (res response.AdminUserPage, err error) {
	adminUser := model.AdminUser{}
	if pid == 0 {
		res, err = adminUser.GetAll(page, pageSize, "pid = 0")
	} else {
		//res, err = adminUser.GetAll(page, pageSize, "user_name like ? or tel like ?", "%"+name+"%", "%"+name+"%")
		res, err = adminUser.GetAll(page, pageSize)
	}
	if err != nil {
		return
	}
	roles := model.GetAllRole()
	mappings := model.GetUserRoleMapping()
	for i, v := range res.Data {
		_, ok := mappings[v.Id]
		if !ok {
			continue
		}
		res.Data[i].Roles = []response.CasRole{}
		for _, role := range mappings[v.Id] {
			if _, ok := roles[role]; !ok {
				logger.Logger.Error("角色获取错误: user_id = ", v.Id)
				continue
			}
			res.Data[i].Roles = append(res.Data[i].Roles, roles[role])
		}
	}
	return
}

// 用户详情
func UserDetail(id int) (res response.AdminUserList, err error) {
	adminUser := model.AdminUser{ID: id}
	res, err = adminUser.Detail()
	if err != nil {
		return
	}
	res.Roles = model.GetUserRole(res.Id)
	return
}

// 创建用户
func UserAdd(user request.UserAdd) error {
	defer casbin.ClearEnforcer()
	u := model.AdminUser{
		UserName:      user.UserName,
		Pid:           user.Pid,
		Tel:           user.Tel,
		Password:      user.Password,
		RealName:      user.RealName,
		Secret:        validate.NewGoogleAuth().GetSecret(), //生成唯一密钥
		Status:        user.Status,
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
		CreateTime:    time.Now(),
		UpdateTime:    time.Now(),
		LastLoginTime: time.Now(),
		IsOpen:        0,
	}
	// 1.判断用户名和手机号是否存在
	if hasName := model.GetUserByWhere("user_name = ?", user.UserName); hasName.ID > 0 {
		return myerr.NewDbValidateError("用户名已存在")
	}
	if hasTel := model.GetUserByWhere("tel = ?", user.Tel); hasTel.ID > 0 {
		return myerr.NewDbValidateError("手机号已存在")
	}

	// 2.判断角色是否存在
	hasRoles, err := model.GetRolesByWhere("id in (?)", user.Roles)
	if err != nil {
		return err
	}
	if len(hasRoles) != len(user.Roles) {
		return myerr.NewDbValidateError("选择的角色不存在")
	}

	// 3.密码加密处理
	bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(bcryptPassword)
	return u.Create(user.Roles)
}

// 编辑用户
func UserEdit(user request.UserEdit) error {
	defer casbin.ClearEnforcer()
	u := model.AdminUser{
		ID:            user.ID,
		UserName:      user.UserName,
		IsOpen:        user.IsOpen,
		Tel:           user.Tel,
		Password:      user.Password,
		RealName:      user.RealName,
		Status:        user.Status,
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
		UpdateTime:    time.Now(),
		CreateTime:    time.Now(),
		LastLoginTime: time.Now(),
	}

	// 1.判断用户名和手机号是否存在
	if hasName := model.GetUserByWhere("user_name = ? and id <> ?", user.UserName, user.ID); hasName.ID > 0 {
		return myerr.NewDbValidateError("用户名已存在")
	}
	if hasTel := model.GetUserByWhere("tel = ? and id <> ?", user.UserName, user.ID); hasTel.ID > 0 {
		return myerr.NewDbValidateError("手机号已存在")
	}
	// 2.判断角色是否存在
	hasRoles, err := model.GetRolesByWhere("id in (?)", user.Roles)
	if err != nil {
		return err
	}
	if len(hasRoles) != len(user.Roles) {
		return myerr.NewDbValidateError("选择的角色不存在")
	}

	// 3.密码加密处理
	if u.Password != "" {
		bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(bcryptPassword)
	}

	return u.Edit(user.Roles)
}

//修改密码
func UpdatePassword(id int) (password response.Password, err error) {
	defer casbin.ClearEnforcer()
	adminUser := model.AdminUser{ID: id}
	//随机获取密码
	password1 := util.RandString(10)
	userPassword := response.Password{Password: password1}
	//加密
	bcryptPassword, _ := bcrypt.GenerateFromPassword([]byte(password1), bcrypt.DefaultCost)
	Password2 := string(bcryptPassword)
	err = adminUser.UpdatePassword(id, Password2)
	return userPassword, err
}

//开启google验证
func OpenGoogle(google request.Google) (err error) {
	defer casbin.ClearEnforcer()
	user := model.GetUserByWhere("id = ?", google.ID)
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
func UserDel(id int) error {
	defer casbin.ClearEnforcer()
	u := model.AdminUser{
		ID: id,
	}
	return u.Del()
}

// 更新Token
func EditToken(token string, username string) error {
	defer casbin.ClearEnforcer()
	u := model.AdminUser{
		UserName: username,
	}
	return u.EditToken(token, username)
}
