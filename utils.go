package main

import (
	"encoding/json"
	"errors"
	uuid "github.com/nu7hatch/gouuid"
	"io/ioutil"
)

type EventData struct {
	ID      string
	Name    string
	Price   int64
	Members int
}

var (
	events = []EventData{
		{
			ID:      "1",
			Name:    "Paper Representation",
			Price:   300,
			Members: 2,
		},
		{
			ID:      "2",
			Name:    "Poster Representation",
			Price:   200,
			Members: 2,
		},
		{
			ID:      "3",
			Name:    "Techinal Quiz",
			Price:   500,
			Members: 4,
		},
		{
			ID:      "4",
			Name:    "Mobile E-Sports",
			Price:   100,
			Members: 4,
		},
		{
			ID:      "5",
			Name:    "Blind Coding",
			Price:   50,
			Members: 1,
		},
		{
			ID:      "6",
			Name:    "Cultural Events",
			Price:   100,
			Members: 1,
		},
	}
)

func GetEventDetails(id string) *EventData {
	for _, event := range events {
		if event.ID == id {
			return &event
		}
	}

	return nil
}

func getUUID() string {
	temp, _ := uuid.NewV4()
	return temp.String()
}

func getConfig() (error, *Config) {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		return errors.New(err.Error()), nil
	}
	config := &Config{}
	err = json.Unmarshal(data, config)
	if err != nil {
		return errors.New(err.Error()), nil
	}
	return nil, config
}
