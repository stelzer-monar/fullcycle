package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type APICEPResponse struct {
	Status   int    `json:"status"`
	Code     string `json:"code"`
	State    string `json:"state"`
	City     string `json:"city"`
	District string `json:"district"`
	Address  string `json:"address"`
}

type VIACEPResponse struct {
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
	viachan := make(chan VIACEPResponse)
	apichan := make(chan APICEPResponse)

	ctx := context.Background()

	go getAPICEP(ctx, os.Args[1], apichan)
	go getVIACEP(ctx, os.Args[1], viachan)

	select {
	case res := <-viachan:
		fmt.Printf("Received from VIA CEP: %v\n", res)
	case res := <-apichan:
		fmt.Printf("Received from API CEP: %v\n", res)
	case <-time.After(time.Second):
		fmt.Printf("Timeout\n")
	}
}

func getAPICEP(ctx context.Context, cep string, response chan<- APICEPResponse) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://cdn.apicep.com/file/apicep/"+cep+".json", nil)
	if err != nil {
		fmt.Println("Error connecting to apicep: " + err.Error())
		return
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error connecting to apicep: " + err.Error())
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response from apicep: " + err.Error())
		return
	}

	conversion := APICEPResponse{}

	err = json.Unmarshal(body, &conversion)
	if err != nil {
		fmt.Println("Error reading response from apicep: " + err.Error())
		return
	}

	response <- conversion
}

func getVIACEP(ctx context.Context, cep string, response chan<- VIACEPResponse) {
	req, err := http.NewRequestWithContext(ctx, "GET", "http://viacep.com.br/ws/"+cep+"/json/", nil)
	if err != nil {
		fmt.Println("Error connecting to viacep: " + err.Error())
		return
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error connecting to viacep: " + err.Error())
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response from viacep: " + err.Error())
		return
	}

	conversion := VIACEPResponse{}

	err = json.Unmarshal(body, &conversion)
	if err != nil {
		fmt.Println("Error reading response from viacep: " + err.Error())
		return
	}

	response <- conversion
}
