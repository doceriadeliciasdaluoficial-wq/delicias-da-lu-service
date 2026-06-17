package cpf

import (
	"strconv"
	"strings"
)

type CPFValidator struct{}

func NewCPFValidator() *CPFValidator {
	return &CPFValidator{}
}

func (cv *CPFValidator) IsValid(cpf string) bool {
	cpf = strings.ReplaceAll(cpf, ".", "")
	cpf = strings.ReplaceAll(cpf, "-", "")
	cpf = strings.TrimSpace(cpf)

	if len(cpf) != 11 {
		return false
	}

	if !isAllDigits(cpf) {
		return false
	}

	if allDigitsEqual(cpf) {
		return false
	}

	firstChecksum := calculateChecksum(cpf[:9], getWeights1())
	if cpf[9:10] != string(rune(firstChecksum+'0')) {
		return false
	}

	secondChecksum := calculateChecksum(cpf[:10], getWeights2())
	if cpf[10:11] != string(rune(secondChecksum+'0')) {
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
	return []int{10, 9, 8, 7, 6, 5, 4, 3, 2}
}

func getWeights2() []int {
	return []int{11, 10, 9, 8, 7, 6, 5, 4, 3, 2}
}
