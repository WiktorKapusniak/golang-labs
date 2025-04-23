package main

import (
	"fmt"
	"math/rand"
	"time"
)
func GenerujPESEL(birthDate time.Time, gender string) [11]int {
	var pesel [11]int

	year := birthDate.Year()
	month := int(birthDate.Month())
	day := birthDate.Day()

	// Wyznaczanie wartości miesiąca zależnie od przedziału wiekowego
	switch {
	case year >= 1800 && year <= 1899:
		month += 80
	case year >= 2000 && year <= 2099:
		month += 20
	case year >= 2100 && year <= 2199:
		month += 40
	case year >= 2200 && year <= 2299:
		month += 60
	}

	yearShort := year % 100
	pesel[0] = yearShort / 10
	pesel[1] = yearShort % 10
	pesel[2] = month / 10
	pesel[3] = month % 10
	pesel[4] = day / 10
	pesel[5] = day % 10

	// Generowanie losowego numeru porządkowego
	serial := rand.Intn(1000) // Liczba trzycyfrowa
	pesel[6] = serial / 100
	pesel[7] = (serial / 10) % 10
	pesel[8] = serial % 10

	// Określenie cyfry płci (ostatnia cyfra)
	if gender == "M" {
		pesel[9] = rand.Intn(5)*2 + 1 // Liczba nieparzysta
	} else {
		pesel[9] = rand.Intn(5) * 2 // Liczba parzysta
	}

	pesel[10] = ObliczCyfreKontrolna(pesel)

	return pesel
}

func ObliczCyfreKontrolna(pesel [11]int) int {
	wagi := []int{1, 3, 7, 9, 1, 3, 7, 9, 1, 3}
	suma := 0

	for i := 0; i < 10; i++ {
		suma += pesel[i] * wagi[i]
	}

	kontrolna := (10 - (suma % 10)) % 10
	return kontrolna
}

func WeryfikujPESEL(pesel [11]int) bool {
	return pesel[10] == ObliczCyfreKontrolna(pesel)
}

func main() {
	birthDate := time.Date(2005, 6, 2, 0, 0, 0, 0, time.FixedZone("CET", 3600))
	pesel := GenerujPESEL(birthDate, "M")

	fmt.Print("Wygenerowany PESEL: ")
	for _, digit := range pesel {
		fmt.Print(digit)
	}
	fmt.Println()
	fmt.Println("Czy numer PESEL jest poprawny?", WeryfikujPESEL(pesel))
}
