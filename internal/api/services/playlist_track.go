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
	ErrPlaylistTrackExists   = errors.New("track is already on this playlist")
	ErrPlaylistTrackNotFound = errors.New("track is not on this playlist")
)

type PlaylistTrack struct {
	ID           int         `json:"id"`
	Name         string      `json:"name"`
	Artist       w2.Dropdown `json:"artist"`
	Album        w2.Dropdown `json:"album"`
	Genre        w2.Dropdown `json:"genre"`
	Milliseconds int         `json:"milliseconds"`
}

type PlaylistTrackForm struct {
	Track w2.Dropdown `json:"track"`
}

func (s *Service) GetPlaylistTrackGrid(ctx context.Context, req w2.GetGridRequest, playlistID int) (w2.GetGridResponse[PlaylistTrack], error) {
	const op = "api.Service.GetPlaylistTrackGrid"
	res, err := w2db.GetGridContext(ctx, s.store.DB(), req, w2db.GetGridOptions[PlaylistTrack]{
		From: "playlist_track as pt",
		Select: []string{
			"t.TrackId as id",
			"t.Name as name",
			"ar.ArtistId as artist_id",
			"ar.Name as artist_text",
			"t.AlbumId as album_id",
			"al.Title as album_text",
			"t.GenreId as genre_id",
			"g.Name as genre_text",
		},
		Where: map[string]string{
			"id":     "t.TrackId",
			"name":   "t.Name",
			"artist": "ar.Name",
			"album":  "al.Title",
		},
		OrderBy: map[string]string{
			"id":     "t.TrackId",
			"name":   "t.Name",
			"artist": "ar.Name",
			"album":  "al.Title",
			"genre":  "g.Name",
		},
		Build: func(sb *sqlbuilder.SelectBuilder) {
			sb.JoinWithOption(sqlbuilder.InnerJoin, "tracks as t", "t.TrackId = pt.TrackId")
			sb.JoinWithOption(sqlbuilder.LeftJoin, "albums as al", "al.AlbumId = t.AlbumId")
			sb.JoinWithOption(sqlbuilder.LeftJoin, "artists as ar", "ar.ArtistId = al.ArtistId")
			sb.JoinWithOption(sqlbuilder.LeftJoin, "genres as g", "g.GenreId = t.GenreId")
			sb.Where(sb.EQ("pt.PlaylistId", playlistID))
			if len(req.Sort) == 0 {
				sb.OrderBy("t.Name")
			}
		},
		Scan: func(rows *sql.Rows, record *PlaylistTrack) error {
			return rows.Scan(
				&record.ID,
				&record.Name,
				&record.Artist.ID,
				&record.Artist.Text,
				&record.Album.ID,
				&record.Album.Text,
				&record.Genre.ID,
				&record.Genre.Text,
			)
		},
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) AddPlaylistTrack(ctx context.Context, req w2.SaveFormRequest[PlaylistTrackForm], playlistID int) error {
	const op = "api.Service.AddPlaylistTrack"
	_, err := w2db.InsertContext(ctx, s.store.DB(), w2db.InsertOptions{
		Into: "playlist_track",
		Values: map[string]any{
			"PlaylistId": playlistID,
			"TrackId":    req.Record.Track.ID,
		},
	})
	if s.store.IsErrConstraintUnique(err) {
		return wrap.IfErr(op, ErrPlaylistTrackExists)
	}
	return wrap.IfErr(op, err)
}

func (s *Service) RemovePlaylistTracks(ctx context.Context, req w2.RemoveGridRequest, playlistID int) error {
	const op = "api.Service.RemovePlaylistTracks"

	if len(req.ID) == 0 {
		return nil
	}

	db := sqlbuilder.DeleteFrom("playlist_track")
	db.Where(
		db.EQ("PlaylistId", playlistID),
		db.In("TrackId", sqlbuilder.Flatten(req.ID)...),
	)

	query, args := db.BuildWithFlavor(sqlbuilder.SQLite)
	res, err := s.store.DB().ExecContext(ctx, query, args...)
	if err != nil {
		return wrap.IfErr(op, err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return wrap.IfErr(op, err)
	} else if int(affected) != len(req.ID) {
		return wrap.IfErr(op, ErrPlaylistTrackNotFound)
	}

	return nil
}
