package prediction

import (
	"hz.ru/hz/util"

	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMakePrediction(t *testing.T) {
	periods := []Period{
		{
			Time:  time.Now().Add(-time.Hour * 24 * 90),
			Value: 97,
		},
		{
			Time:  time.Now().Add(-time.Hour * 24 * 180),
			Value: 156,
		},
		{
			Time:  time.Now().Add(-time.Hour * 24 * 360),
			Value: 131,
		},
		{
			Time:  time.Now().Add(-time.Hour * 24 * 720),
			Value: 164,
		},
	}

	quartile1 := time.Now().Add(util.Quartile)
	quartile2 := time.Now().Add(util.Quartile * 2)
	quartile4 := time.Now().Add(util.Quartile * 4)

	predictionPeriods := []Period{
		{
			Time: quartile1,
		},
		{
			Time: quartile2,
		},
		{
			Time: quartile4,
		},
	}

	got := MakePrediction(periods, predictionPeriods)

	expected := []Period{
		{
			Time:  quartile1,
			Value: 106,
		},
		{
			Time:  quartile2,
			Value: 99,
		},
		{
			Time:  quartile4,
			Value: 87,
		},
	}

	assert.Equal(t, expected, got)
}
