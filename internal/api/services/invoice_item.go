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
	ErrInvoiceItemNotFound = errors.New("invoice line not found")
)

type InvoiceItem struct {
	ID        int               `json:"id"`
	InvoiceID int               `json:"invoice_id"`
	Track     w2.Dropdown       `json:"track"`
	Quantity  w2.Field[int]     `json:"quantity"`
	UnitPrice w2.Field[float64] `json:"unit_price"`
	Total     float64           `json:"total"`
}

func (s *Service) GetInvoiceItemGrid(ctx context.Context, req w2.GetGridRequest, invoiceID int) (w2.GetGridResponse[InvoiceItem], error) {
	const op = "api.Service.GetInvoiceItemGrid"
	res, err := w2db.GetGridContext(ctx, s.store.DB(), req, w2db.GetGridOptions[InvoiceItem]{
		From: "invoice_items as l",
		Select: []string{
			"l.InvoiceLineId as id",
			"l.InvoiceId as invoice_id",
			"l.TrackId as track_id",
			"t.Name as track_text",
			"l.Quantity as quantity",
			"l.UnitPrice as unit_price",
			"l.UnitPrice * l.Quantity as total",
		},
		Build: func(sb *sqlbuilder.SelectBuilder) {
			sb.JoinWithOption(sqlbuilder.InnerJoin, "tracks as t", "t.TrackId = l.TrackId")
			sb.Where(sb.EQ("l.InvoiceId", invoiceID))
			sb.OrderBy("t.Name")
		},
		Scan: func(rows *sql.Rows, record *InvoiceItem) error {
			return rows.Scan(
				&record.ID,
				&record.InvoiceID,
				&record.Track.ID,
				&record.Track.Text,
				&record.Quantity,
				&record.UnitPrice,
				&record.Total,
			)
		},
	})
	if err != nil {
		return res, wrap.IfErr(op, err)
	}

	var total float64
	err = s.store.DB().QueryRowContext(ctx, "select Total from invoices where InvoiceId = ?;", invoiceID).Scan(&total)
	if errors.Is(err, sql.ErrNoRows) {
		return res, nil
	} else if err != nil {
		return res, wrap.IfErr(op, err)
	}
	res.Summary = []InvoiceItem{{Total: total}}
	return res, nil
}

func (s *Service) CreateInvoiceItem(ctx context.Context, req w2.SaveFormRequest[InvoiceItem], invoiceID int) (int, error) {
	const op = "api.Service.CreateInvoiceItem"

	tx, err := s.store.DB().BeginTx(ctx, nil)
	if err != nil {
		return 0, wrap.IfErr(op, err)
	}
	defer tx.Rollback()

	id, err := w2db.InsertContext(ctx, tx, w2db.InsertOptions{
		Into: "invoice_items",
		Values: map[string]any{
			"InvoiceId": invoiceID,
			"TrackId":   req.Record.Track.ID,
			"Quantity":  req.Record.Quantity,
			"UnitPrice": req.Record.UnitPrice,
		},
	})
	if err != nil {
		return 0, wrap.IfErr(op, err)
	}

	if err := recalcInvoiceTotal(ctx, tx, invoiceID); err != nil {
		return 0, wrap.IfErr(op, err)
	}

	return id, wrap.IfErr(op, tx.Commit())
}

func (s *Service) UpdateInvoiceItem(ctx context.Context, req w2.SaveFormRequest[InvoiceItem], invoiceID int) error {
	const op = "api.Service.UpdateInvoiceItem"

	tx, err := s.store.DB().BeginTx(ctx, nil)
	if err != nil {
		return wrap.IfErr(op, err)
	}
	defer tx.Rollback()

	affected, err := w2db.UpdateContext(ctx, tx, w2db.UpdateOptions{
		Update: "invoice_items",
		Values: map[string]any{
			"TrackId":   req.Record.Track.ID,
			"Quantity":  req.Record.Quantity,
			"UnitPrice": req.Record.UnitPrice,
		},
		Where: map[string]any{
			"InvoiceLineId": req.Record.ID,
			"InvoiceId":     invoiceID,
		},
	})
	if err != nil {
		return wrap.IfErr(op, err)
	} else if affected == 0 {
		return wrap.IfErr(op, ErrInvoiceItemNotFound)
	}

	if err := recalcInvoiceTotal(ctx, tx, invoiceID); err != nil {
		return wrap.IfErr(op, err)
	}

	return wrap.IfErr(op, tx.Commit())
}

func (s *Service) DeleteInvoiceItems(ctx context.Context, req w2.RemoveGridRequest, invoiceID int) error {
	const op = "api.Service.DeleteInvoiceItems"

	if len(req.ID) == 0 {
		return nil
	}

	tx, err := s.store.DB().BeginTx(ctx, nil)
	if err != nil {
		return wrap.IfErr(op, err)
	}
	defer tx.Rollback()

	db := sqlbuilder.DeleteFrom("invoice_items")
	db.Where(
		db.EQ("InvoiceId", invoiceID),
		db.In("InvoiceLineId", sqlbuilder.Flatten(req.ID)...),
	)

	query, args := db.BuildWithFlavor(sqlbuilder.SQLite)
	res, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return wrap.IfErr(op, err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return wrap.IfErr(op, err)
	} else if int(affected) != len(req.ID) {
		return wrap.IfErr(op, ErrInvoiceItemNotFound)
	}

	if err := recalcInvoiceTotal(ctx, tx, invoiceID); err != nil {
		return wrap.IfErr(op, err)
	}

	return wrap.IfErr(op, tx.Commit())
}

// the stored total of the invoice follows its lines
func recalcInvoiceTotal(ctx context.Context, tx *sql.Tx, invoiceID int) error {
	_, err := tx.ExecContext(ctx, `
		update invoices set Total = coalesce(round((
			select sum(UnitPrice * Quantity) from invoice_items where InvoiceId = ?
		), 2), 0) where InvoiceId = ?;`, invoiceID, invoiceID)
	return err
}
