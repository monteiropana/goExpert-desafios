package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

//API: protocolo/host/path/???query

const (
	urlApiWeatherAPI   = "http://api.weatherapi.com/v1/current.json?"
	key                = "key=559686bddc72407cbe8173402242101"
	query1             = "&q="
	query2             = "&aqi=no"
	urlApiViaCepPrefix = "https://viacep.com.br/ws/"
	urlSulffix         = "/json/"

	//replace
	//urlReplace = "https://viacep.com.br/ws/CEP/json/"
	// s := strings.Replace(urlReplace, "CEP", cep, 2)
	// fmt.Println(s)
)

func main() {
	var cep string
	fmt.Println("Digite um cep")
	fmt.Scan(&cep)
	if len(cep) != 8 {
		panic("Invalid cep")
	}
	urlCep := urlApiViaCepPrefix + cep + urlSulffix

	response, _ := GetCep(urlCep)
	urlWeather := urlApiWeatherAPI + key + query1 + response.Localidade + query2
	GetWeather(urlWeather)

}

type ResponseWeatherAPI struct {
	Current Current `json:"current"`
}

type Current struct {
	TempC float64 `json:"temp_c"`
	TempF float64 `json:"temp_f"`
}

// type ResponseWeatherAPI struct {
// 	temp_c float32
// 	temp_f float32
// 	//temp_k bool
// }

type ResponseBody struct {
	temp_c float32
	temp_f float32
	//temp_k bool
}

type ResponseCep struct {
	Bairro      string
	Cep         string
	Complemento string
	Ddd         string
	Gia         string
	Ibge        string
	Localidade  string
	Logradouro  string
}

func GetWeather(url string) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Print(err.Error())
	}
	var objectResponse ResponseWeatherAPI
	json.Unmarshal(responseData, &objectResponse)

	fmt.Println("Response Body: ", objectResponse)
}

func GetCep(url string) (ResponseCep, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Print(err.Error())
	}
	var objectResponse ResponseCep
	json.Unmarshal(responseData, &objectResponse)

	fmt.Println("Response Body: ", objectResponse)

	return objectResponse, nil

}

//TODO: criar API e calcular o kelvin
//TODO: response e request
