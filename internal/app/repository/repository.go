package repository

import (
	"fmt"
	"strings"
)

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

func (r *Repository) GetIntervals() ([]Interval, error) {
	intervals := []Interval{
		{
			ID:          1,
			Title:       "Большая терция",
			Tone:        2,
			Image:       "bigtercia.png",
			Description: "Интервал в три ступени",
		},
		{
			ID:          2,
			Title:       "Чистая кварта",
			Tone:        2.5,
			Image:       "purekvarna.png",
			Description: "Интервал в четыре ступени",
		},
		{
			ID:          3,
			Title:       "Большая секунда",
			Tone:        1,
			Image:       "bigecunda.png",
			Description: "Интервал в две ступени",
		},
		{
			ID:          4,
			Title:       "Чистая квинта",
			Tone:        3.5,
			Image:       "purekvinta.png",
			Description: "Интервал в пять ступеней",
		},
		{
			ID:          5,
			Title:       "Чистая октава",
			Tone:        6,
			Image:       "pureoctava.png",
			Description: "Интервал в восемь ступеней",
		},
		{
			ID:          6,
			Title:       "Малая секста",
			Tone:        4,
			Image:       "smallsexta.png",
			Description: "Интервал в шесть ступеней",
		},
		{
			ID:          7,
			Title:       "Большая септима",
			Tone:        5.5,
			Image:       "bigseptima.png",
			Description: "Интервал в семь ступеней",
		},
		{
			ID:          8,
			Title:       "Малая секунда",
			Tone:        0.5,
			Image:       "smallsecunda.png",
			Description: "Интервал в две ступени",
		},
	}
	if len(intervals) == 0 {
		return nil, fmt.Errorf("Массив пустой")
	}

	return intervals, nil
}

func (r *Repository) GetInterval(id int) (Interval, error) {
	intervals, err := r.GetIntervals()
	if err != nil {
		return Interval{}, err // тут у нас уже есть кастомная ошибка из нашего метода, поэтому мы можем просто вернуть ее
	}

	for _, interval := range intervals {
		if interval.ID == id {
			return interval, nil // если нашли, то просто возвращаем найденный заказ (услугу) без ошибок
		}
	}
	return Interval{}, fmt.Errorf("Заказ не найден")
}

func (r *Repository) GetIntervalsByTitle(title string) ([]Interval, error) {
	intervals, err := r.GetIntervals()
	if err != nil {
		return []Interval{}, err
	}

	var result []Interval
	for _, interval := range intervals {
		if strings.Contains(strings.ToLower(interval.Title), strings.ToLower(title)) {
			result = append(result, interval)
		}
	}

	return result, nil
}

func (r *Repository) GetCart() ([]Interval, error) {
	intervals, err := r.GetIntervals()
	if err != nil {
		return []Interval{}, err
	}

	var result []Interval
	for _, interval := range intervals {
		if interval.ID == 4 || interval.ID == 7 {
			result = append(result, interval)
		}
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("Массив пустой")
	}

	return result, nil
}
