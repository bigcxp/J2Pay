package request

type RateEdit struct {
	ID                       int     `json:"id" form:"id"`
	ReceiveWeightType        int     `json:"receive_weight_type" binding:"oneof=0 1" example:"0" form:"receive_weight_type"`                   //代收加权类型：0：百分比，1：固定
	PayWeightType            int     `json:"pay_weight_type" binding:"oneof=0 1" example:"0" form:"pay_weight_type"`                           //代发加权类型：0：百分比，1：固定
	ReceiveWeightValue       float64 `json:"receive_weight_value"  binding:"required" example:"1" form:"float64"`                              //代收加权值
	PayWeightValue           float64 `json:"pay_weight_value"  binding:"required" example:"1" form:"pay_weight_value"`                         //代发加权值
	PayWeightAddOrReduce     int     `json:"pay_weight_add_or_reduce" binding:"oneof=0 1" example:"0" form:"pay_weight_add_or_reduce"`         //代发增加还是减少 0：增加 1：减少
	ReceiveWeightAddOrReduce int     `json:"receive_weight_add_or_reduce" binding:"oneof=0 1" example:"0" form:"receive_weight_add_or_reduce"` //代收增加还是减少 0：增加 1：减少

}
