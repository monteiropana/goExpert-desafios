package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	urlServiceB = "http://172.29.0.3:8080/clima?cep="
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
	http.HandleFunc("/consulta", ExecutePost)
	http.ListenAndServe(":8081", nil)

}

func ExecutePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//responseCep := ResponsePostCep{}
	var cep Cep
	err := json.NewDecoder(r.Body).Decode(&cep)
	if err != nil || len(cep.Cep) != 8 {
		http.Error(w, "invalid zipcode", 422)
	}

	request, err := http.NewRequest(http.MethodGet, urlServiceB+cep.Cep, nil)
	if err != nil {
		fmt.Print(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Print(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}

	defer response.Body.Close()
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Print(err.Error())
		http.Error(w, err.Error(), 500)
	}
	var objResp ResponsePostCep
	err = json.Unmarshal(responseData, &objResp)
	if err != nil {
		fmt.Print(err.Error())
		http.Error(w, err.Error(), 500)
	}
	w.Write(responseData)
}
