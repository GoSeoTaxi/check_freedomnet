package repository

import (
	"errors"
	"sync"

	"github.com/go-resty/resty/v2"
)

type FreedomNetRepo struct {
	Servers    []string
	Client     *resty.Client
	currentIdx int
	mu         sync.Mutex
}

func NewFreedomNetRepo(servers []string) *FreedomNetRepo {
	client := resty.New()
	return &FreedomNetRepo{
		Servers:    servers,
		Client:     client,
		currentIdx: 0,
	}
}

func (r *FreedomNetRepo) FetchFromServers() (string, error) {

	numServers := len(r.Servers)

	r.mu.Lock()
	server := r.Servers[r.currentIdx]
	r.currentIdx = (r.currentIdx + 1) % numServers
	r.mu.Unlock()

	resp, err := r.Client.R().Get(server)
	if err == nil && resp.IsSuccess() {
		return resp.String(), nil
	}

	return "", errors.New("no server responded successfully")
}
