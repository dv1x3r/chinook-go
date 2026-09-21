package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/dv1x3r/chinook-go/internal/lib/wrap"

	"github.com/dv1x3r/w2go/w2"
	"github.com/dv1x3r/w2go/w2db"
	"github.com/huandu/go-sqlbuilder"
)

var (
	ErrPlaylistNotFound   = errors.New("playlist not found")
	ErrPlaylistHasRelated = errors.New("playlist still has tracks")
)

type Playlist struct {
	ID     int              `json:"id"`
	Name   w2.Field[string] `json:"name"`
	Tracks int              `json:"tracks"`
}

func (s *Service) GetPlaylistGrid(ctx context.Context, req w2.GetGridRequest) (w2.GetGridResponse[Playlist], error) {
	const op = "api.Service.GetPlaylistGrid"
	res, err := w2db.GetGridContext(ctx, s.store.DB(), req, w2db.GetGridOptions[Playlist]{
		From: "playlists as t",
		Select: []string{
			"t.PlaylistId as id",
			"t.Name as name",
			"coalesce(r.count, 0) as tracks",
		},
		Where: map[string]string{
			"name": "t.Name",
		},
		OrderBy: map[string]string{
			"id":     "t.PlaylistId",
			"name":   "t.Name",
			"tracks": "r.count",
		},
		Build: func(sb *sqlbuilder.SelectBuilder) {
			sb.JoinWithOption(sqlbuilder.LeftJoin,
				"(select PlaylistId, count(*) as count from playlist_track group by PlaylistId) as r",
				"r.PlaylistId = t.PlaylistId",
			)
		},
		Scan: func(rows *sql.Rows, record *Playlist) error {
			return rows.Scan(
				&record.ID,
				&record.Name,
				&record.Tracks,
			)
		},
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) CreatePlaylist(ctx context.Context, req w2.SaveFormRequest[Playlist]) (int, error) {
	const op = "api.Service.CreatePlaylist"
	id, err := w2db.InsertContext(ctx, s.store.DB(), w2db.InsertOptions{
		Into: "playlists",
		Values: map[string]any{
			"Name": req.Record.Name,
		},
	})
	return id, wrap.IfErr(op, err)
}

func (s *Service) UpdatePlaylist(ctx context.Context, req w2.SaveFormRequest[Playlist]) error {
	const op = "api.Service.UpdatePlaylist"
	affected, err := w2db.UpdateContext(ctx, s.store.DB(), w2db.UpdateOptions{
		Update: "playlists",
		Values: map[string]any{
			"Name": req.Record.Name,
		},
		Where: map[string]any{
			"PlaylistId": req.Record.ID,
		},
	})
	if err != nil {
		return wrap.IfErr(op, err)
	} else if affected == 0 {
		return wrap.IfErr(op, ErrPlaylistNotFound)
	}
	return nil
}

func (s *Service) DeletePlaylists(ctx context.Context, req w2.RemoveGridRequest) error {
	const op = "api.Service.DeletePlaylists"
	affected, err := w2db.RemoveGridContext(ctx, s.store.DB(), req, w2db.RemoveGridOptions{
		From:    "playlists",
		IDField: "PlaylistId",
	})
	if s.store.IsErrConstraintForeignKey(err) {
		return wrap.IfErr(op, ErrPlaylistHasRelated)
	} else if err != nil {
		return wrap.IfErr(op, err)
	} else if affected != len(req.ID) {
		return wrap.IfErr(op, ErrPlaylistNotFound)
	}
	return nil
}
