package storage

import (
	"errors"
)

type Service interface {
	Save(string) (string, error)
	Load(string) (string, error)
	LoadInfo(string) (*Item, error)
	GetAll() []Item
	Close() error
}

type Item struct {
	Id     string `json:"id" redis:"id"`
	URL    string `json:"url" redis:"url"`
	Visits int    `json:"visits" redis:"visits"`
}

var ErrNoLink = errors.New("url not found")
