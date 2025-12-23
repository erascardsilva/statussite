package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {
	site1 := "https://www.electrocode.com.br"
	site2 := "https://www.vmi-informatica.com.br"
	site3 := "https://www.electrocode.com.br/app/authentication/sign-in/"

	var wg sync.WaitGroup
	wg.Add(3)

	fmt.Println("Verificando status dos sites VMI e Electrocode")
	go status(site1, &wg)
	go status(site2, &wg)
	go status(site3, &wg)

	wg.Wait()
}

func status(url string, wg *sync.WaitGroup) {
	defer wg.Done()
	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Head(url)
	if err != nil {
		fmt.Println("Erro ao verificar o status do site:", url)
		return
	}
	defer res.Body.Close()
	fmt.Println("---------------------------------")
	fmt.Println("Site:", url, res.Status)
	if res.StatusCode == http.StatusOK {
		fmt.Println("Status -- ")
		fmt.Println("Esta online")
	} else {
		fmt.Println(" ")
		fmt.Println("Esta offline")
	}

}
