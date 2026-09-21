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
	ErrArtistExists     = errors.New("artist with the same name already exists")
	ErrArtistNotFound   = errors.New("artist not found")
	ErrArtistHasRelated = errors.New("artist still has albums")
)

type Artist struct {
	ID     int              `json:"id"`
	Name   w2.Field[string] `json:"name"`
	Albums int              `json:"albums"`
}

func (s *Service) GetArtistGrid(ctx context.Context, req w2.GetGridRequest) (w2.GetGridResponse[Artist], error) {
	const op = "api.Service.GetArtistGrid"
	res, err := w2db.GetGridContext(ctx, s.store.DB(), req, w2db.GetGridOptions[Artist]{
		From: "artists as t",
		Select: []string{
			"t.ArtistId as id",
			"t.Name as name",
			"coalesce(r.count, 0) as albums",
		},
		Where: map[string]string{
			"name": "t.Name",
		},
		OrderBy: map[string]string{
			"id":     "t.ArtistId",
			"name":   "t.Name",
			"albums": "r.count",
		},
		Build: func(sb *sqlbuilder.SelectBuilder) {
			sb.JoinWithOption(sqlbuilder.LeftJoin,
				"(select ArtistId, count(*) as count from albums group by ArtistId) as r",
				"r.ArtistId = t.ArtistId",
			)
		},
		Scan: func(rows *sql.Rows, record *Artist) error {
			return rows.Scan(
				&record.ID,
				&record.Name,
				&record.Albums,
			)
		},
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) GetArtistDropdown(ctx context.Context, req w2.GetDropdownRequest) (w2.GetDropdownResponse[w2.Dropdown], error) {
	const op = "api.Service.GetArtistDropdown"
	res, err := w2db.GetDropdownContext(ctx, s.store.DB(), req, w2db.GetDropdownOptions{
		From:         "artists",
		IDField:      "ArtistId",
		TextField:    "Name",
		OrderByField: "Name",
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) CreateArtist(ctx context.Context, req w2.SaveFormRequest[Artist]) (int, error) {
	const op = "api.Service.CreateArtist"
	id, err := w2db.InsertContext(ctx, s.store.DB(), w2db.InsertOptions{
		Into: "artists",
		Values: map[string]any{
			"Name": req.Record.Name,
		},
	})
	if s.store.IsErrConstraintUnique(err) {
		return 0, wrap.IfErr(op, ErrArtistExists)
	}
	return id, wrap.IfErr(op, err)
}

func (s *Service) UpdateArtist(ctx context.Context, req w2.SaveFormRequest[Artist]) error {
	const op = "api.Service.UpdateArtist"
	affected, err := w2db.UpdateContext(ctx, s.store.DB(), w2db.UpdateOptions{
		Update: "artists",
		Values: map[string]any{
			"Name": req.Record.Name,
		},
		Where: map[string]any{
			"ArtistId": req.Record.ID,
		},
	})
	if s.store.IsErrConstraintUnique(err) {
		return wrap.IfErr(op, ErrArtistExists)
	} else if err != nil {
		return wrap.IfErr(op, err)
	} else if affected == 0 {
		return wrap.IfErr(op, ErrArtistNotFound)
	}
	return nil
}

func (s *Service) DeleteArtists(ctx context.Context, req w2.RemoveGridRequest) error {
	const op = "api.Service.DeleteArtists"
	affected, err := w2db.RemoveGridContext(ctx, s.store.DB(), req, w2db.RemoveGridOptions{
		From:    "artists",
		IDField: "ArtistId",
	})
	if s.store.IsErrConstraintForeignKey(err) {
		return wrap.IfErr(op, ErrArtistHasRelated)
	} else if err != nil {
		return wrap.IfErr(op, err)
	} else if affected != len(req.ID) {
		return wrap.IfErr(op, ErrArtistNotFound)
	}
	return nil
}
