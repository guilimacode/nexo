package dtos

type CreateAddressDTO struct {
	Street       string `json:"street" validate:"required"`
	Complement   string `json:"complement"`
	Number       int    `json:"number" validate:"required,numeric"`
	Neighborhood string `json:"neighborhood" validate:"required"`
	City         string `json:"city" validate:"required"`
	State        string `json:"state" validate:"required,len=2,uppercase"`
	ZipCode      string `json:"zip_code" validate:"required,len=8,numeric"`
}
