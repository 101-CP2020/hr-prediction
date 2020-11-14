package main

import (
	"hz.ru/hz/io"
	"hz.ru/hz/util"

	"math"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type apiMock struct{}

func (apiMock) GetJobs() ([]io.Job, error) {
	var result []io.Job

	for i := 0; i < 5; i++ {
		result = append(result, io.Job{
			ID:   i,
			Name: strconv.Itoa(i),
		})
	}

	return result, nil
}

func (apiMock) GetPeriods(jobID, periodDays int) ([]io.Period, error) {
	var result []io.Period

	for i := 0; i < 10; i++ {
		result = append(result, io.Period{
			Time:  time.Now().Add(util.Quartile * time.Duration(i)),
			Value: i * 131,
		})
	}

	return result, nil
}

func (apiMock) WritePrediction(prediction io.Prediction) error {
	return nil
}

func (apiMock) GetPredictions() ([]io.Prediction, error) {
	var result []io.Prediction

	for i := 0; i < 5; i++ {
		result = append(result, io.Prediction{
			JobID:   int(math.Mod(float64(i), 5)),
			Time:    time.Now().Add(util.Quartile * time.Duration(i)),
			Month3:  100,
			Month6:  100,
			Month12: 100,
		})
	}

	return result, nil
}

func Test_createQueue(t *testing.T) {
	queue := createQueue(apiMock{})

	expected := []queueItem{
		{
			jobID:     0,
			startTime: time.Now(),
		},
		{
			jobID:     1,
			startTime: time.Now(),
		},
		{
			jobID:     2,
			startTime: time.Now(),
		},
		{
			jobID:     3,
			startTime: time.Now(),
		},
		{
			jobID:     4,
			startTime: time.Now(),
		},
	}

	assert.Equal(t, expected, queue)
}
