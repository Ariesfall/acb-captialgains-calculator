package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Transaction struct {
	Symbol        string
	DateTime      time.Time
	ShareQty      int
	PricePerShare float64
	Proceeds      float64
	Commission    float64
}

type Position struct {
	Symbol string
	Qty    int
	ACB    float64
}

func Calculate(sheet [][]string) []byte {
	transactions, err := readFile(sheet)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	positions := make(map[string]*Position)

	for _, t := range transactions {
		p, ok := positions[t.Symbol]
		if !ok {
			p = &Position{Symbol: t.Symbol}
			positions[t.Symbol] = p
		}

		p.Qty += t.ShareQty
		p.ACB += t.Proceeds + t.Commission
		p.ACB /= float64(p.Qty)
	}

	for _, t := range transactions {
		p := positions[t.Symbol]

		if t.ShareQty > 0 {
			p.ACB = (p.ACB*float64(p.Qty-t.ShareQty) + (t.Proceeds + t.Commission)) / float64(p.Qty)
			p.Qty += t.ShareQty
		} else {
			p.ACB = (p.ACB*float64(p.Qty) - t.Proceeds) / float64(p.Qty-t.ShareQty)
			p.Qty += t.ShareQty
		}
	}

	result := writeResultJson(positions, transactions)
	fmt.Println(string(result))
	return result
}

func writeResultJson(positions map[string]*Position, transactions []*Transaction) []byte {
	var positionsJSON []map[string]interface{}

	for _, p := range positions {
		cg := 0.0
		if p.Qty > 0 {
			cg = (p.ACB * float64(p.Qty)) - p.ACB*float64(len(transactions))
		}

		positionJSON := map[string]interface{}{
			"Symbol":        p.Symbol,
			"Qty":           p.Qty,
			"ACB":           fmt.Sprintf("%.2f", p.ACB),
			"Capital Gains": fmt.Sprintf("%.2f", cg),
		}

		positionsJSON = append(positionsJSON, positionJSON)
	}

	positionsJSONBytes, err := json.Marshal(positionsJSON)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return positionsJSONBytes
}

func writeResultCSVFile(positions map[string]*Position, transactions []*Transaction) *csv.Writer {
	f, err := os.Create("output.csv")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	w.Write([]string{"Symbol", "Qty", "ACB", "Capital Gains"})

	for _, p := range positions {
		cg := 0.0
		if p.Qty > 0 {
			cg = (p.ACB * float64(p.Qty)) - p.ACB*float64(len(transactions))
		}

		w.Write([]string{p.Symbol, strconv.Itoa(p.Qty), fmt.Sprintf("%.2f", p.ACB), fmt.Sprintf("%.2f", cg)})
	}

	return w
}

func readFile(sheet [][]string) ([]*Transaction, error) {
	var transactions []*Transaction

	for i, row := range sheet {
		if i == 0 {
			continue
		}

		t := Transaction{}

		t.Symbol = string(row[0])
		t.DateTime, _ = time.Parse("YYYY-MM-DD, hh:mm:ss", row[1])
		t.ShareQty, _ = strconv.Atoi(row[2])
		t.PricePerShare, _ = strconv.ParseFloat(row[3], 8)
		t.Proceeds, _ = strconv.ParseFloat(row[4], 8)
		t.Commission, _ = strconv.ParseFloat(row[5], 8)

		transactions = append(transactions, &t)
	}

	return transactions, nil
}
