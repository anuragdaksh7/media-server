package dto

type ListQuery struct {
	Path string `form:"path"`
}

type StatQuery struct {
	Path string `form:"path" binding:"required"`
}

type DeleteRequest struct {
	Path string `json:"path" binding:"required"`
}

type MoveRequest struct {
	Source      string `json:"source" binding:"required"`
	Destination string `json:"destination" binding:"required"`
}

type CopyRequest struct {
	Source      string `json:"source" binding:"required"`
	Destination string `json:"destination" binding:"required"`
}

type MkdirRequest struct {
	Path string `json:"path" binding:"required"`
}
