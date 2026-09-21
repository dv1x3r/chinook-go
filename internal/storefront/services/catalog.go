package services

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/dv1x3r/chinook-go/internal/lib/wrap"

	"github.com/dv1x3r/w2go/w2db"
	"github.com/huandu/go-sqlbuilder"
)

// TrackPageSize is how many rows one scroll of the catalog loads.
const TrackPageSize = 25

type Playlist struct {
	ID     int
	Name   string
	Tracks int
}

type Genre struct {
	ID   int
	Name string
}

type MediaType struct {
	ID   int
	Name string
}

type Track struct {
	ID       int
	Name     string
	Album    string
	Artist   string
	Genre    string
	Media    string
	Composer string
	Length   int
	Bytes    int
	Price    float64
}

// TrackRequest carries what the toolbar and the column dropdowns send: the
// search box with the field it searches in, the checked column filters as the
// JSON the table model keeps, and the offset of the infinite scroll.
type TrackRequest struct {
	Search     string
	SearchBy   string
	Filters    string
	PlaylistID int
	Offset     int
}

// GetPlaylists feeds the sidebar, the sample database has a few empty ones which are left out.
func (s *Service) GetPlaylists(ctx context.Context) ([]Playlist, error) {
	const op = "storefront.Service.GetPlaylists"
	records, err := w2db.SelectContext(ctx, s.store.DB(), w2db.SelectOptions[Playlist]{
		Build: func(sb *sqlbuilder.SelectBuilder) {
			sb.Select("p.PlaylistId", "p.Name", "count(*) as tracks")
			sb.From("playlists as p")
			sb.JoinWithOption(sqlbuilder.InnerJoin, "playlist_track as pt", "pt.PlaylistId = p.PlaylistId")
			sb.GroupBy("p.PlaylistId")
			sb.OrderBy("p.Name", "p.PlaylistId")
		},
		Scan: func(rows *sql.Rows, record *Playlist) error {
			return rows.Scan(&record.ID, &record.Name, &record.Tracks)
		},
	})
	return records, wrap.IfErr(op, err)
}

func (s *Service) CountTracks(ctx context.Context) (int, error) {
	const op = "storefront.Service.CountTracks"
	var count int
	err := s.store.DB().QueryRowContext(ctx, "select count(*) from tracks;").Scan(&count)
	return count, wrap.IfErr(op, err)
}

func (s *Service) GetGenres(ctx context.Context) ([]Genre, error) {
	const op = "storefront.Service.GetGenres"
	records, err := w2db.SelectContext(ctx, s.store.DB(), w2db.SelectOptions[Genre]{
		Build: func(sb *sqlbuilder.SelectBuilder) {
			sb.Select("GenreId", "Name")
			sb.From("genres")
			sb.OrderBy("Name")
		},
		Scan: func(rows *sql.Rows, record *Genre) error {
			return rows.Scan(&record.ID, &record.Name)
		},
	})
	return records, wrap.IfErr(op, err)
}

func (s *Service) GetMediaTypes(ctx context.Context) ([]MediaType, error) {
	const op = "storefront.Service.GetMediaTypes"
	records, err := w2db.SelectContext(ctx, s.store.DB(), w2db.SelectOptions[MediaType]{
		Build: func(sb *sqlbuilder.SelectBuilder) {
			sb.Select("MediaTypeId", "Name")
			sb.From("media_types")
			sb.OrderBy("Name")
		},
		Scan: func(rows *sql.Rows, record *MediaType) error {
			return rows.Scan(&record.ID, &record.Name)
		},
	})
	return records, wrap.IfErr(op, err)
}

func (s *Service) GetTracks(ctx context.Context, req TrackRequest) ([]Track, error) {
	const op = "storefront.Service.GetTracks"

	filters, err := parseFilters(req.Filters)
	if err != nil {
		return nil, wrap.IfErr(op, err)
	}

	records, err := w2db.SelectContext(ctx, s.store.DB(), w2db.SelectOptions[Track]{
		Build: func(sb *sqlbuilder.SelectBuilder) {
			sb.Select(
				"t.TrackId",
				"t.Name",
				"coalesce(al.Title, '')",
				"coalesce(ar.Name, '')",
				"coalesce(g.Name, '')",
				"coalesce(mt.Name, '')",
				"coalesce(t.Composer, '')",
				"t.Milliseconds",
				"coalesce(t.Bytes, 0)",
				"t.UnitPrice",
			)
			sb.From("tracks as t")
			sb.JoinWithOption(sqlbuilder.LeftJoin, "albums as al", "al.AlbumId = t.AlbumId")
			sb.JoinWithOption(sqlbuilder.LeftJoin, "artists as ar", "ar.ArtistId = al.ArtistId")
			sb.JoinWithOption(sqlbuilder.LeftJoin, "genres as g", "g.GenreId = t.GenreId")
			sb.JoinWithOption(sqlbuilder.LeftJoin, "media_types as mt", "mt.MediaTypeId = t.MediaTypeId")

			if req.PlaylistID > 0 {
				sb.JoinWithOption(sqlbuilder.InnerJoin, "playlist_track as pt", "pt.TrackId = t.TrackId")
				sb.Where(sb.EQ("pt.PlaylistId", req.PlaylistID))
			}

			if req.Search != "" {
				pattern := "%" + req.Search + "%"
				switch req.SearchBy {
				case "Album":
					sb.Where(sb.Like("al.Title", pattern))
				case "Artist":
					sb.Where(sb.Like("ar.Name", pattern))
				case "Composer":
					sb.Where(sb.Like("t.Composer", pattern))
				default:
					sb.Where(sb.Like("t.Name", pattern))
				}
			}

			// the column dropdowns, one IN per checked column
			mapping := map[string]string{
				"genre": "t.GenreId",
				"media": "t.MediaTypeId",
			}
			for key, values := range filters {
				if field, ok := mapping[key]; ok && len(values) > 0 {
					sb.Where(sb.In(field, sqlbuilder.List(values)))
				}
			}

			sb.OrderBy("t.Name", "t.TrackId")
			sb.Limit(TrackPageSize)
			sb.Offset(req.Offset)
		},
		Scan: func(rows *sql.Rows, record *Track) error {
			return rows.Scan(
				&record.ID,
				&record.Name,
				&record.Album,
				&record.Artist,
				&record.Genre,
				&record.Media,
				&record.Composer,
				&record.Length,
				&record.Bytes,
				&record.Price,
			)
		},
	})

	return records, wrap.IfErr(op, err)
}

func parseFilters(raw string) (map[string][]string, error) {
	if raw == "" {
		return nil, nil
	}
	var filters map[string][]string
	if err := json.Unmarshal([]byte(raw), &filters); err != nil {
		return nil, err
	}
	return filters, nil
}
