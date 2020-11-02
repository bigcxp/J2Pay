package response

//商户实收订单列表
type OrderPage struct {
	Total          int             `json:"total"`           // 总共多少页
	PerPage        int             `json:"per_page"`        // 当前页码
	CurrentPage    int             `json:"current_page"`    // 每页显示多少条
	TotalAmount    float64         `json:"total_amount"`    //总订单金额
	MerchantAmount float64         `json:"merchant_amount"` //总商户总实收金额
	ReallyAmount   float64         `json:"really_amount"`   //总实收金额
	TotalFee       float64         `json:"total_fee"`       //总手续费
	Data           []RealOrderList `json:"data"`
}

//管理端商户实收订单对象
type RealOrderList struct {
	ID               uint    `json:"id"`                 //ID
	UserId           int     `json:"user_id"`            //商户id
	RealName         string  `json:"real_name"`          //组织名称
	OrderCode        string  `json:"order_code"`         //商户订单编号
	Amount           float64 `json:"amount"`             //金额
	TXID             string  `json:"txid"`               //交易hash
	DetailedRecordId float64 `json:"detailed_record_id"` //实收明细订单编号
	ShouldAmount     float64 `json:"should_amount"`      //应收金额
	ReceiptAmount    float64 `json:"receipt_amount"`     //实收金额
	Fee              float64 `json:"fee"`                //手续费
	ReturnAmount     float64 `json:"return_amount"`      //退款金额
	MerchantAmount   float64 `json:"merchant_amount"`    //商户实收金额
	Status           int     `json:"status"`             //状态 -1：收款中，1：已完成，2：异常，3：退款等待中，4：退款中，5：退款失败，6：已退款，7：：已过期
	FinishTime       int64   `json:"finish_time"`        //完成时间
	CreateAt         int64   `json:"create_at"`          //创建时间
	Remark           string  `json:"remark"`             //备注
	Address          string  `json:"address"`            //收款地址

}

//返回给用户充币地址
type UserAddr struct {
	OrderCode      string  `json:"order_code"`      //商户订单编号
	Amount         float64 `json:"amount"`          //金额
	Address        string  `json:"charge_address"`  //分配的收款地址
	ExprireTime    int64   `json:"exprireTime"`     //过期时间
	Currency       string  `json:"currency"`        //币别
	CurrencyAmount float64 `json:"currency_amount"` //币别金额
}
