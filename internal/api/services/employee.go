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
	ErrEmployeeNotFound   = errors.New("employee not found")
	ErrEmployeeHasRelated = errors.New("employee still reports or supports customers")
)

type Employee struct {
	ID         int              `json:"id"`
	FirstName  w2.Field[string] `json:"first_name"`
	LastName   w2.Field[string] `json:"last_name"`
	FullName   string           `json:"full_name"`
	Title      w2.Field[string] `json:"title"`
	ReportsTo  w2.Dropdown      `json:"reports_to"`
	BirthDate  w2.Field[string] `json:"birth_date"`
	HireDate   w2.Field[string] `json:"hire_date"`
	Address    w2.Field[string] `json:"address"`
	City       w2.Field[string] `json:"city"`
	State      w2.Field[string] `json:"state"`
	Country    w2.Field[string] `json:"country"`
	PostalCode w2.Field[string] `json:"postal_code"`
	Phone      w2.Field[string] `json:"phone"`
	Fax        w2.Field[string] `json:"fax"`
	Email      w2.Field[string] `json:"email"`
	Customers  int              `json:"customers"`
}

func employeeValues(record Employee) map[string]any {
	return map[string]any{
		"FirstName":  record.FirstName,
		"LastName":   record.LastName,
		"Title":      record.Title,
		"ReportsTo":  record.ReportsTo.ID,
		"BirthDate":  record.BirthDate,
		"HireDate":   record.HireDate,
		"Address":    record.Address,
		"City":       record.City,
		"State":      record.State,
		"Country":    record.Country,
		"PostalCode": record.PostalCode,
		"Phone":      record.Phone,
		"Fax":        record.Fax,
		"Email":      record.Email,
	}
}

func (s *Service) GetEmployeeGrid(ctx context.Context, req w2.GetGridRequest) (w2.GetGridResponse[Employee], error) {
	const op = "api.Service.GetEmployeeGrid"
	res, err := w2db.GetGridContext(ctx, s.store.DB(), req, w2db.GetGridOptions[Employee]{
		From: "employees as e",
		Select: []string{
			"e.EmployeeId as id",
			"e.FirstName as first_name",
			"e.LastName as last_name",
			"concat(e.FirstName, ' ', e.LastName) as full_name",
			"e.Title as title",
			"e.ReportsTo as reports_to_id",
			"concat(m.FirstName, ' ', m.LastName) as reports_to_text",
			"e.BirthDate as birth_date",
			"e.HireDate as hire_date",
			"e.Address as address",
			"e.City as city",
			"e.State as state",
			"e.Country as country",
			"e.PostalCode as postal_code",
			"e.Phone as phone",
			"e.Fax as fax",
			"e.Email as email",
			"coalesce(c.count, 0) as customers",
		},
		Where: map[string]string{
			"full_name": "concat(e.FirstName, ' ', e.LastName)",
			"title":     "e.Title",
		},
		OrderBy: map[string]string{
			"id":         "e.EmployeeId",
			"full_name":  "e.LastName",
			"title":      "e.Title",
			"birth_date": "e.BirthDate",
			"hire_date":  "e.HireDate",
			"reports_to": "m.LastName",
			"customers":  "c.count",
		},
		Build: func(sb *sqlbuilder.SelectBuilder) {
			sb.JoinWithOption(sqlbuilder.LeftJoin, "employees as m", "m.EmployeeId = e.ReportsTo")
			sb.JoinWithOption(sqlbuilder.LeftJoin,
				"(select SupportRepId, count(*) as count from customers group by SupportRepId) as c",
				"c.SupportRepId = e.EmployeeId",
			)
		},
		Scan: func(rows *sql.Rows, record *Employee) error {
			return rows.Scan(
				&record.ID,
				&record.FirstName,
				&record.LastName,
				&record.FullName,
				&record.Title,
				&record.ReportsTo.ID,
				&record.ReportsTo.Text,
				&record.BirthDate,
				&record.HireDate,
				&record.Address,
				&record.City,
				&record.State,
				&record.Country,
				&record.PostalCode,
				&record.Phone,
				&record.Fax,
				&record.Email,
				&record.Customers,
			)
		},
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) GetEmployeeDropdown(ctx context.Context, req w2.GetDropdownRequest) (w2.GetDropdownResponse[w2.Dropdown], error) {
	const op = "api.Service.GetEmployeeDropdown"
	res, err := w2db.GetDropdownContext(ctx, s.store.DB(), req, w2db.GetDropdownOptions{
		From:         "employees",
		IDField:      "EmployeeId",
		TextField:    "concat(FirstName, ' ', LastName)",
		OrderByField: "LastName, FirstName",
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) CreateEmployee(ctx context.Context, req w2.SaveFormRequest[Employee]) (int, error) {
	const op = "api.Service.CreateEmployee"
	id, err := w2db.InsertContext(ctx, s.store.DB(), w2db.InsertOptions{
		Into:   "employees",
		Values: employeeValues(req.Record),
	})
	return id, wrap.IfErr(op, err)
}

func (s *Service) UpdateEmployee(ctx context.Context, req w2.SaveFormRequest[Employee]) error {
	const op = "api.Service.UpdateEmployee"
	affected, err := w2db.UpdateContext(ctx, s.store.DB(), w2db.UpdateOptions{
		Update: "employees",
		Values: employeeValues(req.Record),
		Where: map[string]any{
			"EmployeeId": req.Record.ID,
		},
	})
	if err != nil {
		return wrap.IfErr(op, err)
	} else if affected == 0 {
		return wrap.IfErr(op, ErrEmployeeNotFound)
	}
	return nil
}

func (s *Service) DeleteEmployees(ctx context.Context, req w2.RemoveGridRequest) error {
	const op = "api.Service.DeleteEmployees"
	affected, err := w2db.RemoveGridContext(ctx, s.store.DB(), req, w2db.RemoveGridOptions{
		From:    "employees",
		IDField: "EmployeeId",
	})
	if s.store.IsErrConstraintForeignKey(err) {
		return wrap.IfErr(op, ErrEmployeeHasRelated)
	} else if err != nil {
		return wrap.IfErr(op, err)
	} else if affected != len(req.ID) {
		return wrap.IfErr(op, ErrEmployeeNotFound)
	}
	return nil
}
