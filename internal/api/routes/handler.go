package routes

import "github.com/dv1x3r/chinook-go/internal/api/services"

type Handler struct {
	svc *services.Service
}

func NewHandler(svc *services.Service) *Handler {
	return &Handler{
		svc: svc,
	}
}
