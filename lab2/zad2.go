package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sort"
)

type Product struct {
	CodeCPV   string
	DE_Label  string
	EN_Label  string
	ES_Label  string
	FR_Label  string
	PT_Label  string
	ShortCode string
}


func loadCSV(filePath string) ([]Product, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';' // Separator ;

	var products []Product
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	// Iteracja po rekordach i tworzenie struktur product
	for _, record := range records {
		if len(record) == 7 {
			product := Product{
				CodeCPV:   record[0],
				DE_Label:  record[1],
				EN_Label:  record[2],
				ES_Label:  record[3],
				FR_Label:  record[4],
				PT_Label:  record[5],
				ShortCode: record[6],
			}
			products = append(products, product)
		}
	}
	return products, nil
}

// Funkcja sortująca po CodeCPV
func sortByCodeCPV(products []Product) {
	sort.Slice(products, func(i, j int) bool {
		return products[i].CodeCPV < products[j].CodeCPV
	})
}

// Funkcja sortująca po EN_Label
func sortByENLabel(products []Product) {
	sort.Slice(products, func(i, j int) bool {
		return products[i].EN_Label < products[j].EN_Label
	})
}

func printStatistics(products []Product) {
	countByLanguage := map[string]int{
		"DE": 0,
		"EN": 0,
		"ES": 0,
		"FR": 0,
		"PT": 0,
	}

	// Zliczamy prod ktore maja nazwe w danym jezyku
	for _, product := range products {
		if product.DE_Label != "" {
			countByLanguage["DE"]++
		}
		if product.EN_Label != "" {
			countByLanguage["EN"]++
		}
		if product.ES_Label != "" {
			countByLanguage["ES"]++
		}
		if product.FR_Label != "" {
			countByLanguage["FR"]++
		}
		if product.PT_Label != "" {
			countByLanguage["PT"]++
		}
	}

	fmt.Println("Statystyka liczby produktów według języka:")
	for lang, count := range countByLanguage {
		fmt.Printf("%s: %d\n", lang, count)
	}
}

func main() {
	products, err := loadCSV("nomenclature-cpv.csv")
	if err != nil {
		log.Fatalf("Błąd podczas wczytywania pliku CSV: %v", err)
	}

	// ile prod zostalo zaladowanych
	fmt.Printf("Załadowano %d produktów\n", len(products))

	sortByCodeCPV(products)
	fmt.Println("Produkty posortowane po CodeCPV:")
	for _, product := range products {
		fmt.Println(product.CodeCPV, product.EN_Label)
	}

	sortByENLabel(products)
	fmt.Println("\nProdukty posortowane po EN_Label:")
	for _, product := range products {
		fmt.Println(product.EN_Label, product.CodeCPV)
	}

	printStatistics(products)
}
