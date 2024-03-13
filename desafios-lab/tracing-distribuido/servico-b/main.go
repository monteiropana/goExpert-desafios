package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
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
	Current  Current  `json:"current"`
	Location Location `json:"location"`
}

type Current struct {
	TempC float64 `json:"temp_c"`
	TempF float64 `json:"temp_f"`
	TempK float64 `json:"temp_k"`
}

type ResponsePostCep struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_c"`
	TempF float64 `json:"temp_f"`
	TempK float64 `json:"temp_k"`
}

type Location struct {
	Name string `json:"name"`
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
	Erro bool
}

func main() {
	log.Print("start servico B")
	InitTracer()
	http.HandleFunc("/clima", Execute)
	log.Fatal(http.ListenAndServe(":8081", nil))
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
	if objectResponseErr.Erro {
		return ResponseCep{}, errors.New("not found")
	}

	var objectResponse ResponseCep
	json.Unmarshal(responseData, &objectResponse)

	fmt.Println("Response Body: ", objectResponse)

	return objectResponse, nil
}

var tracer = otel.Tracer("servico-b")

func Execute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(r.Header))

	_, operationSpan := tracer.Start(ctx, "totalOperation")
	defer operationSpan.End()

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

	_, getWeatherApiSpan := tracer.Start(context.Background(), "api_weather")
	res := GetWeather(urlWeather)
	getWeatherApiSpan.End()

	res.Current.TempK = res.Current.TempC + 273

	objectResponse := ResponsePostCep{
		City:  res.Location.Name,
		TempC: res.Current.TempC,
		TempF: res.Current.TempF,
		TempK: res.Current.TempK + 273,
	}

	returnWeather, _ := json.Marshal(objectResponse)
	w.Write(returnWeather)
}

func InitTracer() {
	ctx := context.Background()

	exp, err := otlptrace.New(
		ctx,
		otlptracegrpc.NewClient(
			otlptracegrpc.WithEndpoint("otel-collector:4317"),
			otlptracegrpc.WithInsecure(),
		),
	)
	if err != nil {
		log.Fatalf("Falha ao criar o exportador OTLP: %v", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("servico-b"),
		)),
	)

	otel.SetTracerProvider(tp)
}
