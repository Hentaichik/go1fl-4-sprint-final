package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	sliceStr := strings.Split(data, ",")
	if len(sliceStr) != 3 {
		return 0, "", 0, errors.New("Введены не все данные")
	}
	stepsCount, err := strconv.Atoi(sliceStr[0])
	if err != nil {
		return 0, "", 0, errors.New("Ошибка: Преобразования количества шагов")
	}
	if stepsCount <= 0 {
		return 0, "", 0, errors.New("Ошибка: Число шагов меньше или равно 0")
	}
	duration, err := time.ParseDuration(sliceStr[2])
	if err != nil {
		return 0, "", 0, errors.New("Ошибка преобразования числа")
	}
	if duration <= 0 {
		return 0, "", 0, errors.New("Ошибка отрицательное время")
	}
	return stepsCount, sliceStr[1], duration, nil
}

func distance(steps int, height float64) float64 {
	len_step := height * stepLengthCoefficient
	count_step := (float64(steps) * len_step) / mInKm
	return count_step
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	dist := distance(steps, height)
	if dist == 0 {
		return 0
	}
	average_speed := dist / duration.Hours()
	return average_speed
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	stepsCount, slice, dur, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}
	str := ""
	switch slice {
	case "Ходьба":
		len_step := distance(stepsCount, height)
		average_speed := meanSpeed(stepsCount, height, dur)
		count_calories, _ := WalkingSpentCalories(stepsCount, weight, height, dur)
		str = fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", slice, dur.Hours(), len_step, average_speed, count_calories)
	case "Бег":
		len_step := distance(stepsCount, height)
		average_speed := meanSpeed(stepsCount, height, dur)
		count_calories, _ := RunningSpentCalories(stepsCount, weight, height, dur)
		str = fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", slice, dur.Hours(), len_step, average_speed, count_calories)
	default:
		return "", errors.New("неизвестный тип тренировки")
	}
	return str, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("Ошибка выполнения функции RunningSpentCalories")
	}
	average_speed := meanSpeed(steps, height, duration)
	minute := duration.Minutes()
	calories_burned := (weight * average_speed * minute) / minInH
	return calories_burned, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("Ошибка выполнения функции WalkingSpentCalories")
	}
	average_speed := meanSpeed(steps, height, duration)
	minute := duration.Minutes()
	count_calories := ((weight * average_speed * minute) / minInH) * walkingCaloriesCoefficient
	return count_calories, nil
}
