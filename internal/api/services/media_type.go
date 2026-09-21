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
	ErrMediaTypeExists     = errors.New("media type with the same name already exists")
	ErrMediaTypeNotFound   = errors.New("media type not found")
	ErrMediaTypeHasRelated = errors.New("media type still has tracks")
)

type MediaType struct {
	ID     int              `json:"id"`
	Name   w2.Field[string] `json:"name"`
	Tracks int              `json:"tracks"`
}

func (s *Service) GetMediaTypeGrid(ctx context.Context, req w2.GetGridRequest) (w2.GetGridResponse[MediaType], error) {
	const op = "api.Service.GetMediaTypeGrid"
	res, err := w2db.GetGridContext(ctx, s.store.DB(), req, w2db.GetGridOptions[MediaType]{
		From: "media_types as t",
		Select: []string{
			"t.MediaTypeId as id",
			"t.Name as name",
			"coalesce(r.count, 0) as tracks",
		},
		Where: map[string]string{
			"name": "t.Name",
		},
		OrderBy: map[string]string{
			"id":     "t.MediaTypeId",
			"name":   "t.Name",
			"tracks": "r.count",
		},
		Build: func(sb *sqlbuilder.SelectBuilder) {
			sb.JoinWithOption(sqlbuilder.LeftJoin,
				"(select MediaTypeId, count(*) as count from tracks group by MediaTypeId) as r",
				"r.MediaTypeId = t.MediaTypeId",
			)
		},
		Scan: func(rows *sql.Rows, record *MediaType) error {
			return rows.Scan(
				&record.ID,
				&record.Name,
				&record.Tracks,
			)
		},
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) GetMediaTypeDropdown(ctx context.Context, req w2.GetDropdownRequest) (w2.GetDropdownResponse[w2.Dropdown], error) {
	const op = "api.Service.GetMediaTypeDropdown"
	res, err := w2db.GetDropdownContext(ctx, s.store.DB(), req, w2db.GetDropdownOptions{
		From:         "media_types",
		IDField:      "MediaTypeId",
		TextField:    "Name",
		OrderByField: "Name",
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) CreateMediaType(ctx context.Context, req w2.SaveFormRequest[MediaType]) (int, error) {
	const op = "api.Service.CreateMediaType"
	id, err := w2db.InsertContext(ctx, s.store.DB(), w2db.InsertOptions{
		Into: "media_types",
		Values: map[string]any{
			"Name": req.Record.Name,
		},
	})
	if s.store.IsErrConstraintUnique(err) {
		return 0, wrap.IfErr(op, ErrMediaTypeExists)
	}
	return id, wrap.IfErr(op, err)
}

func (s *Service) UpdateMediaType(ctx context.Context, req w2.SaveFormRequest[MediaType]) error {
	const op = "api.Service.UpdateMediaType"
	affected, err := w2db.UpdateContext(ctx, s.store.DB(), w2db.UpdateOptions{
		Update: "media_types",
		Values: map[string]any{
			"Name": req.Record.Name,
		},
		Where: map[string]any{
			"MediaTypeId": req.Record.ID,
		},
	})
	if s.store.IsErrConstraintUnique(err) {
		return wrap.IfErr(op, ErrMediaTypeExists)
	} else if err != nil {
		return wrap.IfErr(op, err)
	} else if affected == 0 {
		return wrap.IfErr(op, ErrMediaTypeNotFound)
	}
	return nil
}

func (s *Service) DeleteMediaTypes(ctx context.Context, req w2.RemoveGridRequest) error {
	const op = "api.Service.DeleteMediaTypes"
	affected, err := w2db.RemoveGridContext(ctx, s.store.DB(), req, w2db.RemoveGridOptions{
		From:    "media_types",
		IDField: "MediaTypeId",
	})
	if s.store.IsErrConstraintForeignKey(err) {
		return wrap.IfErr(op, ErrMediaTypeHasRelated)
	} else if err != nil {
		return wrap.IfErr(op, err)
	} else if affected != len(req.ID) {
		return wrap.IfErr(op, ErrMediaTypeNotFound)
	}
	return nil
}
