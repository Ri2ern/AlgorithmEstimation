package handler

import (
	"lab1-algorithms/internal/app/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Repo *repository.Repository
}

func NewHandler(repo *repository.Repository) *Handler {
	return &Handler{Repo: repo}
}

func (h *Handler) GetFeed(ctx *gin.Context) {
	idStr := ctx.Query("id")
	nextParam := ctx.Query("next")
	
	var id uint
	if idStr != "" {
		idUint, _ := strconv.ParseUint(idStr, 10, 32)
		id = uint(idUint)
	}
	
	// Если передан параметр next=true, ищем следующую опубликованную запись
	if nextParam == "true" && id > 0 {
		benchmark, _ := h.Repo.GetFeed(id)
		if benchmark.ID > 0 {
			// Получаем следующую опубликованную запись после текущей
			nextBenchmark, _ := h.Repo.GetNextPublished(benchmark.ID)
			if nextBenchmark.ID > 0 {
				ctx.HTML(http.StatusOK, "feed.html", gin.H{"Benchmark": nextBenchmark})
				return
			}
		}
	}
	
	// Иначе просто получаем запись по ID или первую опубликованную
	benchmark, _ := h.Repo.GetFeed(id)
	ctx.HTML(http.StatusOK, "feed.html", gin.H{"Benchmark": benchmark})
}

func (h *Handler) GetGrid(ctx *gin.Context) {
	minStr := ctx.Query("filter_min")
	maxStr := ctx.Query("filter_max")
	
	min, _ := strconv.Atoi(minStr)
	max, _ := strconv.Atoi(maxStr)
	
	benchmarks, _ := h.Repo.GetGrid(min, max)
	ctx.HTML(http.StatusOK, "grid.html", gin.H{
		"Benchmarks": benchmarks,
		"FilterMin":  minStr,
		"FilterMax":  maxStr,
	})
}

func (h *Handler) GetAdd(ctx *gin.Context) {
	userID := uint(1)
	draft, _ := h.Repo.GetOrCreateDraft(userID)
	ctx.HTML(http.StatusOK, "add.html", gin.H{"Benchmark": draft})
}

func (h *Handler) PostCreateDraft(ctx *gin.Context) {
	name := ctx.PostForm("name")
	img := ctx.PostForm("image_url")
	vid := ctx.PostForm("video_url")
	input, _ := strconv.Atoi(ctx.PostForm("input_size"))
	exec, _ := strconv.Atoi(ctx.PostForm("exec_time_ms"))
	
	h.Repo.CreateDraft(name, img, vid, input, exec, 1)
	ctx.Redirect(http.StatusSeeOther, "/add")
}

func (h *Handler) PostPublish(ctx *gin.Context) {
	idStr := ctx.PostForm("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)
	
	h.Repo.PublishDraft(uint(id))
	ctx.Redirect(http.StatusSeeOther, "/grid")
}

func (h *Handler) PostDelete(ctx *gin.Context) {
	idStr := ctx.PostForm("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)
	
	h.Repo.DeleteServiceRawSQL(uint(id))
	ctx.Redirect(http.StatusSeeOther, "/grid")
}