package dtos

type CreateUserDTO struct {
	FullName        string `json:"full_name"        validate:"required,min=3"`
	Email           string `json:"email"            validate:"required,email"`
	Password        string `json:"password"         validate:"required,min=6"`
	Role            string `json:"role"             validate:"required,oneof=owner manager seller"`
	EstablishmentID *int64 `json:"establishment_id"`
}

type ChangePasswordDTO struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

type UpdateUserDTO struct {
	FullName string `json:"full_name" validate:"required,min=3"`
	Email    string `json:"email"     validate:"required,email"`
	Role     string `json:"role"      validate:"required,oneof=manager seller"`
}
