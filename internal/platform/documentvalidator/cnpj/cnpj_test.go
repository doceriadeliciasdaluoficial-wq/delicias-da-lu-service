package cnpj

import (
	"testing"
)

func TestCNPJValidatorIsValid(t *testing.T) {
	validator := NewCNPJValidator()

	tests := []struct {
		name string
		cnpj string
		want bool
	}{
		{
			name: "Valid CNPJ",
			cnpj: "11222333000181",
			want: true,
		},
		{
			name: "Valid CNPJ with formatting",
			cnpj: "11.222.333/0001-81",
			want: true,
		},
		{
			name: "Invalid CNPJ checksum",
			cnpj: "11222333000180",
			want: false,
		},
		{
			name: "CNPJ all same digits",
			cnpj: "11111111111111",
			want: false,
		},
		{
			name: "CNPJ all zeros",
			cnpj: "00000000000000",
			want: false,
		},
		{
			name: "CNPJ too short",
			cnpj: "1122233300018",
			want: false,
		},
		{
			name: "CNPJ too long",
			cnpj: "112223330001811",
			want: false,
		},
		{
			name: "Empty CNPJ",
			cnpj: "",
			want: false,
		},
		{
			name: "CNPJ with letters",
			cnpj: "11.222.333/000A-81",
			want: false,
		},
		{
			name: "CNPJ with spaces",
			cnpj: "11 222 333 000 180",
			want: false,
		},
		{
			name: "Valid CNPJ variant",
			cnpj: "11222333000180",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validator.IsValid(tt.cnpj)
			if got != tt.want {
				t.Errorf("IsValid() got = %v, want = %v", got, tt.want)
			}
		})
	}
}

func BenchmarkCNPJValidatorIsValid(b *testing.B) {
	validator := NewCNPJValidator()
	cnpj := "11222333000181"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validator.IsValid(cnpj)
	}
}
