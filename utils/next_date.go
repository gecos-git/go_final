package utils

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	startDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", err
	}

	if repeat == "" {
		return "", errors.New("repeat пустой")
	}

	repeatParts := strings.Split(repeat, " ")

	switch repeatParts[0] {
	case "d":
		if len(repeatParts) == 1 {
			return "", errors.New("формат repeat неверный")
		}

		days, err := strconv.Atoi(repeatParts[1])
		if err != nil || days > 400 {
			return "", errors.New("неверное количество дней")
		}

		for {
			startDate = startDate.AddDate(0, 0, days)

			if !startDate.Before(now) && !startDate.Equal(now) {
				break
			}
		}

	case "y":
		startDate = startDate.AddDate(1, 0, 0)

		for now.After(startDate) || now.Equal(startDate) {
			startDate = startDate.AddDate(1, 0, 0)
		}

	default:
		return "", errors.New("не поддерживаемый формат поля repeat")
	}

	return startDate.Format("20060102"), nil
}
