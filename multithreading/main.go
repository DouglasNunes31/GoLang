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

type Brasilapi struct {
	Code         string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
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
	chBrasilapi := make(chan Brasilapi)
	chViaCEP := make(chan ViaCEP)

	go BuscaBrasilapi(chBrasilapi)
	go BuscaViaCEP(chViaCEP)

	select {
	case retornoBrasilapi := <-chBrasilapi:
		fmt.Printf("brasilapi: ", retornoBrasilapi)

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

func BuscaBrasilapi(chBrasilapi chan Brasilapi) {
	var brassilAPI Brasilapi
	BuscaDadosApi("https://brasilapi.com.br/api/cep/v2/89010025.json", &brassilAPI)
	chBrasilapi <- brassilAPI
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
