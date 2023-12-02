package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	urlApiViaCep    = "https://viacep.com.br/ws/86300000/json/"
	urlApiBrasilCep = "https://brasilapi.com.br/api/cep/v1/86300000"
)

//var client = &http.Client{}

func main() {
	c1 := make(chan string) //vazio
	c2 := make(chan string)

	go DoRequest(urlApiViaCep, c1)
	go DoRequest(urlApiBrasilCep, c2)

	select {
	case msg := <-c1:
		fmt.Println("recebido api Via Cep", msg)

	case msg := <-c2:
		fmt.Println("recebido api Brasil Cep", msg)

	case <-time.After(time.Second):
		println("timeout")
	}

}

func DoRequest(url string, channel chan string) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	response, err := http.DefaultClient.Do(request) //executa a request
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()
	bodyBytes, _ := io.ReadAll(response.Body)

	resp := string(bodyBytes) //convertendo pra string
	channel <- resp

}

// fazer as duas requiscoes, recuperar o dado e printar na tela
