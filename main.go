package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	ViaCEP_URL    = "https://viacep.com.br/ws/03734130/json/"
	BrasilAPI_URL = "https://brasilapi.com.br/api/cep/v1/03734130"
)

type BrasilAPIStruct struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

type ViaCEPStruct struct {
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
	channelViaCEP := make(chan ViaCEPStruct)
	channelBrasilAPI := make(chan BrasilAPIStruct)

	go getViaCep(channelViaCEP)
	go getBrasilAPI(channelBrasilAPI)

	select {
	case response := <-channelViaCEP:
		fmt.Printf("ViaCEP: %+v\n", response)
	case response := <-channelBrasilAPI:
		fmt.Printf("BrasilAPI: %+v\n", response)
	case <-time.After(time.Second):
		println("Timeout")
	}
}

func getViaCep(channelAPI chan<- ViaCEPStruct) {
	var response ViaCEPStruct
	// time.Sleep(time.Second * 2)
	requestApi(ViaCEP_URL, &response)
	channelAPI <- response
}

func getBrasilAPI(channelAPI chan<- BrasilAPIStruct) {
	var response BrasilAPIStruct
	// time.Sleep(time.Second * 2)
	requestApi(BrasilAPI_URL, &response)
	channelAPI <- response
}

func requestApi(url string, response interface{}) error {
	request, errorRequest := http.Get(url)
	if errorRequest != nil {
		return errorRequest
	}
	defer request.Body.Close()
	body, errorBody := io.ReadAll(request.Body)
	if errorBody != nil {
		return errorBody
	}
	errorUnmarshal := json.Unmarshal(body, &response)
	if errorUnmarshal != nil {
		return errorUnmarshal
	}
	return nil
}
