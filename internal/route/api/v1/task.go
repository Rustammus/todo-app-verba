package v1

import (
	"ToDoVerba/internal/schemas"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"strconv"
)

func (h *Handler) initTaskHandler(r *httprouter.Router) {
	r.POST("/tasks", h.taskCreate)
	r.GET("/tasks", h.taskList)
	r.GET("/tasks/:id", h.taskFindById)
	r.PUT("/tasks/:id", h.taskUpdateById)
	r.DELETE("/tasks/:id", h.taskDeleteById)
}

// taskCreate godoc
// @Tags         Task API
// @Summary      Create Task Summary
// @Description  Create Task Description
// @Accept       json
// @Produce      json
// @Param Task body schemas.RequestTaskCreate false "Task base"
// @Success      201  {object}  schemas.ResponseTaskRead
// @Failure      400  {object}  errorJSON
// @Failure      500  {object}	errorJSON
// @Router       /tasks [post]
func (h *Handler) taskCreate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	h.logger.Debugf("[%s] %s taskCreate called", r.Method, r.RemoteAddr)
	contentType := r.Header.Get("content-type")
	if contentType != "application/json" {
		writeResponseErr(w, http.StatusBadRequest,
			errors.New("content-type is not application/json"))
		return
	}

	cTask := schemas.RequestTaskCreate{}
	bodyRaw, err := io.ReadAll(r.Body)
	if err != nil {
		writeResponseErr(w, http.StatusInternalServerError, err)
		return
	}

	err = json.Unmarshal(bodyRaw, &cTask)
	if err != nil {
		writeResponseErr(w, http.StatusBadRequest, err)
		return
	}
	err = cTask.Valid()
	if err != nil {
		writeResponseErr(w, http.StatusBadRequest, err)
		return
	}

	rTaskDTO, err := h.service.Task.Create(cTask.ToDTO())
	if err != nil {
		writeResponseErr(w, http.StatusInternalServerError, err)
		return
	}

	rTask := schemas.ResponseTaskRead{}
	rTask.ScanDTO(rTaskDTO)

	writeResponse(w, http.StatusCreated, rTask)
}

// taskList godoc
// @Tags         Task API
// @Summary      List Task Summary
// @Description  List Task Description
// @Accept       json
// @Produce      json
// @Success      200  {object}  []schemas.ResponseTaskRead
// @Failure      500  {object}	errorJSON
// @Router       /tasks [get]
func (h *Handler) taskList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	h.logger.Debugf("[%s] %s taskList called", r.Method, r.RemoteAddr)

	rTasksDTO, err := h.service.Task.List()
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		writeResponseErr(w, http.StatusInternalServerError, err)
		return
	}

	rTasks := make([]schemas.ResponseTaskRead, 0, len(rTasksDTO))
	for i := 0; i < len(rTasksDTO); i++ {
		rTask := schemas.ResponseTaskRead{}
		rTask.ScanDTO(&rTasksDTO[i])
		rTasks = append(rTasks, rTask)
	}

	writeResponse(w, http.StatusOK, rTasks)
}

// taskFindById godoc
// @Tags         Task API
// @Summary      Find Task by id Summary
// @Description  Find Task by id Description
// @Accept       json
// @Produce      json
// @Param id path int false "Task id"
// @Success      200  {object}  schemas.ResponseTaskRead
// @Failure      400  {object}	errorJSON
// @Failure      404  {object}	errorJSON
// @Failure      500  {object}	errorJSON
// @Router       /tasks/{id} [get]
func (h *Handler) taskFindById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	h.logger.Debugf("[%s] %s taskFindById called", r.Method, r.RemoteAddr)

	idStr := ps.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeResponseErr(w, http.StatusBadRequest, err)
		return
	}

	rTaskDTO, err := h.service.Task.FindByID(id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeResponseErr(w, http.StatusNotFound, err)
			return
		}
		writeResponseErr(w, http.StatusInternalServerError, err)
		return
	}

	rTask := schemas.ResponseTaskRead{}
	rTask.ScanDTO(rTaskDTO)
	writeResponse(w, http.StatusOK, rTask)
}

// taskUpdateById godoc
// @Tags         Task API
// @Summary      Update Task Summary
// @Description  Update Task Description
// @Accept       json
// @Produce      json
// @Param id path int false "Task id"
// @Param Task body schemas.RequestTaskUpdate false "Task update"
// @Success      200  {object}  schemas.ResponseTaskRead
// @Failure      400  {object}  errorJSON
// @Failure      404  {object}  errorJSON
// @Failure      500  {object}	errorJSON
// @Router       /tasks/{id} [put]
func (h *Handler) taskUpdateById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	h.logger.Debugf("[%s] %s taskUpdateById called", r.Method, r.RemoteAddr)

	contentType := r.Header.Get("content-type")
	if contentType != "application/json" {
		writeResponseErr(w, http.StatusBadRequest,
			errors.New("content-type is not application/json"))
		return
	}

	idStr := ps.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeResponseErr(w, http.StatusBadRequest, err)
		return
	}

	uTask := schemas.RequestTaskUpdate{}
	bodyRaw, err := io.ReadAll(r.Body)
	if err != nil {
		writeResponseErr(w, http.StatusInternalServerError, err)
		return
	}

	err = json.Unmarshal(bodyRaw, &uTask)
	if err != nil {
		writeResponseErr(w, http.StatusBadRequest, err)
		return
	}
	err = uTask.Valid()
	if err != nil {
		writeResponseErr(w, http.StatusBadRequest, err)
		return
	}

	rTaskDTO, err := h.service.Task.UpdateById(id, uTask.ToDTO())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeResponseErr(w, http.StatusNotFound, err)
			return
		}
		writeResponseErr(w, http.StatusInternalServerError, err)
		return
	}

	rTask := schemas.ResponseTaskRead{}
	rTask.ScanDTO(rTaskDTO)
	writeResponse(w, http.StatusOK, rTask)
}

// taskDeleteById godoc
// @Tags         Task API
// @Summary      Delete Task by id Summary
// @Description  Delete Task by id Description
// @Accept       json
// @Produce      json
// @Param id path int false "Task id"
// @Success      204  {object}  schemas.ResponseTaskRead
// @Failure      400  {object}	errorJSON
// @Failure      404  {object}	errorJSON
// @Failure      500  {object}	errorJSON
// @Router       /tasks/{id} [delete]
func (h *Handler) taskDeleteById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	h.logger.Debugf("[%s] %s taskDeleteById called", r.Method, r.RemoteAddr)

	idStr := ps.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeResponseErr(w, http.StatusBadRequest, err)
		return
	}

	err = h.service.Task.DeleteById(id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeResponseErr(w, http.StatusNotFound, err)
			return
		}
		writeResponseErr(w, http.StatusInternalServerError, err)
		return
	}

	writeResponse(w, http.StatusNoContent, nil)
}
