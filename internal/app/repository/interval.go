package repository

import (
	"Backend-RIP/internal/app/ds"
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IntervalRepository struct {
	db          *gorm.DB
	minioClient *minio.Client
}

func NewIntervalRepository(db *gorm.DB, minioClient *minio.Client) *IntervalRepository {
	return &IntervalRepository{
		db:          db,
		minioClient: minioClient,
	}
}

const (
	intervalImagesBucket = "interval-image"
)

// GET список интервалов с фильтрацией
func (r *IntervalRepository) GetIntervals(title string, toneMin, toneMax float64) ([]ds.Interval, error) {
	var intervals []ds.Interval
	query := r.db.Where("is_delete = ?", false)

	if title != "" {
		query = query.Where("title ILIKE ?", "%"+title+"%")
	}
	if toneMin > 0 {
		query = query.Where("tone >= ?", toneMin)
	}
	if toneMax > 0 {
		query = query.Where("tone <= ?", toneMax)
	}

	err := query.Find(&intervals).Error
	if err != nil {
		return nil, err
	}

	return intervals, nil
}

// GET один интервал
func (r *IntervalRepository) GetInterval(id int) (ds.Interval, error) {
	interval := ds.Interval{}
	err := r.db.Where("id = ? AND is_delete = ?", id, false).First(&interval).Error
	if err != nil {
		return ds.Interval{}, err
	}
	return interval, nil
}

// POST добавление интервала (без изображения)
func (r *IntervalRepository) CreateInterval(interval *ds.Interval) error {
	interval.IsDelete = false
	return r.db.Create(interval).Error
}

// PUT изменение интервала
func (r *IntervalRepository) UpdateInterval(id uint, updates map[string]interface{}) error {
	result := r.db.Model(&ds.Interval{}).Where("id = ? AND is_delete = ?", id, false).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("interval with id %d not found or deleted", id)
	}
	return nil
}

// DELETE удаление интервала
func (r *IntervalRepository) DeleteInterval(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var interval ds.Interval
		if err := tx.Where("id = ?", id).First(&interval).Error; err != nil {
			return err
		}

		if interval.Photo != "" {
			if err := r.deleteIntervalImage(interval.Photo); err != nil {
				return err
			}
		}

		result := tx.Model(&ds.Interval{}).Where("id = ?", id).Update("is_delete", true)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("interval with id %d not found", id)
		}
		return nil
	})
}

// POST добавления интервала в произведение-черновик
func (r *IntervalRepository) AddIntervalToComposition(intervalID uint, creatorID uint, amount uint) error {
	var composition ds.Composition

	err := r.db.Where("creator_id = ? AND status = ?", creatorID, "Черновик").
		First(&composition).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		composition = ds.Composition{
			Status:     "Черновик",
			DateCreate: time.Now(),
			CreatorID:  creatorID,
		}
		if err := r.db.Create(&composition).Error; err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	var existingItem ds.CompositorInterval
	err = r.db.Where("composition_id = ? AND interval_id = ?", composition.ID, intervalID).
		First(&existingItem).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		compositionItem := ds.CompositorInterval{
			CompositionID: composition.ID,
			IntervalID:    intervalID,
			Amount:        amount,
		}
		if err := r.db.Create(&compositionItem).Error; err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {
		existingItem.Amount = amount
		if err := r.db.Save(&existingItem).Error; err != nil {
			return err
		}
	}

	return nil
}

// POST добавление изображения нотной записи интервала
func (r *IntervalRepository) UpdateIntervalPhoto(id uint, fileHeader *multipart.FileHeader) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var interval ds.Interval
		if err := tx.Where("is_delete = false").First(&interval, id).Error; err != nil {
			return err
		}

		if interval.Photo != "" {
			if err := r.deleteIntervalImage(interval.Photo); err != nil {
				return err
			}
		}

		fileExt := filepath.Ext(fileHeader.Filename)
		newFileName := fmt.Sprintf("interval_%d_%d%s", id, time.Now().Unix(), fileExt)
		newFileName = strings.ToLower(newFileName)

		imageURL, err := r.saveIntervalImageToMinIO(newFileName, fileHeader)
		if err != nil {
			return err
		}

		return tx.Model(&interval).Update("photo", imageURL).Error
	})
}

// saveIntervalImageToMinIO сохраняет изображение интервала в MinIO
func (r *IntervalRepository) saveIntervalImageToMinIO(fileName string, fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	fileSize := fileHeader.Size

	contentType := "application/octet-stream"
	if strings.HasSuffix(strings.ToLower(fileName), ".jpg") || strings.HasSuffix(strings.ToLower(fileName), ".jpeg") {
		contentType = "image/jpeg"
	} else if strings.HasSuffix(strings.ToLower(fileName), ".png") {
		contentType = "image/png"
	} else if strings.HasSuffix(strings.ToLower(fileName), ".gif") {
		contentType = "image/gif"
	}

	_, err = r.minioClient.PutObject(context.Background(), intervalImagesBucket, fileName, file, fileSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s:%s/%s/%s", os.Getenv("MINIO_HOST"), os.Getenv("MINIO_SERVER_PORT"), intervalImagesBucket, fileName), nil
}

func (r *IntervalRepository) deleteIntervalImage(imageURL string) error {
	// Если imageURL пустой, ничего не удаляем
	if imageURL == "" {
		return nil
	}

	minioOrigin := os.Getenv("MINIO_HOST") + ":" + os.Getenv("MINIO_SERVER_PORT")

	// Проверяем содержит ли URL правильный origin
	if !strings.Contains(imageURL, minioOrigin) {
		logrus.Printf("Image URL %s doesn't contain MinIO origin, skipping deletion", imageURL)
		return nil
	}

	parts := strings.Split(imageURL, "/")
	if len(parts) == 0 {
		return errors.New("invalid image URL format")
	}

	fileName := parts[len(parts)-1]

	// Проверяем существует ли файл в MinIO
	_, err := r.minioClient.StatObject(context.Background(), intervalImagesBucket, fileName, minio.StatObjectOptions{})
	if err != nil {
		// Файл не существует - это нормально, логируем и продолжаем
		logrus.Printf("File %s not found in MinIO bucket %s, skipping deletion", fileName, intervalImagesBucket)
		return nil
	}

	// Файл существует - удаляем
	err = r.minioClient.RemoveObject(context.Background(), intervalImagesBucket, fileName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete object from MinIO: %v", err)
	}

	logrus.Printf("Successfully deleted interval image from MinIO: %s", fileName)
	return nil
}
