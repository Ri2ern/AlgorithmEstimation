package models

import "time"

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"size:50;not null"`
}

type AlgorithmBenchmark struct {
	ID           uint      `gorm:"primaryKey"`
	Name         string    `gorm:"size:100;not null"`
	Description  string    `gorm:"type:text"`
	Status       string    `gorm:"size:20;default:'draft'"`
	ImageURL     string    `gorm:"size:255;column:image_url"`
	VideoURL     string    `gorm:"size:255;column:video_url"`
	InputSize    int       `gorm:"column:input_size"`       // Поле по теме 1
	ExecTimeMs   int       `gorm:"column:exec_time_ms"`     // Поле по теме 2
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
	UserID       uint      `gorm:"column:user_id"`
	User         User      `gorm:"foreignKey:UserID"`
	Likes        []User    `gorm:"many2many:benchmark_likes;joinForeignKey:benchmark_id;joinReferences:user_id"`
}

type BenchmarkLike struct {
	UserID      uint `gorm:"primaryKey;column:user_id"`
	BenchmarkID uint `gorm:"primaryKey;column:benchmark_id"`
}