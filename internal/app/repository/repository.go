package repository

import (
	"lab1-algorithms/internal/app/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository() (*Repository, error) {
	// ВАЖНО: Если ты менял пароль при установке PostgreSQL, измени password=postgres здесь
	dsn := "host=localhost user=postgres password=postgres dbname=lab2 port=5432 sslmode=disable TimeZone=Europe/Moscow"
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Не удалось подключиться к БД:", err)
	}

	// Автоматическая синхронизация структуры (безопасно, если таблицы уже есть)
	db.AutoMigrate(&models.User{}, &models.AlgorithmBenchmark{}, &models.BenchmarkLike{})
	
	return &Repository{DB: db}, nil
}

// 1. GET: Получение одной услуги для ленты (возвращает сразу 1 строку из БД)
func (r *Repository) GetFeed(id uint) (models.AlgorithmBenchmark, error) {
	var benchmark models.AlgorithmBenchmark
	
	// Если ID передан, ищем его. Иначе берем первую опубликованную
	if id > 0 {
		err := r.DB.Preload("Likes").Where("id = ? AND status = 'published'", id).First(&benchmark).Error
		if err == nil {
			return benchmark, nil
		}
	}
	
	// Если не нашли или ID не передан, берем любую опубликованную
	r.DB.Preload("Likes").Where("status = 'published'").Limit(1).First(&benchmark)
	return benchmark, nil
}

// 2. GET: Получение списка для плитки с фильтрацией через ORM
func (r *Repository) GetGrid(minSize, maxSize int) ([]models.AlgorithmBenchmark, error) {
	var benchmarks []models.AlgorithmBenchmark
	query := r.DB.Preload("Likes").Where("status IN ('published', 'draft')")
	
	if minSize > 0 {
		query = query.Where("input_size >= ?", minSize)
	}
	if maxSize > 0 {
		query = query.Where("input_size <= ?", maxSize)
	}
	
	err := query.Find(&benchmarks).Error
	return benchmarks, err
}

// 3. GET: Получение или создание черновика для текущего пользователя (ID=1)
func (r *Repository) GetOrCreateDraft(userID uint) (models.AlgorithmBenchmark, error) {
	var draft models.AlgorithmBenchmark
	err := r.DB.Where("user_id = ? AND status = 'draft'", userID).First(&draft).Error
	
	if err == gorm.ErrRecordNotFound {
		draft = models.AlgorithmBenchmark{
			Name:       "Новый бенчмарк",
			Status:     "draft",
			ImageURL:   "default.png",
			VideoURL:   "default.mp4",
			InputSize:  100,
			ExecTimeMs: 10,
			UserID:     userID,
		}
		r.DB.Create(&draft)
	}
	return draft, nil
}

// 4. POST: Создание новой услуги через ORM
func (r *Repository) CreateDraft(name, img, vid string, input, exec int, userID uint) error {
	benchmark := models.AlgorithmBenchmark{
		Name:       name,
		Status:     "draft",
		ImageURL:   img,
		VideoURL:   vid,
		InputSize:  input,
		ExecTimeMs: exec,
		UserID:     userID,
	}
	return r.DB.Create(&benchmark).Error
}

// 5. POST: Публикация услуги (смена статуса) через ORM
func (r *Repository) PublishDraft(id uint) error {
	return r.DB.Model(&models.AlgorithmBenchmark{}).Where("id = ?", id).Update("status", "published").Error
}

// 6. POST: Логическое удаление через сырой SQL UPDATE (без ORM, как в задании)
func (r *Repository) DeleteServiceRawSQL(id uint) error {
	query := "UPDATE algorithm_benchmarks SET status = 'deleted', updated_at = CURRENT_TIMESTAMP WHERE id = ?"
	return r.DB.Exec(query, id).Error
}