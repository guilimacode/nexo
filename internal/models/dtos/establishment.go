package dtos

type EstablishmentInputDTO struct {
	Name     string           `json:"name" validate:"required,min=3"`
	Nickname string           `json:"nickname"`
	Document string           `json:"document" validate:"required,min=11,max=18"`
	Address  CreateAddressDTO `json:"address" validate:"required"`
}
