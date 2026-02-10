package dtos

type SignupDTO struct {
	FullName        string `json:"full_name"    validate:"required,min=3"`
	Email           string `json:"email"        validate:"required,email"`
	Password        string `json:"password"     validate:"required,min=6"`
	CompanyName     string `json:"company_name" validate:"required"`
	CompanyDocument string `json:"document"     validate:"required,min=11,max=18"`
}

type LoginInput struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
