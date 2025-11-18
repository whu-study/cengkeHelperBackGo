package dto

// ChatRequestDTO
type ChatRequestDTO struct {
	Query string `json:"query" binding:"required"`
}
