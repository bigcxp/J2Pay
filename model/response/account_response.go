package response

//账户返回列表
type AccountPage struct {
	Total       int           `json:"total"`        // 总共多少页
	PerPage     int           `json:"per_page"`     // 当前页码
	CurrentPage int           `json:"current_page"` // 每页显示多少条
	Data        []AccountList `json:"data"`
}

//角色
type CasRole struct {
	ID   int    `json:"id"`   // 角色ID
	Name string `json:"name"` // 角色名
}

//组织
type User struct {
	ID       int64  `json:"id"`        // 组织ID
	RealName string `json:"real_name"` // 组织名称
}

//账户列表
type AccountList struct {
	ID            int64   `json:"id"`
	UID           int64   `json:"uid"`             //组织ID
	RID           int     `json:"rid"`             // 角色ID
	QrcodeUrl     string  `json:"qr_code_url"`     //google 二维码地址
	Secret        string  `json:"secret"`          //google密钥
	UserName      string  `json:"user_name"`       // 登录名
	Token         string  `json:"token"`           //token
	Status        int     `json:"status"`          // 用户状态
	CreateTime    int64   `json:"create_time"`     //创建时间
	UpdateTime    int64   `json:"update_time"`     //更新时间
	LastLoginTime int64   `json:"last_login_time"` //最后登录时间
	User          User    `json:"user"`            //商户信息
	Roles         CasRole `json:"roles"`           // 角色信息
	IsOpen        int     `json:"is_open"`         //是否开启google双重验证 默认1：开启 ，2：关闭 ';"
}

//返回修改的密码
type Password struct {
	Password string `json:"password"` //密码
}
