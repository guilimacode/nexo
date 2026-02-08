package utils

import (
	"errors"
	"regexp"

	"github.com/klassmann/cpfcnpj"
)

func FormatAndValidateCpfCnpj(doc string) (string, error) {
	formatted := regexp.MustCompile(`\D`).ReplaceAllString(doc, "")

	switch len(formatted) {
	case 11:
		if !cpfcnpj.ValidateCPF(formatted) {
			return "", errors.New("CPF inválido")
		}
		return formatted, nil
	case 14:
		if !cpfcnpj.ValidateCNPJ(formatted) {
			return "", errors.New("CNPJ inválido")
		}
		return formatted, nil
	default:
		return "", errors.New("tamanho do documento inválido")
	}
}
