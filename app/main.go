package main

import (
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/hashicorp/go-multierror"
	"log"
	"time"

	"hz.ru/hz/io"
	"hz.ru/hz/prediction"
	"hz.ru/hz/util"
)

func main() {
	scheduler := gocron.NewScheduler(time.UTC)

	queue := createQueue(api())

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

func predictionJob(jobID int) {
	fmt.Printf("prediction job started for %d\n", jobID)

	api := api()

	periods, err := api.GetPeriods(jobID, 30)
	if err != nil {
		log.Println(err)
		return
	}

	if len(periods) < 2 {
		fmt.Printf("job %d: not enough data for prediction\n", jobID)

		return
	}

	fmt.Printf("job %d: got previous periods, making prediction\n", jobID)

	now := time.Now()

	predictions := []io.Period{
		{
			Time: now.Add(util.Quartile),
		},
		{
			Time: now.Add(util.Quartile * 2), //nolint:gomnd // because
		},
		{
			Time: now.Add(util.Quartile * 4), //nolint:gomnd // because
		},
	}

	predictionsResult := prediction.MakePrediction(periods, predictions)

	result := io.Prediction{
		JobID:   jobID,
		Time:    now.Unix(),
		Month3:  predictionsResult[0].Value,
		Month6:  predictionsResult[1].Value,
		Month12: predictionsResult[2].Value,
	}

	fmt.Printf("job %d: writing prediction result\n", jobID)

	err = api.WritePrediction(result)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Printf("job %d: prediction done\n", jobID)
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

		startTime := time.Now().Add(time.Duration(i) * time.Second)
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

var apiInstance io.API //nolint:gochecknoglobals // because

func api() io.API {
	if apiInstance == nil {
		apiInstance = io.NewAPI()
	}

	return apiInstance
}

func indexPredictions(predictions []io.Prediction) map[int]time.Time {
	result := map[int]time.Time{}

	for _, prediction := range predictions {
		result[prediction.JobID] = time.Unix(prediction.Time, 0)
	}

	return result
}
