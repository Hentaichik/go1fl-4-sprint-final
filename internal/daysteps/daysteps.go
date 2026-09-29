package daysteps

import (
	"errors"
	"fmt"

	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {

	str := strings.Split(data, ",")

	if len(str) != 2 {
		return 0, 0, errors.New("неправильный ввод данных")
	}
	stepsCount, err := strconv.Atoi(str[0])
	if err != nil {
		return 0, 0, errors.New("неправильное преобразования числа")
	}
	if stepsCount <= 0 {
		return 0, 0, errors.New("количество шагов должно быть больше 0")
	}
	duration, err := time.ParseDuration(str[1])
	if err != nil {
		return 0, 0, errors.New("неправильное преобразования числа")
	}
	if duration <= 0 {
		return 0, 0, errors.New("время меньше 0")
	}
	return stepsCount, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	stepsCount, timeSpent, err := parsePackage(data)
	if err != nil {
		
		return ""
	}
	dist := float64(stepsCount) * stepLength
	dist /= mInKm
	caloriesSpent, err := spentcalories.WalkingSpentCalories(stepsCount, weight, height, timeSpent)
	if err != nil {
		return ""
	}
	str := fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", stepsCount, dist, caloriesSpent)
	return str
}
