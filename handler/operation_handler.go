package handler

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"io"
	"net/http"
)

type OperationInputData struct {
	Year   int    `json:"year"`
	Month  int    `json:"month"`
	Limit  uint64 `json:"limit" validate:"required,min=1,max=10"`
	Offset uint64 `json:"offset"`
}

// GetOperations  godoc
// @Summary      Get operations by year by month
// @Description  Get operations by year by month
// @Tags         operations
// @Accept       json
// @Produce      json
// @Success      200
// @Failure      400
// @Failure      500
// @Param        time body handler.OperationInputData true "time"
// @Router       /operations [post]
func (h *Handler) GetOperations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		NewResponse(w, http.StatusBadRequest, "request body read error: "+err.Error())
		return
	}

	var data OperationInputData

	err = json.Unmarshal(body, &data)

	if err != nil {
		NewResponse(w, http.StatusBadRequest, "unable to bind: "+err.Error())
		return
	}

	err = h.validate.Struct(data)

	if err != nil {
		errs := err.(validator.ValidationErrors)
		NewResponse(w, http.StatusBadRequest, errs.Error())
		return
	}

	res, err := h.serv.GetOperations(r.Context(), data.Year, data.Month, data.Limit, data.Offset)

	if err != nil {
		NewResponse(w, http.StatusInternalServerError, err.Error())
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	_, _ = w.Write(res)

	return
}
