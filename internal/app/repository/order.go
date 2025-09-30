package repository

import (
	"Backend-RIP/internal/app/ds"
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (r *Repository) GetIntervals() ([]ds.Interval, error) {
	var intervals []ds.Interval
	err := r.db.Find(&intervals).Error
	if err != nil {
		return nil, err
	}
	if len(intervals) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}

	return intervals, nil
}

func (r *Repository) GetInterval(id int) (ds.Interval, error) {
	interval := ds.Interval{}
	err := r.db.Where("id = ?", id).Find(&interval).Error
	if err != nil {
		return ds.Interval{}, err
	}
	return interval, nil
}

func (r *Repository) GetIntervalsByTitle(title string) ([]ds.Interval, error) {
	var intervals []ds.Interval
	err := r.db.Where("title ILIKE ?", "%"+title+"%").Find(&intervals).Error
	if err != nil {
		return nil, err
	}
	return intervals, nil
}

func (r *Repository) GetCompositionCount() int64 {
	var CompositionID uint
	var count int64
	creatorID := 1

	err := r.db.Model(&ds.Composition{}).Where("creator_id = ? AND status = ?", creatorID, "Черновик").Select("id").First(&CompositionID).Error
	if err != nil {
		return 0
	}

	err = r.db.Model(&ds.CompositorInterval{}).Where("composition_id = ?", CompositionID).Count(&count).Error
	if err != nil {
		logrus.Println("Error counting records in lists_chats:", err)
	}

	return count
}

func (r *Repository) GetActiveCompositionID() uint {
	var CompositionID uint
	err := r.db.Model(&ds.Composition{}).Where("status = ?", "Черновик").Select("id").First(&CompositionID).Error
	if err != nil {
		return 0
	}
	return CompositionID
}

func (r *Repository) GetComposition(id int) ([]ds.CompositorInterval, error) {
	var compositionItems []ds.CompositorInterval
	err := r.db.Where("composition_id = ?", id).Preload("Interval").Find(&compositionItems).Error
	if err != nil {
		return nil, err
	}

	return compositionItems, nil
}

func (r *Repository) AddInterval(intervalsID uint, creatorID uint) error {
	var app ds.Composition

	err := r.db.Where("creator_id = ? AND status = ?", creatorID, "Черновик").
		First(&app).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		app = ds.Composition{
			Status:      "Черновик",
			DateCreate:  time.Now(),
			CreatorID:   creatorID,
			ModeratorID: 2,
		}
		if err := r.db.Create(&app).Error; err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	var count int64
	r.db.Model(&ds.CompositorInterval{}).
		Where("composition_id = ? AND device_id = ?", app.ID, intervalsID).Preload("Interval").
		Count(&count)

	if count == 0 {
		var interval ds.Interval
		if err := r.db.First(&interval, intervalsID).Error; err != nil {
			return err
		}

		appDev := ds.CompositorInterval{
			CompositionID: app.ID,
			IntervalID:    intervalsID,
			Amount:        1,
		}
		if err := r.db.Create(&appDev).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) DeleteComposition(comID uint) error {
	query := `
		UPDATE compositions 
		SET status = 'Удалён', date_update = NOW()
		WHERE id = $1;
	`
	result := r.db.Exec(query, comID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("composition with id %d not found", comID)
	}
	return nil
}

func (r *Repository) IsDraftComposition(comID int) (bool, error) {
	var com ds.Composition
	err := r.db.Select("status").Where("id = ?", comID).First(&com).Error
	if err != nil {
		return false, err
	}
	return com.Status == "Черновик", nil
}
