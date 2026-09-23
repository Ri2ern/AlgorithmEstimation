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
	return &Handler{Repository: r}
}

type FeedData struct {
	Benchmark repository.AlgoBenchmark
}

type GridData struct {
	Benchmarks []repository.AlgoBenchmark
	FilterMin  string
	FilterMax  string
}

func (h *Handler) FeedHandler(ctx *gin.Context) {
	id := ctx.Query("id")
	nextParam := ctx.Query("next")
	var current repository.AlgoBenchmark

	if id != "" {
		b, err := h.Repository.GetBenchmarkByID(id)
		if err == nil {
			current = b
			if nextParam == "true" {
				benchmarks, _ := h.Repository.GetBenchmarks()
				for i, b := range benchmarks {
					if b.BenchmarkID == id {
						for j := 1; j <= len(benchmarks); j++ {
							idx := (i + j) % len(benchmarks)
							if benchmarks[idx].Status == "published" {
								current = benchmarks[idx]
								break
							}
						}
						break
					}
				}
			}
		}
	}

	if current.BenchmarkID == "" {
		benchmarks, _ := h.Repository.GetBenchmarks()
		for _, b := range benchmarks {
			if b.Status == "published" {
				current = b
				break
			}
		}
	}

	current.PreviewImageKey = repository.GetMinioURL(current.PreviewImageKey)
	current.VideoKey = repository.GetMinioURL(current.VideoKey)

	ctx.HTML(http.StatusOK, "feed.html", FeedData{Benchmark: current})
}

func (h *Handler) AddHandler(ctx *gin.Context) {
	draft, err := h.Repository.GetDraft()
	if err != nil {
		logrus.Error(err)
	}
	draft.PreviewImageKey = repository.GetMinioURL(draft.PreviewImageKey)
	draft.VideoKey = repository.GetMinioURL(draft.VideoKey)
	ctx.HTML(http.StatusOK, "add.html", FeedData{Benchmark: draft})
}

func (h *Handler) GridHandler(ctx *gin.Context) {
	minStr := ctx.Query("filter_min")
	maxStr := ctx.Query("filter_max")

	benchmarks, err := h.Repository.GetBenchmarks()
	if err != nil {
		logrus.Error(err)
	}

	var filtered []repository.AlgoBenchmark
	for _, b := range benchmarks {
		if b.Status == "published" || b.Status == "draft" {
			include := true
			if minStr != "" {
				minVal, _ := strconv.Atoi(minStr)
				if b.InputSize < minVal {
					include = false
				}
			}
			if maxStr != "" {
				maxVal, _ := strconv.Atoi(maxStr)
				if b.InputSize > maxVal {
					include = false
				}
			}

			if include {
				b.PreviewImageKey = repository.GetMinioURL(b.PreviewImageKey)
				filtered = append(filtered, b)
			}
		}
	}

	ctx.HTML(http.StatusOK, "grid.html", GridData{
		Benchmarks: filtered,
		FilterMin:  minStr,
		FilterMax:  maxStr,
	})
}