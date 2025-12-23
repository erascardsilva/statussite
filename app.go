// Erasmo Cardoso - Dev
package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

type SiteStatus struct {
	URL      string `json:"url"`
	Status   string `json:"status"`
	Message  string `json:"message"`
	IsOnline bool   `json:"isOnline"`
}

// CheckSites - verifica todos os sites em paralelo
func (a *App) CheckSites() []SiteStatus {
	// lista dos sites que vou monitorar
	sites := []string{
		"https://www.electrocode.com.br",
		"https://www.vmi-informatica.com.br",
		"https://www.electrocode.com.br/app/authentication/sign-in/",
	}

	results := make([]SiteStatus, len(sites))
	var wg sync.WaitGroup
	wg.Add(len(sites))

	// checa cada site em uma goroutine separada
	for i, site := range sites {
		go func(idx int, siteUrl string) {
			defer wg.Done()
			results[idx] = checkSite(siteUrl)
		}(i, site)
	}

	wg.Wait()
	return results
}

func checkSite(url string) SiteStatus {
	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Head(url)

	if err != nil {
		// falhou a conexão
		return SiteStatus{
			URL:      url,
			Status:   "Erro",
			Message:  fmt.Sprintf("Erro ao verificar: %v", err),
			IsOnline: false,
		}
	}
	defer res.Body.Close()

	online := res.StatusCode == http.StatusOK
	msg := "Online"
	if !online {
		msg = "Offline"
	}

	return SiteStatus{
		URL:      url,
		Status:   res.Status,
		Message:  msg,
		IsOnline: online,
	}
}
