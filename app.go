// Erasmo Cardoso - Dev
package main

import (
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"strings"
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

type APIResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

// createHTTPClient - cria cliente HTTP otimizado para evitar connection reset
func createHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives: true, // evita reutilização de conexões
			MaxIdleConns:      1,
			IdleConnTimeout:   1 * time.Second,
			DialContext: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 0,
			}).DialContext,
		},
	}
}

// formatError - formata mensagens de erro de forma amigável
func formatError(err error) string {
	errMsg := err.Error()

	if strings.Contains(errMsg, "connection reset") {
		return "Conexão interrompida pelo servidor"
	}
	if strings.Contains(errMsg, "timeout") {
		return "Tempo de resposta excedido"
	}
	if strings.Contains(errMsg, "no such host") {
		return "Site não encontrado"
	}
	if strings.Contains(errMsg, "connection refused") {
		return "Conexão recusada"
	}

	return "Erro de conexão"
}

// CheckSites - verifica todos os sites em paralelo
func (a *App) CheckSites() []SiteStatus {
	// lista dos sites que vou monitorar
	sites := []string{
		"https://www.electrocode.com.br",
		"https://www.vmi-informatica.com.br",
		"https://www.electrocode.com.br/app/authentication/sign-in/",
		"https://www.electrocode.com.br/api/",
	}

	results := make([]SiteStatus, len(sites))
	var wg sync.WaitGroup
	wg.Add(len(sites))

	// checa cada site em uma goroutine separada
	for i, site := range sites {
		go func(idx int, siteUrl string) {
			defer wg.Done()
			// verifica se é um endpoint de API
			if strings.Contains(siteUrl, "/api/") {
				results[idx] = checkAPIEndpoint(siteUrl)
			} else {
				results[idx] = checkSite(siteUrl)
			}
		}(i, site)
	}

	wg.Wait()
	return results
}

func checkSite(url string) SiteStatus {
	maxRetries := 3
	var lastErr error
	var res *http.Response

	// retry logic com backoff exponencial
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			// backoff: 500ms, 1s, 2s
			backoff := time.Duration(500*(1<<uint(attempt-1))) * time.Millisecond
			time.Sleep(backoff)
		}

		client := createHTTPClient()
		var err error
		res, err = client.Head(url)

		if err == nil {
			// sucesso!
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

		lastErr = err

		// se não for erro de rede temporário, não tenta novamente
		if !strings.Contains(err.Error(), "connection reset") &&
			!strings.Contains(err.Error(), "timeout") {
			break
		}
	}

	// todas as tentativas falharam
	return SiteStatus{
		URL:      url,
		Status:   "Erro",
		Message:  formatError(lastErr),
		IsOnline: false,
	}
}

func checkAPIEndpoint(url string) SiteStatus {
	maxRetries := 3
	var lastErr error

	// retry logic com backoff exponencial
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			// backoff: 500ms, 1s, 2s
			backoff := time.Duration(500*(1<<uint(attempt-1))) * time.Millisecond
			time.Sleep(backoff)
		}

		client := createHTTPClient()
		res, err := client.Get(url)

		if err != nil {
			lastErr = err

			// se não for erro de rede temporário, não tenta novamente
			if !strings.Contains(err.Error(), "connection reset") &&
				!strings.Contains(err.Error(), "timeout") {
				break
			}
			continue
		}

		defer res.Body.Close()

		// lê o corpo da resposta
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return SiteStatus{
				URL:      url,
				Status:   "Erro",
				Message:  "Erro ao ler resposta da API",
				IsOnline: false,
			}
		}

		// decodifica o JSON
		var apiRes APIResponse
		if err := json.Unmarshal(body, &apiRes); err != nil {
			return SiteStatus{
				URL:      url,
				Status:   res.Status,
				Message:  "Resposta inválida da API",
				IsOnline: false,
			}
		}

		// valida se o backend está realmente online
		online := res.StatusCode == http.StatusOK && apiRes.Ok
		msg := "Backend Online"
		if !online {
			msg = "Backend Offline"
		}

		// usa a mensagem da API se disponível
		if apiRes.Message != "" {
			msg = apiRes.Message
		}

		return SiteStatus{
			URL:      url,
			Status:   res.Status,
			Message:  msg,
			IsOnline: online,
		}
	}

	// todas as tentativas falharam
	return SiteStatus{
		URL:      url,
		Status:   "Erro",
		Message:  formatError(lastErr),
		IsOnline: false,
	}
}
