package service

import (
	"context"
	model2 "fileserver/internal/auth/model"
	"fileserver/internal/config"
	"fileserver/internal/realtime/handler"
	"fileserver/internal/torrent/dto"
	"fileserver/internal/torrent/manager"
	"fileserver/internal/torrent/model"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/google/uuid"
)

var conf config.Config

func init() {
	var err error
	conf, err = config.LoadConfig(".")
	if err != nil {
		panic(err)
	}
}

type TorrentService struct {
	Client *torrent.Client

	Manager *manager.TorrentManager

	Realtime *handler.RealtimeHandler
}

func NewTorrentService(
	manager *manager.TorrentManager,
	realtime *handler.RealtimeHandler,
) (*TorrentService, error) {

	client, err := torrent.NewClient(&torrent.ClientConfig{
		DataDir: conf.StoragePath,
	})

	if err != nil {
		return nil, err
	}

	return &TorrentService{
		Client: client,

		Manager: manager,

		Realtime: realtime,
	}, nil
}

func (s *TorrentService) AddTorrent(
	ctx context.Context,
	req dto.AddTorrentRequest,
) (*model.Torrent, error) {

	torrentModel := &model.Torrent{
		Base: model2.Base{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},

		Magnet: req.Magnet,

		DownloadPath: req.DownloadPath,

		Status: model.TorrentPending,
	}

	s.Manager.Add(torrentModel)

	go s.download(torrentModel)

	return torrentModel, nil
}

func (s *TorrentService) download(
	torrentModel *model.Torrent,
) {

	torrentModel.Status = model.TorrentDownloading

	t, err := s.Client.AddMagnet(
		torrentModel.Magnet,
	)

	if err != nil {

		torrentModel.Status = model.TorrentFailed

		torrentModel.Error = err.Error()

		return
	}

	<-t.GotInfo()

	torrentModel.Name = t.Name()

	torrentModel.Size = t.Length()

	t.DownloadAll()

	torrentModel.Status = model.TorrentDownloading

	ticker := time.NewTicker(
		2 * time.Second,
	)

	defer ticker.Stop()

	for range ticker.C {

		completed := t.BytesCompleted()

		total := t.Length()

		progress := float64(completed) /
			float64(total) * 100

		torrentModel.Progress = progress

		torrentModel.Peers = len(
			t.PeerConns(),
		)

		s.Realtime.BroadcastJSON(
			"torrent_progress",
			torrentModel,
		)

		if completed >= total {

			torrentModel.Progress = 100

			torrentModel.Status = model.TorrentCompleted

			s.Realtime.BroadcastJSON(
				"torrent_completed",
				torrentModel,
			)

			return
		}
	}
}
