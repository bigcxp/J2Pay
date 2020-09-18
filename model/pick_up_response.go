package model

type PickUpPage struct {
	Total       int    `json:"total"`        // 总共多少页
	PerPage     int    `json:"per_page"`     // 当前页码
	CurrentPage int    `json:"current_page"` // 每页显示多少条
	Data        []AdminUser `json:"data"`
}
