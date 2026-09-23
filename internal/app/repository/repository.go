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
	dsn := "host=localhost user=postgres password=postgres dbname=lab2 port=5432 sslmode=disable TimeZone=Europe/Moscow"
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Не удалось подключиться к БД:", err)
	}

	db.AutoMigrate(&models.User{}, &models.AlgorithmBenchmark{}, &models.BenchmarkLike{})
	
	return &Repository{DB: db}, nil
}

// Получить одну запись для ленты. Если next=true, ищем следующую после currentID
func (r *Repository) GetFeed(currentID uint, next bool) (models.AlgorithmBenchmark, error) {
	var benchmark models.AlgorithmBenchmark

	// Получаем все опубликованные записи
	var benchmarks []models.AlgorithmBenchmark
	r.DB.Preload("Likes").Where("status = 'published'").Order("id ASC").Find(&benchmarks)

	if len(benchmarks) == 0 {
		return benchmark, nil
	}

	// Если нужно найти следующую запись после текущей
	if next && currentID > 0 {
		for i, b := range benchmarks {
			if b.ID == currentID {
				// Берём следующую по кругу
				nextIdx := (i + 1) % len(benchmarks)
				return benchmarks[nextIdx], nil
			}
		}
	}

	// Если ID передан и запись найдена — возвращаем её
	if currentID > 0 {
		for _, b := range benchmarks {
			if b.ID == currentID {
				return b, nil
			}
		}
	}

	// Иначе возвращаем первую опубликованную
	return benchmarks[0], nil
}

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

func (r *Repository) GetOrCreateDraft(userID uint) (models.AlgorithmBenchmark, error) {
	var draft models.AlgorithmBenchmark
	err := r.DB.Where("user_id = ? AND status = 'draft'", userID).First(&draft).Error
	
	if err == gorm.ErrRecordNotFound {
		draft = models.AlgorithmBenchmark{
			Name:       "Новый бенчмарк",
			Status:     "draft",
			ImageURL:   "/static/default.png",
			VideoURL:   "/static/default.mp4",
			InputSize:  100,
			ExecTimeMs: 10,
			UserID:     userID,
		}
		r.DB.Create(&draft)
	}
	return draft, nil
}

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

func (r *Repository) PublishDraft(id uint) error {
	return r.DB.Model(&models.AlgorithmBenchmark{}).Where("id = ?", id).Update("status", "published").Error
}

func (r *Repository) DeleteServiceRawSQL(id uint) error {
	query := "UPDATE algorithm_benchmarks SET status = 'deleted', updated_at = CURRENT_TIMESTAMP WHERE id = ?"
	return r.DB.Exec(query, id).Error
}