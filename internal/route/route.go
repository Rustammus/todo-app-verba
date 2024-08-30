package route

import (
	v1 "ToDoVerba/internal/route/api/v1"
	"ToDoVerba/internal/service"
	"ToDoVerba/pkg/logging"
	"github.com/julienschmidt/httprouter"
)

type Handler struct {
	services service.Services //TODO
	logger   logging.Logger
}

type Deps struct {
	Services service.Services //TODO
	Logger   logging.Logger
}

func NewHandler(d Deps) *Handler {
	return &Handler{services: d.Services, logger: d.Logger}
}

func (h *Handler) Init(r *httprouter.Router) {
	hv1 := v1.NewHandler(v1.Deps{
		Service: h.services,
		Logger:  h.logger,
	})
	hv1.Init(r)
}
