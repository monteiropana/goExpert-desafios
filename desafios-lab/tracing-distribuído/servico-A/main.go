package main

import (
	"encoding/json"
	"net/http"
)

type Cep struct {
	Cep string `json:"cep"`
}

type ResponsePostCep struct {
	TempC float64 `json:"temp_c"`
	TempF float64 `json:"temp_f"`
	TempK float64 `json:"temp_k"`
	City  string  `json:"city"`
}

func main() {
	http.HandleFunc("/consulta", ExecutePost)
	http.ListenAndServe(":8081", nil)

}

func ExecutePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var cep Cep
	err := json.NewDecoder(r.Body).Decode(&cep)
	if err != nil || len(cep.Cep) != 8 {
		http.Error(w, "invalid zipcode", 422)
		return
	}

	var res ResponsePostCep
	returnJson, _ := json.Marshal(res)
	w.Write(returnJson)
}
