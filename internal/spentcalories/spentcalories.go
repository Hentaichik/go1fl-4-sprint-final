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
		return 0, "", 0, errors.New("введены не все данные")
	}
	stepsCount, err := strconv.Atoi(sliceStr[0])
	if err != nil {
		return 0, "", 0, errors.New("преобразования количества шагов")
	}
	if stepsCount <= 0 {
		return 0, "", 0, errors.New("число шагов меньше или равно 0")
	}
	duration, err := time.ParseDuration(sliceStr[2])
	if err != nil {
		return 0, "", 0, errors.New("преобразования числа")
	}
	if duration <= 0 {
		return 0, "", 0, errors.New("отрицательное время")
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
	return distance(steps, height) / duration.Hours()
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	str := ""
	var countCalories float64
	stepsCount, slice, dur, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}
	switch slice {
	case "Ходьба":
		countCalories, err = WalkingSpentCalories(stepsCount, weight, height, dur)
		if err != nil {
			return "", err
		}
	case "Бег":
		countCalories, err = RunningSpentCalories(stepsCount, weight, height, dur)
		if err != nil {
			return "", err
		}
	default:
		return "", errors.New("неизвестный тип тренировки")
	}
	len_step := distance(stepsCount, height)
	averageSpeed := meanSpeed(stepsCount, height, dur)
	str = fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", slice, dur.Hours(), len_step, averageSpeed, countCalories)
	return str, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("неправильные шаги")
	}else if weight <= 0{
		return 0, errors.New("неправильный вес")
	}else if height <= 0{
		return 0, errors.New("неправильный рост")
	}else if  duration <= 0{
		return 0, errors.New("неправильная продолжительность")
	}
	averageSpeed := meanSpeed(steps, height, duration)
	minute := duration.Minutes()
	caloriesBurned := (weight * averageSpeed * minute) / minInH
	return caloriesBurned, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("неправильные шаги")
	}else if weight <= 0{
		return 0, errors.New("неправильный вес")
	}else if height <= 0{
		return 0, errors.New("неправильный рост")
	}else if  duration <= 0{
		return 0, errors.New("неправильная продолжительность")
	}
	averageSpeed := meanSpeed(steps, height, duration)
	minute := duration.Minutes()
	countCalories := ((weight * averageSpeed * minute) / minInH) * walkingCaloriesCoefficient
	return countCalories, nil
}
