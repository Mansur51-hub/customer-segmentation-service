package handler

import (
	"encoding/json"
	"errors"
	"github.com/Mansur51-hub/customer-segmentation-service/repository/repoerrs"
	"github.com/go-playground/validator/v10"
	"io"
	"net/http"
)

type SegmentInputData struct {
	Slug    string `json:"slug" validate:"required,min=1,max=255"`
	Percent int    `json:"percent,omitempty" validate:"min=1,max=100"`
}

// CreateSegment             godoc
// @Summary      Create new segment
// @Description  Create new segment
// @Tags         segments
// @Accept       json
// @Produce      json
// @Success      201 {object} model.Segment "segment"
// @Failure      400
// @Failure      409
// @Failure      500
// @Param        slug body handler.SegmentInputData true "Segment slug"
// @Router       /segments [post]
func (h *Handler) CreateSegment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		NewResponse(w, http.StatusBadRequest, "request body read error: "+err.Error())
		return
	}

	var segment SegmentInputData

	err = json.Unmarshal(body, &segment)

	if err != nil {
		NewResponse(w, http.StatusBadRequest, "unable to bind: "+err.Error())
		return
	}

	err = h.validate.Struct(segment)

	if err != nil {
		errs := err.(validator.ValidationErrors)
		NewResponse(w, http.StatusBadRequest, errs.Error())
		return
	}

	seg, err := h.serv.CreateSegment(r.Context(), segment.Slug, segment.Percent)

	if errors.Is(err, repoerrs.ErrAlreadyExists) {
		NewResponse(w, http.StatusConflict, repoerrs.ErrAlreadyExists.Error())
		return
	}

	if err != nil {
		NewResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	NewResponse(w, http.StatusCreated, seg)
}

type SegmentDeleteData struct {
	Slug string `json:"slug" validate:"required,min=1,max=255"`
}

// DeleteSegment             godoc
// @Summary      Delete segment
// @Description  Delete segment
// @Tags         segments
// @Accept       json
// @Produce      json
// @Success      200
// @Failure      400
// @Failure      404
// @Failure      500
// @Param        slug body handler.SegmentDeleteData true "Segment slug"
// @Router       /segments [delete]
func (h *Handler) DeleteSegment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		NewResponse(w, http.StatusBadRequest, "request body read error: "+err.Error())
		return
	}

	var segment SegmentDeleteData

	err = json.Unmarshal(body, &segment)

	if err != nil {
		NewResponse(w, http.StatusBadRequest, "unable to bind: "+err.Error())
		return
	}

	err = h.validate.Struct(segment)

	if err != nil {
		errs := err.(validator.ValidationErrors)
		NewResponse(w, http.StatusBadRequest, errs.Error())
		return
	}

	err = h.serv.DeleteSegment(r.Context(), segment.Slug)

	if errors.Is(err, repoerrs.ErrNotExists) {
		NewResponse(w, http.StatusNotFound, repoerrs.ErrNotExists.Error())
		return
	}

	NewResponse(w, http.StatusCreated, "success")
}
