package dto

import "github.com/asaskevich/govalidator"

type (
	AuthRequestDTO struct {
		Password string `json:"password" valid:"required,length(8|255)"`
		Username string `json:"username" valid:"required,matches(^[a-zA-Z0-9_]+$)"`
	}
	AuthResponseDTO struct {
		Token string `json:"token"`
	}
	ReviewDTO struct {
		Mark    uint32 `valid:"int,range(1|10),required"`
		Comment string `valid:"optional,length(10|10000)"`
	}
)

func (authReqDTO *AuthRequestDTO) Validate() []string {
	_, err := govalidator.ValidateStruct(authReqDTO)
	return collectErrors(err)
}

func (reviewDTO *ReviewDTO) Validate() []string {
	_, err := govalidator.ValidateStruct(reviewDTO)
	return collectErrors(err)
}

func collectErrors(err error) []string {
	validationErrors := make([]string, 0)
	if err == nil {
		return validationErrors
	}
	if allErrs, ok := err.(govalidator.Errors); ok {
		for _, fld := range allErrs {
			validationErrors = append(validationErrors, fld.Error())
		}
	}
	return validationErrors
}
