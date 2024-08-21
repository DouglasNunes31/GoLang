package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type RetornoApi interface{}

type ApiCEP struct {
	Code       string `json:"code"`
	State      string `json:"state"`
	City       string `json:"city"`
	District   string `json:"district"`
	Address    string `json:"address"`
	Status     int    `json:"status"`
	Ok         bool   `json:"ok"`
	StatusText string `json:"statusText"`
}
type ViaCEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func main() {
	chApiCEP := make(chan ApiCEP)
	chViaCEP := make(chan ViaCEP)

	go BuscaApiCEP(chApiCEP)
	go BuscaViaCEP(chViaCEP)

	select {
	//case retornoApiCEP := <-chApiCEP:
	//		fmt.Printf("brasilapi: ", retornoApiCEP)

	case retornoViaCEP := <-chViaCEP:
		fmt.Printf("ViaCEP: ", retornoViaCEP)

	case <-time.After(time.Second):
		fmt.Printf("TimeOut")

	}

}

func BuscaViaCEP(chBuscaViaCep chan ViaCEP) {
	var viaCEP ViaCEP
	BuscaDadosApi("https://viacep.com.br/ws/70070080/json/", &viaCEP)
	chBuscaViaCep <- viaCEP
}

func BuscaApiCEP(chBuscaApiCep chan ApiCEP) {
	var apiCEP ApiCEP
	BuscaDadosApi("https://brasilapi.com.br/api/cep/v1/01153000.json", &apiCEP)
	chBuscaApiCep <- apiCEP
}

func BuscaDadosApi(url string, res RetornoApi) error {
	req, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao fazer requisicao: %v \n", err)
	}
	defer req.Body.Close()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao ler resposta: %v\n", err)
	}
	err = json.Unmarshal(body, res)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao fazer o parse da resposta %v\n", err)
	}
	return nil

}
