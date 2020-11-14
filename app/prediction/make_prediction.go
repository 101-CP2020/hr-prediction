package prediction

import (
	"fmt"
	"hz.ru/hz/io"
	"time"

	linearmodel "github.com/pa-m/sklearn/linear_model"
	"gonum.org/v1/gonum/mat"
)

type Period = io.Period

func MakePrediction(periods, predictionPeriods []Period) []Period {
	predictor := linearmodel.NewLinearRegression()

	timeAttr := mapPeriods(periods, func(period Period) float64 {
		return timeToFloat(period.Time)
	})

	valueAttr := mapPeriods(periods, func(period Period) float64 {
		return float64(period.Value)
	})

	predictor.Fit(timeAttr, valueAttr)

	fmt.Printf("prediction data :%v\n", valueAttr)

	predictionValuesCount := len(predictionPeriods)

	var predictionTimes []float64
	for _, period := range predictionPeriods {
		predictionTimes = append(predictionTimes, timeToFloat(period.Time))
	}

	timePrediction := mat.NewDense(predictionValuesCount, 1, predictionTimes)
	valuePrediction := mat.NewDense(predictionValuesCount, 1, nil)

	predictor.Predict(timePrediction, valuePrediction)

	fmt.Printf("prediction result :%v\n", valuePrediction)

	for i := 0; i < predictionValuesCount; i++ {
		predictionValue := int(valuePrediction.At(i, 0))
		if predictionValue < 0 {
			predictionValue = 0
		}

		predictionPeriods[i].Value = predictionValue
	}

	return predictionPeriods
}

func mapPeriods(periods []Period, mapper func(period Period) float64) *mat.Dense {
	var data []float64

	for _, period := range periods {
		data = append(data, mapper(period))
	}

	return mat.NewDense(len(periods), 1, data)
}

func timeToFloat(date time.Time) float64 {
	return float64(date.Unix())
}
