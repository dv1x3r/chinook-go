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
	ErrAlbumNotFound   = errors.New("album not found")
	ErrAlbumHasRelated = errors.New("album still has tracks")
)

type Album struct {
	ID     int              `json:"id"`
	Title  w2.Field[string] `json:"title"`
	Artist w2.Dropdown      `json:"artist"`
	Tracks int              `json:"tracks"`
}

func (s *Service) GetAlbumGrid(ctx context.Context, req w2.GetGridRequest, artistID int) (w2.GetGridResponse[Album], error) {
	const op = "api.Service.GetAlbumGrid"
	res, err := w2db.GetGridContext(ctx, s.store.DB(), req, w2db.GetGridOptions[Album]{
		From: "albums as al",
		Select: []string{
			"al.AlbumId as id",
			"al.Title as title",
			"al.ArtistId as artist_id",
			"ar.Name as artist_text",
			"coalesce(r.count, 0) as tracks",
		},
		Where: map[string]string{
			"title": "al.Title",
		},
		OrderBy: map[string]string{
			"id":     "al.AlbumId",
			"title":  "al.Title",
			"artist": "ar.Name",
			"tracks": "r.count",
		},
		Build: func(sb *sqlbuilder.SelectBuilder) {
			sb.JoinWithOption(sqlbuilder.InnerJoin, "artists as ar", "ar.ArtistId = al.ArtistId")
			sb.JoinWithOption(sqlbuilder.LeftJoin,
				"(select AlbumId, count(*) as count from tracks group by AlbumId) as r",
				"r.AlbumId = al.AlbumId",
			)
			if artistID != 0 {
				sb.Where(sb.EQ("al.ArtistId", artistID))
			}
		},
		Scan: func(rows *sql.Rows, record *Album) error {
			return rows.Scan(
				&record.ID,
				&record.Title,
				&record.Artist.ID,
				&record.Artist.Text,
				&record.Tracks,
			)
		},
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) GetAlbumDropdown(ctx context.Context, req w2.GetDropdownRequest) (w2.GetDropdownResponse[w2.Dropdown], error) {
	const op = "api.Service.GetAlbumDropdown"
	res, err := w2db.GetDropdownContext(ctx, s.store.DB(), req, w2db.GetDropdownOptions{
		From:         "albums",
		IDField:      "AlbumId",
		TextField:    "Title",
		OrderByField: "Title",
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) CreateAlbum(ctx context.Context, req w2.SaveFormRequest[Album]) (int, error) {
	const op = "api.Service.CreateAlbum"
	id, err := w2db.InsertContext(ctx, s.store.DB(), w2db.InsertOptions{
		Into: "albums",
		Values: map[string]any{
			"Title":    req.Record.Title,
			"ArtistId": req.Record.Artist.ID,
		},
	})
	return id, wrap.IfErr(op, err)
}

func (s *Service) UpdateAlbum(ctx context.Context, req w2.SaveFormRequest[Album]) error {
	const op = "api.Service.UpdateAlbum"
	affected, err := w2db.UpdateContext(ctx, s.store.DB(), w2db.UpdateOptions{
		Update: "albums",
		Values: map[string]any{
			"Title":    req.Record.Title,
			"ArtistId": req.Record.Artist.ID,
		},
		Where: map[string]any{
			"AlbumId": req.Record.ID,
		},
	})
	if err != nil {
		return wrap.IfErr(op, err)
	} else if affected == 0 {
		return wrap.IfErr(op, ErrAlbumNotFound)
	}
	return nil
}

func (s *Service) DeleteAlbums(ctx context.Context, req w2.RemoveGridRequest) error {
	const op = "api.Service.DeleteAlbums"
	affected, err := w2db.RemoveGridContext(ctx, s.store.DB(), req, w2db.RemoveGridOptions{
		From:    "albums",
		IDField: "AlbumId",
	})
	if s.store.IsErrConstraintForeignKey(err) {
		return wrap.IfErr(op, ErrAlbumHasRelated)
	} else if err != nil {
		return wrap.IfErr(op, err)
	} else if affected != len(req.ID) {
		return wrap.IfErr(op, ErrAlbumNotFound)
	}
	return nil
}
