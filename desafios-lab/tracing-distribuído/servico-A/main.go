package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

const (
	urlApiWeatherAPI   = "http://api.weatherapi.com/v1/current.json?"
	key                = "key=559686bddc72407cbe8173402242101"
	query1             = "&q="
	query2             = "&aqi=no"
	urlApiViaCepPrefix = "https://viacep.com.br/ws/"
	urlSulffix         = "/json/"
)

type ResponseWeatherAPI struct {
	Current Current `json:"current"`
}

type Current struct {
	TempC float64 `json:"temp_c"`
	TempF float64 `json:"temp_f"`
	TempK float64 `json:"temp_k"`
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

type ResponseError struct {
	Erro string
}

type ResponsePostCep struct {
	TempC float64 `json:"temp_c"`
	TempF float64 `json:"temp_f"`
	TempK float64 `json:"temp_k"`
	City  string  `json:"city"`
}

func main() {
	http.HandleFunc("/clima", Execute)
	http.HandleFunc("/consulta", ExecutePost)
	http.ListenAndServe(":8080", nil)

}

func ExecutePost(w http.ResponseWriter, r *http.Request) {
	var cep ResponsePostCep

	err := json.NewDecoder(r.Body).Decode(&cep)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "Response Body: %+v", cep)
}

func GetWeather(url string) ResponseWeatherAPI {
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

	return objectResponse
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
		return ResponseCep{}, err
	}

	var objectResponseErr ResponseError
	json.Unmarshal(responseData, &objectResponseErr)
	if objectResponseErr.Erro == "true" {
		return ResponseCep{}, errors.New("not found")
	}

	var objectResponse ResponseCep
	json.Unmarshal(responseData, &objectResponse)

	fmt.Println("Response Body: ", objectResponse)

	return objectResponse, nil
}

func Execute(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	cep := v.Get("cep")

	if len(cep) != 8 {
		http.Error(w, "invalid zipcode", 422)
		return
	}

	urlCep := urlApiViaCepPrefix + cep + urlSulffix

	response, err := GetCep(urlCep)
	if err != nil {
		http.Error(w, "cannot find zipcode", 404)
		return
	}
	urlWeather := urlApiWeatherAPI + key + query1 + url.QueryEscape(response.Localidade) + query2

	res := GetWeather(urlWeather)
	res.Current.TempK = res.Current.TempC + 273

	returnWeather, _ := json.Marshal(res)
	w.Write(returnWeather)

}
