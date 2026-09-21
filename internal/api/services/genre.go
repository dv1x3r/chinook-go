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
	ErrGenreExists     = errors.New("genre with the same name already exists")
	ErrGenreNotFound   = errors.New("genre not found")
	ErrGenreHasRelated = errors.New("genre still has tracks")
)

type Genre struct {
	ID     int              `json:"id"`
	Name   w2.Field[string] `json:"name"`
	Tracks int              `json:"tracks"`
}

func (s *Service) GetGenreGrid(ctx context.Context, req w2.GetGridRequest) (w2.GetGridResponse[Genre], error) {
	const op = "api.Service.GetGenreGrid"
	res, err := w2db.GetGridContext(ctx, s.store.DB(), req, w2db.GetGridOptions[Genre]{
		From: "genres as t",
		Select: []string{
			"t.GenreId as id",
			"t.Name as name",
			"coalesce(r.count, 0) as tracks",
		},
		Where: map[string]string{
			"name": "t.Name",
		},
		OrderBy: map[string]string{
			"id":     "t.GenreId",
			"name":   "t.Name",
			"tracks": "r.count",
		},
		Build: func(sb *sqlbuilder.SelectBuilder) {
			sb.JoinWithOption(sqlbuilder.LeftJoin,
				"(select GenreId, count(*) as count from tracks group by GenreId) as r",
				"r.GenreId = t.GenreId",
			)
		},
		Scan: func(rows *sql.Rows, record *Genre) error {
			return rows.Scan(
				&record.ID,
				&record.Name,
				&record.Tracks,
			)
		},
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) GetGenreDropdown(ctx context.Context, req w2.GetDropdownRequest) (w2.GetDropdownResponse[w2.Dropdown], error) {
	const op = "api.Service.GetGenreDropdown"
	res, err := w2db.GetDropdownContext(ctx, s.store.DB(), req, w2db.GetDropdownOptions{
		From:         "genres",
		IDField:      "GenreId",
		TextField:    "Name",
		OrderByField: "Name",
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) CreateGenre(ctx context.Context, req w2.SaveFormRequest[Genre]) (int, error) {
	const op = "api.Service.CreateGenre"
	id, err := w2db.InsertContext(ctx, s.store.DB(), w2db.InsertOptions{
		Into: "genres",
		Values: map[string]any{
			"Name": req.Record.Name,
		},
	})
	if s.store.IsErrConstraintUnique(err) {
		return 0, wrap.IfErr(op, ErrGenreExists)
	}
	return id, wrap.IfErr(op, err)
}

func (s *Service) UpdateGenre(ctx context.Context, req w2.SaveFormRequest[Genre]) error {
	const op = "api.Service.UpdateGenre"
	affected, err := w2db.UpdateContext(ctx, s.store.DB(), w2db.UpdateOptions{
		Update: "genres",
		Values: map[string]any{
			"Name": req.Record.Name,
		},
		Where: map[string]any{
			"GenreId": req.Record.ID,
		},
	})
	if s.store.IsErrConstraintUnique(err) {
		return wrap.IfErr(op, ErrGenreExists)
	} else if err != nil {
		return wrap.IfErr(op, err)
	} else if affected == 0 {
		return wrap.IfErr(op, ErrGenreNotFound)
	}
	return nil
}

func (s *Service) DeleteGenres(ctx context.Context, req w2.RemoveGridRequest) error {
	const op = "api.Service.DeleteGenres"
	affected, err := w2db.RemoveGridContext(ctx, s.store.DB(), req, w2db.RemoveGridOptions{
		From:    "genres",
		IDField: "GenreId",
	})
	if s.store.IsErrConstraintForeignKey(err) {
		return wrap.IfErr(op, ErrGenreHasRelated)
	} else if err != nil {
		return wrap.IfErr(op, err)
	} else if affected != len(req.ID) {
		return wrap.IfErr(op, ErrGenreNotFound)
	}
	return nil
}
