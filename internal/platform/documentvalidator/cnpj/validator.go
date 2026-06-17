package cnpj

import (
	"strconv"
	"strings"
)

type CNPJValidator struct{}

func NewCNPJValidator() *CNPJValidator {
	return &CNPJValidator{}
}

func (cv *CNPJValidator) IsValid(cnpj string) bool {
	cnpj = strings.ReplaceAll(cnpj, ".", "")
	cnpj = strings.ReplaceAll(cnpj, "/", "")
	cnpj = strings.ReplaceAll(cnpj, "-", "")
	cnpj = strings.TrimSpace(cnpj)

	if len(cnpj) != 14 {
		return false
	}

	if !isAllDigits(cnpj) {
		return false
	}

	if allDigitsEqual(cnpj) {
		return false
	}

	firstChecksum := calculateChecksum(cnpj[:12], getWeights1())
	if cnpj[12:13] != string(rune(firstChecksum+'0')) {
		return false
	}

	secondChecksum := calculateChecksum(cnpj[:13], getWeights2())
	if cnpj[13:14] != string(rune(secondChecksum+'0')) {
		return false
	}

	return true
}

func isAllDigits(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func allDigitsEqual(s string) bool {
	if len(s) == 0 {
		return false
	}
	first := rune(s[0])
	for _, c := range s {
		if c != first {
			return false
		}
	}
	return true
}

func calculateChecksum(digits string, weights []int) int {
	sum := 0
	for i, w := range weights {
		digit, _ := strconv.Atoi(string(digits[i]))
		sum += digit * w
	}
	remainder := sum % 11
	if remainder < 2 {
		return 0
	}
	return 11 - remainder
}

func getWeights1() []int {
	return []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
}

func getWeights2() []int {
	return []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
}
