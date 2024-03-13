package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

const (
	urlServiceB = "http://servico-b:8081/clima?cep="
	//urlServiceB = "http://localhost:8081/clima?cep="
)

type Cep struct {
	Cep string `json:"cep"`
}

type ResponsePostCep struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_c"`
	TempF float64 `json:"temp_f"`
	TempK float64 `json:"temp_k"`
}

func main() {
	log.Print("start servico A")
	InitTracer()
	http.HandleFunc("/consulta", ExecutePost)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

var tracer = otel.Tracer("servico-a")

func ExecutePost(w http.ResponseWriter, r *http.Request) {
	_, validateCEPSpan := tracer.Start(context.Background(), "ConsultaCEP")
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(r.Header))
	_, operationSpan := tracer.Start(ctx, "totalOperation")
	defer operationSpan.End()

	var cep Cep
	err := json.NewDecoder(r.Body).Decode(&cep)
	if err != nil || len(cep.Cep) != 8 {
		http.Error(w, "invalid zipcode", 422)
		return
	}

	validateCEPSpan.End()
	request, err := http.NewRequest(http.MethodGet, urlServiceB+cep.Cep, nil)
	if err != nil {
		http.Error(w, "erro na requisicao", 404)
		return
	}
	_, span := otel.Tracer("servico-a").Start(context.Background(), "callServiceB")
	defer span.End()

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		http.Error(w, "cannot find zipcode", 404)
		return
	}

	defer response.Body.Close()
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		http.Error(w, "erro na leitura dos dados", 404)
		return

	}
	if response.StatusCode != 200 {
		http.Error(w, string(responseData), 404)
		return
	}

	var objResp ResponsePostCep
	err = json.Unmarshal(responseData, &objResp)
	if err != nil {
		http.Error(w, "erro de serialização", 404)
		return
	}
	w.Write(responseData)
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
			semconv.ServiceNameKey.String("servico-a"),
		)),
	)

	otel.SetTracerProvider(tp)
}
