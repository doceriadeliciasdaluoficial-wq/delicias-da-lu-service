package documentvalidator

import (
	"testing"
)

func TestDocumentValidatorIsValid(t *testing.T) {
	validator := NewDocumentValidator()

	tests := []struct {
		name      string
		document  Document
		length    int
		wantValid bool
		wantError bool
	}{
		{
			name:      "Valid CPF",
			document:  "11144477735",
			length:    11,
			wantValid: true,
			wantError: false,
		},
		{
			name:      "Valid CNPJ",
			document:  "11222333000181",
			length:    14,
			wantValid: true,
			wantError: false,
		},
		{
			name:      "Invalid CPF checksum",
			document:  "11144477734",
			length:    11,
			wantValid: false,
			wantError: false,
		},
		{
			name:      "Invalid CNPJ checksum",
			document:  "11222333000180",
			length:    14,
			wantValid: false,
			wantError: false,
		},
		{
			name:      "CPF all same digits",
			document:  "11111111111",
			length:    11,
			wantValid: false,
			wantError: false,
		},
		{
			name:      "CNPJ all same digits",
			document:  "11111111111111",
			length:    14,
			wantValid: false,
			wantError: false,
		},
		{
			name:      "Invalid length 10",
			document:  "1234567890",
			length:    10,
			wantValid: false,
			wantError: true,
		},
		{
			name:      "Invalid length 12",
			document:  "123456789012",
			length:    12,
			wantValid: false,
			wantError: true,
		},
		{
			name:      "Invalid length 13",
			document:  "1234567890123",
			length:    13,
			wantValid: false,
			wantError: true,
		},
		{
			name:      "Invalid length 15",
			document:  "123456789012345",
			length:    15,
			wantValid: false,
			wantError: true,
		},
		{
			name:      "Auto-detect valid CPF",
			document:  "11144477735",
			length:    0,
			wantValid: true,
			wantError: false,
		},
		{
			name:      "Auto-detect valid CNPJ",
			document:  "11222333000181",
			length:    0,
			wantValid: true,
			wantError: false,
		},
		{
			name:      "Auto-detect invalid length",
			document:  "1234567890",
			length:    0,
			wantValid: false,
			wantError: true,
		},
		{
			name:      "Empty document",
			document:  "",
			length:    0,
			wantValid: false,
			wantError: true,
		},
		{
			name:      "CPF with formatting",
			document:  "111.444.777-35",
			length:    11,
			wantValid: true,
			wantError: false,
		},
		{
			name:      "CNPJ with formatting",
			document:  "11.222.333/0001-81",
			length:    14,
			wantValid: true,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validator.IsValid(tt.document, tt.length)

			if (err != nil) != tt.wantError {
				t.Errorf("IsValid() error = %v, wantError = %v", err, tt.wantError)
				return
			}

			if got != tt.wantValid {
				t.Errorf("IsValid() got = %v, want = %v", got, tt.wantValid)
			}
		})
	}
}

func TestDocumentValidatorSanitize(t *testing.T) {
	validator := NewDocumentValidator()

	tests := []struct {
		name    string
		input   Document
		wantOut Document
	}{
		{
			name:    "CPF with dots and dash",
			input:   "111.444.777-35",
			wantOut: "11144477735",
		},
		{
			name:    "CNPJ with dots slash and dash",
			input:   "11.222.333/0001-81",
			wantOut: "11222333000181",
		},
		{
			name:    "Simple numeric CPF",
			input:   "11144477735",
			wantOut: "11144477735",
		},
		{
			name:    "Simple numeric CNPJ",
			input:   "11222333000181",
			wantOut: "11222333000181",
		},
		{
			name:    "Empty document",
			input:   "",
			wantOut: "",
		},
		{
			name:    "Document with spaces",
			input:   "111 444 777 35",
			wantOut: "11144477735",
		},
		{
			name:    "Document with all special chars",
			input:   "***111%%%444###777@@@35",
			wantOut: "11144477735",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validator.Sanitize(tt.input)
			if got != tt.wantOut {
				t.Errorf("Sanitize() got = %v, want = %v", got, tt.wantOut)
			}
		})
	}
}

func BenchmarkDocumentValidatorIsValid(b *testing.B) {
	validator := NewDocumentValidator()
	document := Document("11144477735")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validator.IsValid(document, 11)
	}
}

func BenchmarkDocumentValidatorSanitize(b *testing.B) {
	validator := NewDocumentValidator()
	document := Document("111.444.777-35")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validator.Sanitize(document)
	}
}
