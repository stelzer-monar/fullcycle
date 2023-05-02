package main

import (
	"context"
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"net/http"
	"time"
)

type Conversion struct {
	UsdBrl CotationFullResponse `json:"USDBRlL`
}
type CotationBidResponse struct {
	Bid string `json:"bid"`
}

type CotationFullResponse struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

func main() {
	http.HandleFunc("/cotacao", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/goexpert")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
		return
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(r.Context(), time.Millisecond*200)
	defer cancel()

	cotation, err := getCotation(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
		return
	}

	ctx, cancel = context.WithTimeout(r.Context(), time.Millisecond*10)
	defer cancel()
	err = saveToBD(ctx, cotation, db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(CotationBidResponse{cotation.Bid})
}

func saveToBD(ctx context.Context, cttn CotationFullResponse, db *sql.DB) error {
	stmt, err := db.PrepareContext(ctx, "insert into cotations(code, codein, name, high, low, varBid, pctChange, bid, ask, timestamp, createDate) values (?,?,?,?,?,?,?,?,?,?,?);")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, cttn.Code, cttn.Codein, cttn.Name, cttn.High, cttn.Low, cttn.VarBid, cttn.PctChange, cttn.Bid, cttn.Ask, cttn.Timestamp, cttn.CreateDate)
	if err != nil {
		return err
	}
	return nil
}

func getCotation(ctx context.Context) (CotationFullResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return CotationFullResponse{}, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return CotationFullResponse{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return CotationFullResponse{}, err
	}

	conversion := Conversion{}

	err = json.Unmarshal(body, &conversion)
	if err != nil {
		return CotationFullResponse{}, err
	}

	return conversion.UsdBrl, nil
}
