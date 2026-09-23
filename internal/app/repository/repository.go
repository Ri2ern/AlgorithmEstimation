package repository

import (
	"fmt"
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

func (r *Repository) GetFeed(id uint) (models.AlgorithmBenchmark, error) {
	var benchmark models.AlgorithmBenchmark
	
	if id > 0 {
		err := r.DB.Preload("Likes").Where("id = ? AND status = 'published'", id).First(&benchmark).Error
		if err == nil {
			return benchmark, nil
		}
	}
	
	r.DB.Preload("Likes").Where("status = 'published'").Limit(1).First(&benchmark)
	return benchmark, nil
}

func (r *Repository) GetNextPublished(currentID uint) (models.AlgorithmBenchmark, error) {
	var benchmark models.AlgorithmBenchmark
	
	err := r.DB.Preload("Likes").
		Where("status = 'published' AND id > ?", currentID).
		Order("id ASC").
		Limit(1).
		First(&benchmark).Error
	
	if err != nil {
		r.DB.Preload("Likes").
			Where("status = 'published'").
			Order("id ASC").
			Limit(1).
			First(&benchmark)
	}
	
	return benchmark, nil
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

// SQL КУРСОР 
func (r *Repository) DeleteServiceRawSQL(id uint) error {
	query := fmt.Sprintf(`
		DO $$
		DECLARE
			cur CURSOR FOR 
				SELECT id FROM algorithm_benchmarks 
				WHERE id = %d 
				FOR UPDATE;
			rec RECORD;
		BEGIN
			OPEN cur;
			FETCH cur INTO rec;
			IF FOUND THEN
				UPDATE algorithm_benchmarks 
				SET status = 'deleted', updated_at = CURRENT_TIMESTAMP 
				WHERE CURRENT OF cur;
			END IF;
			CLOSE cur;
		END $$;
	`, id)
	
	return r.DB.Exec(query).Error
}