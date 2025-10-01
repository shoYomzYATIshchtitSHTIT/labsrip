package repository

import (
	"fmt"
	"strings"
)

// =====================
// Основные интервалы
// =====================

type Repository struct {
}

func NewRepository() (*Repository, error) {
	return &Repository{}, nil
}

type Interval struct {
	ID          int
	Title       string
	Tone        float64
	Image       string
	Description string
}

var intervals = []Interval{
	{ID: 1, Title: "Большая терция", Tone: 2, Image: "bigtercia.png", Description: "Интервал в три ступени"},
	{ID: 2, Title: "Чистая кварта", Tone: 2.5, Image: "purekvarna.png", Description: "Интервал в четыре ступени"},
	{ID: 3, Title: "Большая секунда", Tone: 1, Image: "bigecunda.png", Description: "Интервал в две ступени"},
	{ID: 4, Title: "Чистая квинта", Tone: 3.5, Image: "purekvinta.png", Description: "Интервал в пять ступеней"},
	{ID: 5, Title: "Чистая октава", Tone: 6, Image: "pureoctava.png", Description: "Интервал в восемь ступеней"},
	{ID: 6, Title: "Малая секста", Tone: 4, Image: "smallsexta.png", Description: "Интервал в шесть ступеней"},
	{ID: 7, Title: "Большая септима", Tone: 5.5, Image: "bigseptima.png", Description: "Интервал в семь ступеней"},
	{ID: 8, Title: "Малая секунда", Tone: 0.5, Image: "smallsecunda.png", Description: "Интервал в две ступени"},
}

func (r *Repository) GetIntervals() ([]Interval, error) {
	if len(intervals) == 0 {
		return nil, fmt.Errorf("Массив пустой")
	}
	return intervals, nil
}

func (r *Repository) GetInterval(id int) (Interval, error) {
	for _, interval := range intervals {
		if interval.ID == id {
			return interval, nil
		}
	}
	return Interval{}, fmt.Errorf("Интервал не найден")
}

func (r *Repository) GetIntervalsByTitle(title string) ([]Interval, error) {
	var result []Interval
	for _, interval := range intervals {
		if strings.Contains(strings.ToLower(interval.Title), strings.ToLower(title)) {
			result = append(result, interval)
		}
	}
	return result, nil
}

func (r *Repository) GetComposition() ([]Interval, error) {
	var result []Interval
	for _, interval := range intervals {
		if interval.ID == 1 || interval.ID == 2 { // пример для композиции
			result = append(result, interval)
		}
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("Массив пустой")
	}
	return result, nil
}

// =====================
// Заявки на композицию
// =====================

type CompositionRequestRepository struct {
}

func NewCompositionRequestRepository() (*CompositionRequestRepository, error) {
	return &CompositionRequestRepository{}, nil
}

type CompositionRequest struct {
	ID        int
	Belonging string
}

// запись в заявке: интервал + количество
type CompositionRequestViewEntry struct {
	Interval Interval
	Amount   int
}

// заявка с набором интервалов
type CompositionRequestView struct {
	CompositionRequest CompositionRequest
	Entries            []CompositionRequestViewEntry
}

// тестовые данные
var compositionRequestViewByID = map[int]CompositionRequestView{
	1: {
		CompositionRequest: CompositionRequest{
			ID:        1,
			Belonging: "является классическим",
		},
		Entries: []CompositionRequestViewEntry{
			{Interval: intervals[0], Amount: 2}, // ID=1
			{Interval: intervals[1], Amount: 1}, // ID=2
		},
	},
}

// методы

func (*CompositionRequestRepository) GetCompositionEntryCntByID(id int) (int, error) {
	view, found := compositionRequestViewByID[id]
	if !found {
		return 0, fmt.Errorf("не найдено")
	}
	return len(view.Entries), nil
}

func (*CompositionRequestRepository) GetCompositionRequestViewByID(id int) (*CompositionRequestView, error) {
	view, found := compositionRequestViewByID[id]
	if !found {
		return nil, fmt.Errorf("не найдено")
	}
	return &view, nil
}
