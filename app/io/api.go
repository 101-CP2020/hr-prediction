package io

import (
	"time"

	"github.com/imroc/req"
)

const apiURL = ""

type Job struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Period struct {
	Time  time.Time `json:"time"`
	Value int       `json:"value"`
}

type Prediction struct {
	JobID   int       `json:"job_id"`
	Time    time.Time `json:"time"`
	Month3  int       `json:"month_3"`
	Month6  int       `json:"month_6"`
	Month12 int       `json:"month_12"`
}

type API interface {
	GetJobs() ([]Job, error)
	GetPeriods(jobID int, periodDays int) ([]Period, error)
	GetPredictions() ([]Prediction, error)
	WritePrediction(prediction Prediction) error
}

func NewAPI() API {
	return apiImpl{}
}

type apiImpl struct{}

func (apiImpl) GetJobs() ([]Job, error) {
	response, err := req.Get(apiURL + "jobs")
	if err != nil {
		return nil, err
	}

	var result []Job

	err = response.ToJSON(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (apiImpl) GetPeriods(jobID, periodDays int) ([]Period, error) {
	params := req.Param{
		"id":     jobID,
		"period": periodDays,
	}

	response, err := req.Get(apiURL+"periods", params)
	if err != nil {
		return nil, err
	}

	var result []Period

	err = response.ToJSON(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (apiImpl) GetPredictions() ([]Prediction, error) {
	response, err := req.Get(apiURL + "predictions")
	if err != nil {
		return nil, err
	}

	var result []Prediction

	err = response.ToJSON(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (apiImpl) WritePrediction(prediction Prediction) error {
	_, err := req.Post(apiURL+"predictions", req.BodyJSON(&prediction))
	if err != nil {
		return err
	}

	return nil
}
