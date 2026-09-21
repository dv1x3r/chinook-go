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
	ErrCustomerNotFound   = errors.New("customer not found")
	ErrCustomerHasRelated = errors.New("customer still has invoices")
)

type Customer struct {
	ID         int              `json:"id"`
	FirstName  w2.Field[string] `json:"first_name"`
	LastName   w2.Field[string] `json:"last_name"`
	FullName   string           `json:"full_name"`
	Company    w2.Field[string] `json:"company"`
	Address    w2.Field[string] `json:"address"`
	City       w2.Field[string] `json:"city"`
	State      w2.Field[string] `json:"state"`
	Country    w2.Field[string] `json:"country"`
	PostalCode w2.Field[string] `json:"postal_code"`
	Phone      w2.Field[string] `json:"phone"`
	Fax        w2.Field[string] `json:"fax"`
	Email      w2.Field[string] `json:"email"`
	SupportRep w2.Dropdown      `json:"support_rep"`
	Invoices   int              `json:"invoices"`
}

func customerValues(record Customer) map[string]any {
	return map[string]any{
		"FirstName":    record.FirstName,
		"LastName":     record.LastName,
		"Company":      record.Company,
		"Address":      record.Address,
		"City":         record.City,
		"State":        record.State,
		"Country":      record.Country,
		"PostalCode":   record.PostalCode,
		"Phone":        record.Phone,
		"Fax":          record.Fax,
		"Email":        record.Email,
		"SupportRepId": record.SupportRep.ID,
	}
}

func (s *Service) GetCustomerGrid(ctx context.Context, req w2.GetGridRequest) (w2.GetGridResponse[Customer], error) {
	const op = "api.Service.GetCustomerGrid"
	res, err := w2db.GetGridContext(ctx, s.store.DB(), req, w2db.GetGridOptions[Customer]{
		From: "customers as c",
		Select: []string{
			"c.CustomerId as id",
			"c.FirstName as first_name",
			"c.LastName as last_name",
			"concat(c.FirstName, ' ', c.LastName) as full_name",
			"c.Company as company",
			"c.Address as address",
			"c.City as city",
			"c.State as state",
			"c.Country as country",
			"c.PostalCode as postal_code",
			"c.Phone as phone",
			"c.Fax as fax",
			"c.Email as email",
			"c.SupportRepId as support_rep_id",
			"concat(e.FirstName, ' ', e.LastName) as support_rep_text",
			"coalesce(i.count, 0) as invoices",
		},
		Where: map[string]string{
			"full_name": "concat(c.FirstName, ' ', c.LastName)",
			"company":   "c.Company",
			"phone":     "c.Phone",
			"fax":       "c.Fax",
			"email":     "c.Email",
		},
		OrderBy: map[string]string{
			"id":        "c.CustomerId",
			"full_name": "c.LastName",
			"company":   "c.Company",
			"phone":     "c.Phone",
			"fax":       "c.Fax",
			"email":     "c.Email",
			"invoices":  "i.count",
		},
		Build: func(sb *sqlbuilder.SelectBuilder) {
			sb.JoinWithOption(sqlbuilder.LeftJoin, "employees as e", "e.EmployeeId = c.SupportRepId")
			sb.JoinWithOption(sqlbuilder.LeftJoin,
				"(select CustomerId, count(*) as count from invoices group by CustomerId) as i",
				"i.CustomerId = c.CustomerId",
			)
		},
		Scan: func(rows *sql.Rows, record *Customer) error {
			return rows.Scan(
				&record.ID,
				&record.FirstName,
				&record.LastName,
				&record.FullName,
				&record.Company,
				&record.Address,
				&record.City,
				&record.State,
				&record.Country,
				&record.PostalCode,
				&record.Phone,
				&record.Fax,
				&record.Email,
				&record.SupportRep.ID,
				&record.SupportRep.Text,
				&record.Invoices,
			)
		},
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) GetCustomerDropdown(ctx context.Context, req w2.GetDropdownRequest) (w2.GetDropdownResponse[w2.Dropdown], error) {
	const op = "api.Service.GetCustomerDropdown"
	res, err := w2db.GetDropdownContext(ctx, s.store.DB(), req, w2db.GetDropdownOptions{
		From:         "customers",
		IDField:      "CustomerId",
		TextField:    "concat(FirstName, ' ', LastName)",
		OrderByField: "LastName, FirstName",
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) CreateCustomer(ctx context.Context, req w2.SaveFormRequest[Customer]) (int, error) {
	const op = "api.Service.CreateCustomer"
	id, err := w2db.InsertContext(ctx, s.store.DB(), w2db.InsertOptions{
		Into:   "customers",
		Values: customerValues(req.Record),
	})
	return id, wrap.IfErr(op, err)
}

func (s *Service) UpdateCustomer(ctx context.Context, req w2.SaveFormRequest[Customer]) error {
	const op = "api.Service.UpdateCustomer"
	affected, err := w2db.UpdateContext(ctx, s.store.DB(), w2db.UpdateOptions{
		Update: "customers",
		Values: customerValues(req.Record),
		Where: map[string]any{
			"CustomerId": req.Record.ID,
		},
	})
	if err != nil {
		return wrap.IfErr(op, err)
	} else if affected == 0 {
		return wrap.IfErr(op, ErrCustomerNotFound)
	}
	return nil
}

func (s *Service) DeleteCustomers(ctx context.Context, req w2.RemoveGridRequest) error {
	const op = "api.Service.DeleteCustomers"
	affected, err := w2db.RemoveGridContext(ctx, s.store.DB(), req, w2db.RemoveGridOptions{
		From:    "customers",
		IDField: "CustomerId",
	})
	if s.store.IsErrConstraintForeignKey(err) {
		return wrap.IfErr(op, ErrCustomerHasRelated)
	} else if err != nil {
		return wrap.IfErr(op, err)
	} else if affected != len(req.ID) {
		return wrap.IfErr(op, ErrCustomerNotFound)
	}
	return nil
}
