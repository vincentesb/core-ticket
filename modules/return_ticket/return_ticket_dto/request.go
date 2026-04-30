package return_ticket_dto

type Request struct {
	ProductTypeID int    `json:"productTypeID" binding:"required,number,gte=1"`
	IssueType     string `json:"issueType" binding:"omitempty,ascii"`
}
