package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Mansur51-hub/customer-segmentation-service/repository/pgx"
	"github.com/Mansur51-hub/customer-segmentation-service/repository/repoerrs"
	"github.com/go-playground/validator/v10"
	"io"
	"net/http"
	"time"
)

type SegmentInfo struct {
	Slug string `json:"slug" validate:"required,min=1,max=255"`
	Ttl  string `json:"ttl,omitempty"`
}

type UserSegmentsInputData struct {
	UserId           int           `json:"user_id"`
	SegmentsToAdd    []SegmentInfo `json:"segments_to_add" validate:"required,max=10"`
	SegmentsToDelete []string      `json:"segments_to_delete" validate:"required,max=10"`
}

// CreateUserSegments  godoc
// @Summary      Create new user segments
// @Description  Create new user segments
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200
// @Failure      400
// @Failure      404
// @Failure      409
// @Param        segments body handler.UserSegmentsInputData true "segments info"
// @Router       /users [post]
func (h *Handler) CreateUserSegments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		NewResponse(w, http.StatusBadRequest, "request body read error: "+err.Error())
		return
	}

	var data UserSegmentsInputData

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

	slugs, ttl, err := handleSegmentInfo(data.SegmentsToAdd)

	if err != nil {
		NewResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	memberships, err := h.serv.CreateUserSegments(r.Context(), data.UserId, slugs, ttl, data.SegmentsToDelete)

	if errors.Is(err, repoerrs.ErrNotExists) {
		NewResponse(w, http.StatusNotFound, err.Error())
		return
	}

	if errors.Is(err, repoerrs.ErrAlreadyExists) {
		NewResponse(w, http.StatusConflict, err.Error())
		return
	}

	if err != nil {
		NewResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	NewResponse(w, http.StatusOK, memberships)
}

type UserGetSegmentsInputData struct {
	UserId int    `json:"user_id"`
	Limit  uint64 `json:"limit" validate:"required,min=1,max=10"`
	Offset uint64 `json:"offset"`
}

// GetUserSegments  godoc
// @Summary      Get user active segments
// @Description  Get user active segments
// @ID Get user active segments by id
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200
// @Failure      400
// @Param        user body handler.UserGetSegmentsInputData true "user id"
// @Router       /users [get]
func (h *Handler) GetUserSegments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		NewResponse(w, http.StatusBadRequest, "request body read error: "+err.Error())
		return
	}

	var data UserGetSegmentsInputData

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

	memberships, err := h.serv.GetUserSegments(r.Context(), data.UserId, data.Limit, data.Offset)

	if err != nil {
		NewResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if memberships == nil {
		memberships = make([]string, 0, 0)
	}

	NewResponse(w, http.StatusOK, memberships)
}

func handleSegmentInfo(info []SegmentInfo) (slugs []string, ttl []time.Duration, err error) {
	for _, val := range info {
		slugs = append(slugs, val.Slug)
		if val.Ttl == "" {
			d, _ := time.ParseDuration(pgx.DurationNilValue)
			ttl = append(ttl, d)
		} else {
			d, err := time.ParseDuration(val.Ttl)

			if err != nil {
				return nil, nil, fmt.Errorf("error parse ttl: %w", err)
			}

			ttl = append(ttl, d)
		}
	}

	return slugs, ttl, err
}
