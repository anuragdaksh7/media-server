package dto

type StreamQuery struct {
	Path string `form:"path" binding:"required"`
}

type DownloadQuery struct {
	Path string `form:"path" binding:"required"`
}
