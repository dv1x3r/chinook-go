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
	ErrTrackNotFound   = errors.New("track not found")
	ErrTrackHasRelated = errors.New("track is part of an invoice or a playlist")
)

type Track struct {
	ID           int               `json:"id"`
	Name         w2.Field[string]  `json:"name"`
	Album        w2.Dropdown       `json:"album"`
	Genre        w2.Dropdown       `json:"genre"`
	MediaType    w2.Dropdown       `json:"media_type"`
	Composer     w2.Field[string]  `json:"composer"`
	Milliseconds w2.Field[int]     `json:"milliseconds"`
	Bytes        w2.Field[int]     `json:"bytes"`
	UnitPrice    w2.Field[float64] `json:"unit_price"`
}

func trackValues(record Track) map[string]any {
	return map[string]any{
		"Name":         record.Name,
		"AlbumId":      record.Album.ID,
		"GenreId":      record.Genre.ID,
		"MediaTypeId":  record.MediaType.ID,
		"Composer":     record.Composer,
		"Milliseconds": record.Milliseconds,
		"Bytes":        record.Bytes,
		"UnitPrice":    record.UnitPrice,
	}
}

func (s *Service) GetTrackGrid(ctx context.Context, req w2.GetGridRequest, artistID int, albumID int) (w2.GetGridResponse[Track], error) {
	const op = "api.Service.GetTrackGrid"
	res, err := w2db.GetGridContext(ctx, s.store.DB(), req, w2db.GetGridOptions[Track]{
		From: "tracks as t",
		Select: []string{
			"t.TrackId as id",
			"t.Name as name",
			"t.AlbumId as album_id",
			"al.Title as album_text",
			"t.GenreId as genre_id",
			"g.Name as genre_text",
			"t.MediaTypeId as media_type_id",
			"mt.Name as media_type_text",
			"t.Composer as composer",
			"t.Milliseconds as milliseconds",
			"t.Bytes as bytes",
			"t.UnitPrice as unit_price",
		},
		Where: map[string]string{
			"id":         "t.TrackId",
			"name":       "t.Name",
			"composer":   "t.Composer",
			"album":      "t.AlbumId",
			"genre":      "t.GenreId",
			"media_type": "t.MediaTypeId",
		},
		OrderBy: map[string]string{
			"id":           "t.TrackId",
			"name":         "t.Name",
			"album":        "al.Title",
			"genre":        "g.Name",
			"media_type":   "mt.Name",
			"composer":     "t.Composer",
			"milliseconds": "t.Milliseconds",
			"bytes":        "t.Bytes",
			"unit_price":   "t.UnitPrice",
		},
		Build: func(sb *sqlbuilder.SelectBuilder) {
			sb.JoinWithOption(sqlbuilder.LeftJoin, "albums as al", "al.AlbumId = t.AlbumId")
			sb.JoinWithOption(sqlbuilder.LeftJoin, "genres as g", "g.GenreId = t.GenreId")
			sb.JoinWithOption(sqlbuilder.LeftJoin, "media_types as mt", "mt.MediaTypeId = t.MediaTypeId")
			if artistID > 0 {
				sb.Where(sb.EQ("al.ArtistId", artistID))
			}
			if albumID > 0 {
				sb.Where(sb.EQ("t.AlbumId", albumID))
			}
		},
		Scan: func(rows *sql.Rows, record *Track) error {
			return rows.Scan(
				&record.ID,
				&record.Name,
				&record.Album.ID,
				&record.Album.Text,
				&record.Genre.ID,
				&record.Genre.Text,
				&record.MediaType.ID,
				&record.MediaType.Text,
				&record.Composer,
				&record.Milliseconds,
				&record.Bytes,
				&record.UnitPrice,
			)
		},
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) GetTrackDropdown(ctx context.Context, req w2.GetDropdownRequest) (w2.GetDropdownResponse[w2.Dropdown], error) {
	const op = "api.Service.GetTrackDropdown"
	res, err := w2db.GetDropdownContext(ctx, s.store.DB(), req, w2db.GetDropdownOptions{
		From:         "tracks",
		IDField:      "TrackId",
		TextField:    "Name",
		OrderByField: "Name",
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) CreateTrack(ctx context.Context, req w2.SaveFormRequest[Track]) (int, error) {
	const op = "api.Service.CreateTrack"
	id, err := w2db.InsertContext(ctx, s.store.DB(), w2db.InsertOptions{
		Into:   "tracks",
		Values: trackValues(req.Record),
	})
	return id, wrap.IfErr(op, err)
}

func (s *Service) UpdateTrack(ctx context.Context, req w2.SaveFormRequest[Track]) error {
	const op = "api.Service.UpdateTrack"
	affected, err := w2db.UpdateContext(ctx, s.store.DB(), w2db.UpdateOptions{
		Update: "tracks",
		Values: trackValues(req.Record),
		Where: map[string]any{
			"TrackId": req.Record.ID,
		},
	})
	if err != nil {
		return wrap.IfErr(op, err)
	} else if affected == 0 {
		return wrap.IfErr(op, ErrTrackNotFound)
	}
	return nil
}

func (s *Service) SaveTrackGrid(ctx context.Context, req w2.SaveGridRequest[Track]) error {
	const op = "api.Service.SaveTrackGrid"
	affected, err := w2db.SaveGridContext(ctx, s.store.DB(), req, w2db.SaveGridOptions[Track]{
		BuildOptions: func(change Track) w2db.UpdateOptions {
			return w2db.UpdateOptions{
				Update: "tracks",
				Values: trackValues(change),
				Where: map[string]any{
					"TrackId": change.ID,
				},
			}
		},
	})
	if err != nil {
		return wrap.IfErr(op, err)
	} else if affected != len(req.Changes) {
		return wrap.IfErr(op, ErrTrackNotFound)
	}
	return nil
}

func (s *Service) DeleteTracks(ctx context.Context, req w2.RemoveGridRequest) error {
	const op = "api.Service.DeleteTracks"
	affected, err := w2db.RemoveGridContext(ctx, s.store.DB(), req, w2db.RemoveGridOptions{
		From:    "tracks",
		IDField: "TrackId",
	})
	if s.store.IsErrConstraintForeignKey(err) {
		return wrap.IfErr(op, ErrTrackHasRelated)
	} else if err != nil {
		return wrap.IfErr(op, err)
	} else if affected != len(req.ID) {
		return wrap.IfErr(op, ErrTrackNotFound)
	}
	return nil
}
