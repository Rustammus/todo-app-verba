package v1

import (
	"ToDoVerba/internal/service"
	"ToDoVerba/pkg/logging"
	"github.com/julienschmidt/httprouter"
)

type Handler struct {
	service service.Services
	logger  logging.Logger
}

type Deps struct {
	Service service.Services
	Logger  logging.Logger
}

func NewHandler(d Deps) *Handler {
	return &Handler{
		service: d.Service,
		logger:  d.Logger,
	}
}

func (h *Handler) Init(r *httprouter.Router) {
	h.initTaskHandler(r)
}
