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
	ErrInvoiceNotFound = errors.New("invoice not found")
)

type Invoice struct {
	ID          int               `json:"id"`
	Customer    w2.Dropdown       `json:"customer"`
	InvoiceDate w2.Field[string]  `json:"invoice_date"`
	Address     w2.Field[string]  `json:"billing_address"`
	City        w2.Field[string]  `json:"billing_city"`
	State       w2.Field[string]  `json:"billing_state"`
	Country     w2.Field[string]  `json:"billing_country"`
	PostalCode  w2.Field[string]  `json:"billing_postal_code"`
	Total       w2.Field[float64] `json:"total"`
}

func invoiceValues(record Invoice) map[string]any {
	return map[string]any{
		"CustomerId":        record.Customer.ID,
		"InvoiceDate":       record.InvoiceDate,
		"BillingAddress":    record.Address,
		"BillingCity":       record.City,
		"BillingState":      record.State,
		"BillingCountry":    record.Country,
		"BillingPostalCode": record.PostalCode,
	}
}

func (s *Service) GetInvoiceGrid(ctx context.Context, req w2.GetGridRequest) (w2.GetGridResponse[Invoice], error) {
	const op = "api.Service.GetInvoiceGrid"
	res, err := w2db.GetGridContext(ctx, s.store.DB(), req, w2db.GetGridOptions[Invoice]{
		From: "invoices as i",
		Select: []string{
			"i.InvoiceId as id",
			"i.CustomerId as customer_id",
			"concat(c.FirstName, ' ', c.LastName) as customer_text",
			"i.InvoiceDate as invoice_date",
			"i.BillingAddress as billing_address",
			"i.BillingCity as billing_city",
			"i.BillingState as billing_state",
			"i.BillingCountry as billing_country",
			"i.BillingPostalCode as billing_postal_code",
			"i.Total as total",
		},
		Where: map[string]string{
			"customer":        "i.CustomerId",
			"billing_city":    "i.BillingCity",
			"billing_country": "i.BillingCountry",
		},
		OrderBy: map[string]string{
			"id":              "i.InvoiceId",
			"customer":        "c.LastName",
			"invoice_date":    "i.InvoiceDate",
			"billing_city":    "i.BillingCity",
			"billing_country": "i.BillingCountry",
			"total":           "i.Total",
		},
		Build: func(sb *sqlbuilder.SelectBuilder) {
			sb.JoinWithOption(sqlbuilder.InnerJoin, "customers as c", "c.CustomerId = i.CustomerId")
		},
		Scan: func(rows *sql.Rows, record *Invoice) error {
			return rows.Scan(
				&record.ID,
				&record.Customer.ID,
				&record.Customer.Text,
				&record.InvoiceDate,
				&record.Address,
				&record.City,
				&record.State,
				&record.Country,
				&record.PostalCode,
				&record.Total,
			)
		},
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) CreateInvoice(ctx context.Context, req w2.SaveFormRequest[Invoice]) (int, error) {
	const op = "api.Service.CreateInvoice"
	values := invoiceValues(req.Record)
	values["Total"] = 0 // a fresh invoice has no lines yet
	id, err := w2db.InsertContext(ctx, s.store.DB(), w2db.InsertOptions{
		Into:   "invoices",
		Values: values,
	})
	return id, wrap.IfErr(op, err)
}

func (s *Service) UpdateInvoice(ctx context.Context, req w2.SaveFormRequest[Invoice]) error {
	const op = "api.Service.UpdateInvoice"
	affected, err := w2db.UpdateContext(ctx, s.store.DB(), w2db.UpdateOptions{
		Update: "invoices",
		Values: invoiceValues(req.Record),
		Where: map[string]any{
			"InvoiceId": req.Record.ID,
		},
	})
	if err != nil {
		return wrap.IfErr(op, err)
	} else if affected == 0 {
		return wrap.IfErr(op, ErrInvoiceNotFound)
	}
	return nil
}

func (s *Service) DeleteInvoices(ctx context.Context, req w2.RemoveGridRequest) error {
	const op = "api.Service.DeleteInvoices"

	tx, err := s.store.DB().BeginTx(ctx, nil)
	if err != nil {
		return wrap.IfErr(op, err)
	}
	defer tx.Rollback()

	// the sample database has no cascade on the lines, an invoice takes them along
	for _, id := range req.ID {
		if _, err := tx.ExecContext(ctx, "delete from invoice_items where InvoiceId = ?;", id); err != nil {
			return wrap.IfErr(op, err)
		}
	}

	affected, err := w2db.RemoveGridContext(ctx, tx, req, w2db.RemoveGridOptions{
		From:    "invoices",
		IDField: "InvoiceId",
	})
	if err != nil {
		return wrap.IfErr(op, err)
	} else if affected != len(req.ID) {
		return wrap.IfErr(op, ErrInvoiceNotFound)
	}

	return wrap.IfErr(op, tx.Commit())
}
