package response

import "time"

//汇率表返回列表
type RatePage struct {
	Data        []Rate`json:"data"`
}


type Rate struct {
	ID                       uint    `json:"id"`                            //ID
	Currency                 string  `json:"currency"`                      //币别
	OriginalRate             float64 `json:"original_rate"`                 //原汇率
	Collection               float64 `json:"collection"`                    //代收加权
	Payment                  float64 `json:"payment"`                       //代发加权
	ReceiveWeightType        int     `json:"receive_weight_type"`           //代收加权类型：0：百分比，1：固定
	PayWeightType            int     `json:"pay_weight_type"`               //代发加权类型：0：百分比，1：固定
	ReceiveWeightValue       float64 `json:"receive_weight_value"`          //代收加权值
	PayWeightValue           float64 `json:"pay_weight_value"`              //代发加权值
	PayWeightAddOrReduce     int     `json:"pay_weight_add_or_reduce" `     //代发增加还是减少 0：增加 1：减少
	ReceiveWeightAddOrReduce int     `json:"receive_weight_add_or_reduce" ` //代收增加还是减少 0：增加 1：减少

	UpdatedAt time.Time `json:"updated_at"` //更新时间
}
