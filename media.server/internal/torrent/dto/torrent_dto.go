package dto

type AddTorrentRequest struct {
	Magnet string `json:"magnet" binding:"required"`

	DownloadPath string `json:"download_path" binding:"required"`
}
