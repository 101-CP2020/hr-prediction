package io

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

const apiURL = ""

type Job struct {
	ID   int    `gorm:"column:okpdtr"`
	Name string `gorm:"column:title"`
}

type PeriodInt struct {
	Time  int `gorm:"column:month_year"`
	Value int `gorm:"column:total"`
}

type Period struct {
	Time  time.Time
	Value int
}

type Prediction struct {
	ID      int   `gorm:"column:id"`
	JobID   int   `gorm:"column:okpdtr"`
	Time    int64 `gorm:"column:created_at"`
	Month3  int   `gorm:"column:month_3_value"`
	Month6  int   `gorm:"column:month_6_value"`
	Month12 int   `gorm:"column:month_12_value"`
}

type API interface {
	GetJobs() ([]Job, error)
	GetPeriods(jobID int, periodDays int) ([]Period, error)
	GetPredictions() ([]Prediction, error)
	WritePrediction(prediction Prediction) error
}

func NewAPI() API {
	dsn := "user=db_user password=db_pwd dbname=hr_db host=92.63.103.157 port=7080"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return apiImpl{
		db: db,
	}
}

type apiImpl struct {
	db *gorm.DB
}

func (api apiImpl) GetJobs() ([]Job, error) {
	var result []Job

	queryResult := api.db.Table("tbl_okpdtr").Find(&result)
	if queryResult.Error != nil {
		return nil, queryResult.Error
	}

	return result, nil
}

func (api apiImpl) GetPeriods(jobID, periodDays int) ([]Period, error) {
	var dbResult []PeriodInt

	queryResult := api.db.
		Table("tbl_vacancies").
		Select("month_year, sum(number) as total").
		Group("month_year").
		Where("okpdtr = ?", jobID).
		Find(&dbResult)

	if queryResult.Error != nil {
		return nil, queryResult.Error
	}

	var result []Period

	for _, periodInt := range dbResult {
		year := periodInt.Time % 10000
		month := periodInt.Time / 10000

		date := time.Date(year, time.Month(month), 0, 0, 0, 0, 0, time.UTC)

		result = append(result, Period{
			Time:  date,
			Value: periodInt.Value,
		})
	}

	return result, nil
}

func (api apiImpl) GetPredictions() ([]Prediction, error) {
	var result []Prediction

	queryResult := api.db.Table("tbl_predictions").Find(&result).Limit(10000)
	if queryResult.Error != nil {
		return nil, queryResult.Error
	}

	return result, nil
}

func (api apiImpl) WritePrediction(prediction Prediction) error {
	result := api.db.Table("tbl_predictions").Create(&prediction)

	return result.Error
}
