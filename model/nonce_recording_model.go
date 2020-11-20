package model

type NonceRecording struct {
	ID         int64  `gorm:"default:0;comment:'ID'"`
	Nonce      int64  `gorm:"default:0;comment:'交易索引'"`
	OthersUsed int    `gorm:"default:0;comment:'是否被第三方使用'"`
	ErrMsg     string `gorm:"default:0;comment:'错误信息'"`
	FormAddr   string `gorm:"default:0;comment:'交易发起地址'"`
}

type MaxAndMinNonce struct {
	MaxNonce int64 `gorm:"max_nonce"`
	MinNonce int64 `gorm:"min_nonce"`
}

//通过交易发起地址查询nonce最大值和最小值
func (*NonceRecording) SQLGetTSendMaxAndMinNonceByFrom(address string, rpcNonce int64) (*MaxAndMinNonce, error) {
	row := MaxAndMinNonce{}
	err := GetDb().Raw("select IFNULL(MAX(nonce), -2) max_nonce,IFNULL(min(nonce), -2) min_nonce  from nonce_recording where form_addr = ? and nonce> ?", address, rpcNonce).Scan(&row).Error
	return &row, err
}

// 保存nonce信息
func (a *NonceRecording) Create() error {
	if err := GetDb().Table("nonce_recording").Create(a).Error; err != nil {
		return err
	}
	return nil
}
