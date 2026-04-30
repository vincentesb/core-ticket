package dto

type FilterDto struct {
	Page    int `json:"page,omitempty"    form:"page,default=1"     binding:"numeric,gte=1"`
	PerPage int `json:"perPage,omitempty" form:"perPage,default=20" binding:"numeric,gte=1,lte=50"`
}
