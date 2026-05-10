package dto

type UploadQuery struct {
	Path string `form:"path" binding:"required"`
}
