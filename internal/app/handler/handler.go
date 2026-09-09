package handler

import (
	"lab1-algorithms/internal/app/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

type FeedData struct {
	Estimation repository.AlgorithmEstimation
}

type GridData struct {
	Estimations []repository.AlgorithmEstimation
	FilterValue string
}

func (h *Handler) FeedHandler(ctx *gin.Context) {
	id := ctx.Query("id")
	nextParam := ctx.Query("next")

	var current repository.AlgorithmEstimation

	if id != "" {
		est, err := h.Repository.GetEstimationByID(id)
		if err != nil {
			logrus.Error(err)
		}
		current = est

		if nextParam == "true" {
			estimations, _ := h.Repository.GetEstimations()
			for i, e := range estimations {
				if e.EstimationID == id {
					for j := 1; j <= len(estimations); j++ {
						idx := (i + j) % len(estimations)
						if estimations[idx].PublicationStatus == "published" {
							current = estimations[idx]
							break
						}
					}
					break
				}
			}
		}
	}

	if current.EstimationID == "" {
		estimations, _ := h.Repository.GetEstimations()
		for _, est := range estimations {
			if est.PublicationStatus == "published" {
				current = est
				break
			}
		}
	}

	current.PreviewImageKey = repository.GetMinioURL(current.PreviewImageKey)
	current.DemonstrationVideoKey = repository.GetMinioURL(current.DemonstrationVideoKey)

	data := FeedData{Estimation: current}
	ctx.HTML(http.StatusOK, "feed.html", data)
}

func (h *Handler) AddHandler(ctx *gin.Context) {
	draft, err := h.Repository.GetDraft()
	if err != nil {
		logrus.Error(err)
	}

	draft.PreviewImageKey = repository.GetMinioURL(draft.PreviewImageKey)
	draft.DemonstrationVideoKey = repository.GetMinioURL(draft.DemonstrationVideoKey)

	data := FeedData{Estimation: draft}
	ctx.HTML(http.StatusOK, "add.html", data)
}

func (h *Handler) GridHandler(ctx *gin.Context) {
	filterParam := ctx.Query("filter_volume")

	estimations, err := h.Repository.GetEstimations()
	if err != nil {
		logrus.Error(err)
	}

	var filtered []repository.AlgorithmEstimation
	for _, est := range estimations {
		if est.PublicationStatus == "published" || est.PublicationStatus == "draft" {
			if filterParam != "" {
				filterVol, _ := strconv.Atoi(filterParam)
				if est.InputDataVolume == filterVol {
					est.PreviewImageKey = repository.GetMinioURL(est.PreviewImageKey)
					filtered = append(filtered, est)
				}
			} else {
				est.PreviewImageKey = repository.GetMinioURL(est.PreviewImageKey)
				filtered = append(filtered, est)
			}
		}
	}

	data := GridData{
		Estimations: filtered,
		FilterValue: filterParam,
	}
	ctx.HTML(http.StatusOK, "grid.html", data)
}