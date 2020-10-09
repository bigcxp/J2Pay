package request

type AddressEdit struct {
	Id           int       `json:"-"`
	UserId    int     `json:"user_id" binding:"required" example:"1"`       //组织id

}
