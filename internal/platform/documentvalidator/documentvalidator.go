package documentvalidator

import (
	"regexp"
	"strconv"
	"strings"
)

type Document string

type DocumentValidator interface {
	IsValid(Document, int) (bool, error)
	Sanitize(Document) Document
}

type documentValidatorImpl struct{}

func NewDocumentValidator() DocumentValidator {
	return &documentValidatorImpl{}
}

func (ref *documentValidatorImpl) IsValid(document Document, length int) (bool, error) {
	sanitized := ref.Sanitize(document)

	if length == 0 {
		length = len(sanitized)
	}

	switch length {
	case 11:
		return ref.validateCPF(sanitized)
	case 14:
		return ref.validateCNPJ(sanitized)
	default:
		return false, ErrInvalidDocumentLenght
	}
}

func (ref *documentValidatorImpl) Sanitize(document Document) Document {
	re := regexp.MustCompile(`[^0-9]`)
	return Document(re.ReplaceAllString(string(document), ""))
}

func (ref *documentValidatorImpl) validateCPF(document Document) (bool, error) {
	cpf := strings.TrimSpace(string(document))

	if len(cpf) != 11 {
		return false, ErrInvalidDocumentLenght
	}

	if !isNumeric(cpf) {
		return false, ErrInvalidDocumentLenght
	}

	if allDigitsSame(cpf) {
		return false, nil
	}

	firstChecksum := calculateCPFChecksum(cpf[:9], cpfWeights1())
	if cpf[9:10] != string(rune(firstChecksum+'0')) {
		return false, nil
	}

	secondChecksum := calculateCPFChecksum(cpf[:10], cpfWeights2())
	if cpf[10:11] != string(rune(secondChecksum+'0')) {
		return false, nil
	}

	return true, nil
}

func (ref *documentValidatorImpl) validateCNPJ(document Document) (bool, error) {
	cnpj := strings.TrimSpace(string(document))

	if len(cnpj) != 14 {
		return false, ErrInvalidDocumentLenght
	}

	if !isNumeric(cnpj) {
		return false, ErrInvalidDocumentLenght
	}

	if allDigitsSame(cnpj) {
		return false, nil
	}

	firstChecksum := calculateCNPJChecksum(cnpj[:12], cnpjWeights1())
	if cnpj[12:13] != string(rune(firstChecksum+'0')) {
		return false, nil
	}

	secondChecksum := calculateCNPJChecksum(cnpj[:13], cnpjWeights2())
	if cnpj[13:14] != string(rune(secondChecksum+'0')) {
		return false, nil
	}

	return true, nil
}

func isNumeric(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func allDigitsSame(s string) bool {
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

func calculateCPFChecksum(digits string, weights []int) int {
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

func calculateCNPJChecksum(digits string, weights []int) int {
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

func cpfWeights1() []int {
	return []int{10, 9, 8, 7, 6, 5, 4, 3, 2}
}

func cpfWeights2() []int {
	return []int{11, 10, 9, 8, 7, 6, 5, 4, 3, 2}
}

func cnpjWeights1() []int {
	return []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
}

func cnpjWeights2() []int {
	return []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
}
