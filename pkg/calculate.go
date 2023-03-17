package pkg

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
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
	Symbol string  `json:"Symbol"`
	Qty    int     `json:"Qty"`
	ACB    float64 `json:"ACB"`
	CG     float64 `json:"Capital Gains"`
	TACB   float64 `json:"Total Proceeds ACB"`
	TPCD   float64 `json:"Total Proceeds"`
	TCOM   float64 `json:"Total Commision"`
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
		// fmt.Println(t)
		p.TCOM += t.Commission
		if (p.Qty >= 0 && t.ShareQty > 0) || (p.Qty <= 0 && t.ShareQty < 0) {
			p.ACB = (p.ACB + t.Proceeds + t.Commission)
			p.Qty += t.ShareQty
			p.TACB += t.Proceeds
		} else {
			// shareQtyAbs := math.Abs(float64(t.ShareQty))
			p.CG = p.CG + (t.PricePerShare*float64(-t.ShareQty) - t.Commission - (p.ACB / float64(p.Qty) * float64(-t.ShareQty)))
			p.ACB = p.ACB * (float64(p.Qty + t.ShareQty)) / float64(p.Qty)
			p.Qty += t.ShareQty
			p.TPCD += t.Proceeds
		}
		log.Printf("%+v\n", *p)
	}

	result := writeResultJson(positions, transactions)
	return result
}

func writeResultJson(positions map[string]*Position, transactions []*Transaction) []byte {
	positionsJSON := []*Position{{Symbol: "Total"}}

	for _, p := range positions {
		positionsJSON = append(positionsJSON, p)
		positionsJSON[0].Qty += p.Qty
		positionsJSON[0].ACB += p.ACB
		positionsJSON[0].CG += p.CG
		positionsJSON[0].TACB += p.TACB
		positionsJSON[0].TPCD += p.TPCD
		positionsJSON[0].TCOM += p.TCOM
	}

	for _, pj := range positionsJSON {
		pj.ACB = math.Round(pj.ACB*100) / 100
		pj.CG = math.Round(pj.CG*100) / 100
		pj.TACB = math.Round((pj.TACB-pj.ACB)*100) / 100
		pj.TPCD = math.Round(pj.TPCD*100) / 100
		pj.TCOM = math.Round(pj.TCOM*100) / 100
	}

	positionsJSONBytes, err := json.Marshal(positionsJSON)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return positionsJSONBytes
}

// func writeResultCSVFile(positions map[string]*Position, transactions []*Transaction) *csv.Writer {
// 	f, err := os.Create("output.csv")
// 	if err != nil {
// 		log.Println(err)
// 		return nil
// 	}
// 	defer f.Close()

// 	w := csv.NewWriter(f)
// 	defer w.Flush()

// 	w.Write([]string{"Symbol", "Qty", "ACB", "Capital Gains"})

// 	for _, p := range positions {
// 		w.Write([]string{p.Symbol, strconv.Itoa(p.Qty), fmt.Sprintf("%.2f", p.ACB), fmt.Sprintf("%.2f", p.CG)})
// 	}

// 	return w
// }

func readFile(sheet [][]string) ([]*Transaction, error) {
	var transactions []*Transaction

	for i, row := range sheet {
		if i == 0 || row[0] == "" {
			continue
		}
		// log.Println(row)

		t := Transaction{}

		t.Symbol = string(row[0])
		t.DateTime, _ = time.Parse("2006-01-02, 15:04:05", row[1])
		t.ShareQty, _ = strconv.Atoi(row[2])
		t.PricePerShare, _ = strconv.ParseFloat(row[3], 64)
		// t.Proceeds, _ = strconv.ParseFloat(row[4], 64)
		t.Commission, _ = strconv.ParseFloat(row[5], 64)

		t.Proceeds = float64(t.ShareQty) * t.PricePerShare
		t.Commission = math.Abs(t.Commission)

		transactions = append(transactions, &t)
	}

	return transactions, nil
}
