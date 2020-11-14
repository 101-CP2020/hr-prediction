package prediction

import (
	"hz.ru/hz/io"

	linearmodel "github.com/pa-m/sklearn/linear_model"
	"gonum.org/v1/gonum/mat"
)

type Period = io.Period

func MakePrediction(periods, predictionPeriods []Period) []Period {
	predictor := linearmodel.NewLinearRegression()

	timeAttr := mapPeriods(periods, func(period Period) float64 {
		return float64(period.Time.Unix())
	})

	valueAttr := mapPeriods(periods, func(period Period) float64 {
		return float64(period.Value)
	})

	predictor.Fit(timeAttr, valueAttr)

	predictionValuesCount := len(predictionPeriods)

	var predictionTimes []float64
	for _, period := range predictionPeriods {
		predictionTimes = append(predictionTimes, float64(period.Time.Unix()))
	}

	timePrediction := mat.NewDense(predictionValuesCount, 1, predictionTimes)
	valuePrediction := mat.NewDense(predictionValuesCount, 1, nil)

	predictor.Predict(timePrediction, valuePrediction)

	for i := 0; i < predictionValuesCount; i++ {
		predictionPeriods[i].Value = int(valuePrediction.At(i, 0))
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
