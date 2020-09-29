package request

type RateEdit struct {
	Id                       int     `json:"-"`
	ReceiveWeightType        int     `json:"receive_weight_type" binding:"oneof=0 1" example:"0"`          //代收加权类型：0：百分比，1：固定
	PayWeightType            int     `json:"pay_weight_type" binding:"oneof=0 1" example:"0"`              //代发加权类型：0：百分比，1：固定
	ReceiveWeightValue       float64 `json:"receive_weight_value"  binding:"required" example:"1"`         //代收加权值
	PayWeightValue           float64 `json:"pay_weight_value"  binding:"required" example:"1"`             //代发加权值
	PayWeightAddOrReduce     int     `json:"pay_weight_add_or_reduce" binding:"oneof=0 1" example:"0"`     //代发增加还是减少 0：增加 1：减少
	ReceiveWeightAddOrReduce int     `json:"receive_weight_add_or_reduce" binding:"oneof=0 1" example:"0"` //代收增加还是减少 0：增加 1：减少

}
