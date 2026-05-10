package model

import "fileserver/internal/auth/model"

type TorrentStatus string

const (
	TorrentPending     TorrentStatus = "pending"
	TorrentDownloading TorrentStatus = "downloading"
	TorrentCompleted   TorrentStatus = "completed"
	TorrentPaused      TorrentStatus = "paused"
	TorrentFailed      TorrentStatus = "failed"
)

type Torrent struct {
	model.Base
	Name          string        `json:"name"`
	Magnet        string        `json:"magnet"`
	DownloadPath  string        `json:"download_path"`
	Status        TorrentStatus `json:"status"`
	Progress      float64       `json:"progress"`
	DownloadSpeed int64         `json:"download_speed"`
	Size          int64         `json:"size"`
	Peers         int           `json:"peers"`
	Error         string        `json:"error,omitempty"`
}
