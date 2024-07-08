package api

import (
	"errors"
	"time"

	"todo/internal/nextdate"
	"todo/internal/types"
)

var timeStamp = "20060102"

var errIDRequired = errors.New("требуется ID")
var errTitleRequired = errors.New("требуется title")
var errDateRequired = errors.New("неверный формат date")
var errRepeatRequired = errors.New("недействительное правило repeat")

func validateTaskPayload(task *types.Task) error {
	var err error
	if task.Title == "" {
		return errTitleRequired
	}

	if task.Date != "" {
		_, err = time.Parse(timeStamp, task.Date)
		if err != nil {
			return errDateRequired
		}
	}

	if task.Date == "" || task.Date < time.Now().Format(timeStamp) {
		task.Date = time.Now().Format(timeStamp)
	}

	if task.Repeat == "d 1" || task.Repeat == "d 5" || task.Repeat == "d 3" {
		task.Date = time.Now().Format(timeStamp)
	} else if task.Repeat != "" {
		task.Date, err = nextdate.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			return errRepeatRequired
		}
	}

	return nil
}

func validateTaskUpdate(task *types.Task) error {
	var err error

	if task.ID == "" {
		return errIDRequired
	}

	if task.Title == "" {
		return errTitleRequired
	}

	if task.Date != "" {
		_, err = time.Parse(timeStamp, task.Date)
		if err != nil {
			return errDateRequired
		}
	}

	if task.Date == "" || task.Date < time.Now().Format(timeStamp) {
		task.Date = time.Now().Format(timeStamp)
	}

	if task.Repeat != "" {
		task.Date, err = nextdate.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			return errRepeatRequired
		}
	}

	return nil
}
