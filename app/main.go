package main

import (
	"fmt"
	"log"
	"time"

	"hz.ru/hz/io"
	"hz.ru/hz/prediction"
	"hz.ru/hz/util"

	"github.com/go-co-op/gocron"
	"github.com/hashicorp/go-multierror"
)

func predictionJob(jobID int) {
	fmt.Printf("prediction job started for %d\n", jobID)

	api := io.NewAPI()

	periods, err := api.GetPeriods(jobID, 30)
	if err != nil {
		log.Println(err)
		return
	}

	now := time.Now()

	predictions := []io.Period{
		{
			Time:  now.Add(util.Quartile),
			Value: 0,
		},
		{
			Time:  now.Add(util.Quartile * 2), //nolint:gomnd // because
			Value: 0,
		},
		{
			Time:  now.Add(util.Quartile * 4), //nolint:gomnd // because
			Value: 0,
		},
	}

	predictionsResult := prediction.MakePrediction(periods, predictions)

	result := io.Prediction{
		JobID:   jobID,
		Time:    now,
		Month3:  predictionsResult[0].Value,
		Month6:  predictionsResult[1].Value,
		Month12: predictionsResult[2].Value,
	}

	err = api.WritePrediction(result)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Printf("prediction job done for %d\n", jobID)
}

func main() {
	scheduler := gocron.NewScheduler(time.UTC)

	queue := createQueue(io.NewAPI())

	var err error

	for _, item := range queue {
		_, sErr := scheduler.Every(1).Month(1).StartAt(item.startTime).Do(predictionJob, item.jobID)
		if sErr != nil {
			err = multierror.Append(err, sErr)
		}
	}

	if err != nil {
		log.Fatal(err)
	}

	scheduler.StartBlocking()
}

type queueItem struct {
	jobID     int
	startTime time.Time
}

func createQueue(api io.API) []queueItem {
	predictions, err := api.GetPredictions()
	if err != nil {
		log.Fatal(err)
	}

	jobs, err := api.GetJobs()
	if err != nil {
		log.Fatal(err)
	}

	index := indexPredictions(predictions)

	var result []queueItem

	for i, job := range jobs {
		lastPredictionTime, found := index[job.ID]

		startTime := time.Now().Add(time.Duration(i))
		if found {
			startTime = lastPredictionTime.Add(util.Month)
		}

		result = append(result, queueItem{
			jobID:     job.ID,
			startTime: startTime,
		})
	}

	return result
}

func indexPredictions(predictions []io.Prediction) map[int]time.Time {
	result := map[int]time.Time{}

	for _, prediction := range predictions {
		result[prediction.JobID] = prediction.Time
	}

	return result
}
