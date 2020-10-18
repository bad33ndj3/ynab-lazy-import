package main

import (
	"encoding/csv"
	"fmt"
	"github.com/gocarina/gocsv"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type INGExport struct {
	Datum            int    `csv:"Datum"`
	NaamOmschrijving string `csv:"Naam / Omschrijving"`
	Rekening         string `csv:"Rekening"`
	Tegenrekening    string `csv:"Tegenrekening"`
	Code             string `csv:"Code"`
	AfBij            string `csv:"Af Bij"`
	BedragEUR        string `csv:"Bedrag (EUR)"`
	Mutatiesoort     string `csv:"Mutatiesoort"`
	Mededelingen     string `csv:"Mededelingen"`
	SaldoNaMutatie   string `csv:"Saldo na mutatie"`
	Tag              string `csv:"Tag"`
}

var errFailedToGetPath error = fmt.Errorf("failed to get path")

func getLines(account string, path string) ([]*INGExport, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, errFailedToGetPath
	}
	var exportLines []*INGExport
	var files []string
	err := filepath.Walk(path, visit(&files))
	if err != nil {
		return nil, fmt.Errorf("failed to check files in download path: %w", err)
	}

	// check for bank export files
	for _, file := range files {
		if strings.Contains(file, ".csv") && strings.Contains(file, account) {
			exportFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, os.ModePerm)
			if err != nil {
				panic(err)
			}
			var fileExportLines []*INGExport

			gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
				r := csv.NewReader(in)
				r.Comma = ';'
				return r // Allows use pipe as delimiter
			})

			if err := gocsv.UnmarshalFile(exportFile, &fileExportLines); err != nil {
				return nil, fmt.Errorf("failed to unmarshal csv: %w", err)
			}
			exportLines = append(exportLines, fileExportLines...)

			err = exportFile.Close()
			if err != nil {
				return nil, fmt.Errorf("failed to close export file: %w", err)
			}
		}
	}
	return exportLines, nil
}

func visit(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		*files = append(*files, path)
		return nil
	}
}
