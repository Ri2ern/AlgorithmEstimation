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

// GET 1: Лента
func (h *Handler) GetFeed(ctx *gin.Context) {
	idStr := ctx.Query("id")
	var id uint
	if idStr != "" {
		idUint, _ := strconv.ParseUint(idStr, 10, 32)
		id = uint(idUint)
	}
	
	benchmark, _ := h.Repo.GetFeed(id)
	ctx.HTML(http.StatusOK, "feed.html", gin.H{"Benchmark": benchmark})
}

// GET 2: Плитка с поиском
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

// GET 3: Страница добавления
func (h *Handler) GetAdd(ctx *gin.Context) {
	userID := uint(1) // Хардкод ID пользователя для лабораторной
	draft, _ := h.Repo.GetOrCreateDraft(userID)
	ctx.HTML(http.StatusOK, "add.html", gin.H{"Benchmark": draft})
}

// POST 1: Создание черновика (кнопка "Далее")
func (h *Handler) PostCreateDraft(ctx *gin.Context) {
	name := ctx.PostForm("name")
	img := ctx.PostForm("image_url")
	vid := ctx.PostForm("video_url")
	input, _ := strconv.Atoi(ctx.PostForm("input_size"))
	exec, _ := strconv.Atoi(ctx.PostForm("exec_time_ms"))
	
	h.Repo.CreateDraft(name, img, vid, input, exec, 1)
	ctx.Redirect(http.StatusSeeOther, "/add")
}

// POST 2: Публикация карточки (кнопка "Опубликовать")
func (h *Handler) PostPublish(ctx *gin.Context) {
	idStr := ctx.PostForm("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)
	
	h.Repo.PublishDraft(uint(id))
	ctx.Redirect(http.StatusSeeOther, "/grid")
}

// POST 3: Удаление услуги через сырой SQL (кнопка "Удалить")
func (h *Handler) PostDelete(ctx *gin.Context) {
	idStr := ctx.PostForm("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)
	
	h.Repo.DeleteServiceRawSQL(uint(id))
	ctx.Redirect(http.StatusSeeOther, "/grid")
}