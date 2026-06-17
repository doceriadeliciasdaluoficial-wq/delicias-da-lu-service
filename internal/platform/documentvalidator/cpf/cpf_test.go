package cpf

import (
	"testing"
)

func TestCPFValidatorIsValid(t *testing.T) {
	validator := NewCPFValidator()

	tests := []struct {
		name string
		cpf  string
		want bool
	}{
		{
			name: "Valid CPF",
			cpf:  "11144477735",
			want: true,
		},
		{
			name: "Valid CPF with formatting",
			cpf:  "111.444.777-35",
			want: true,
		},
		{
			name: "Invalid CPF checksum",
			cpf:  "11144477734",
			want: false,
		},
		{
			name: "CPF all same digits",
			cpf:  "11111111111",
			want: false,
		},
		{
			name: "CPF all zeros",
			cpf:  "00000000000",
			want: false,
		},
		{
			name: "CPF too short",
			cpf:  "1114447773",
			want: false,
		},
		{
			name: "CPF too long",
			cpf:  "111444777355",
			want: false,
		},
		{
			name: "Empty CPF",
			cpf:  "",
			want: false,
		},
		{
			name: "CPF with letters",
			cpf:  "111.444.77A-35",
			want: false,
		},
		{
			name: "CPF with spaces",
			cpf:  "111 444 777 34",
			want: false,
		},
		{
			name: "Valid CPF variant",
			cpf:  "12345678909",
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validator.IsValid(tt.cpf)
			if got != tt.want {
				t.Errorf("IsValid() got = %v, want = %v", got, tt.want)
			}
		})
	}
}

func BenchmarkCPFValidatorIsValid(b *testing.B) {
	validator := NewCPFValidator()
	cpf := "11144477735"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validator.IsValid(cpf)
	}
}
