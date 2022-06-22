package model

import (
	uuid "github.com/satori/go.uuid"
	"net/url"
	"regexp"
	"time"
	apiError "xm/error"
)

// Company contains the data of each company
type Company struct {
	ID        uuid.UUID  `gorm:"type:varchar(36);primary_key;"`
	CreatedAt time.Time  `gorm:"column:createdOn"`
	UpdatedAt time.Time  `gorm:"column:modifiedOn"`
	DeletedAt *time.Time `sql:"index" gorm:"column:deletedOn"`
	Name      string     `gorm:"column:name"`
	Code      string     `gorm:"column:code"`
	Country   string     `gorm:"column:country"`
	Website   string     `gorm:"column:website"`
	Phone     string     `gorm:"column:phone"`
}

// NewCompany creates new company
func NewCompany(name, code, country, website, phone string) (*Company, error) {
	if err := validateCompany(name, code, country, website, phone); err != nil {
		return nil, err
	}

	return &Company{
		ID:      uuid.NewV4(),
		Name:    name,
		Code:    code,
		Country: country,
		Website: website,
		Phone:   phone,
	}, nil
}

// Update updates existing company
func (company *Company) Update(name, code, country, website, phone string) error {
	if err := validateCompany(name, code, country, website, phone); err != nil {
		return err
	}

	company.Name = name
	company.Code = code
	company.Country = country
	company.Website = website
	company.Phone = phone
	return nil
}

func validateCompany(name, code, country, website, phone string) error {
	if len(name) == 0 {
		return apiError.NewInvalidFieldsError(map[string]string{"name": apiError.ErrorCodeRequired})
	}
	if len(code) == 0 {
		return apiError.NewInvalidFieldsError(map[string]string{"code": apiError.ErrorCodeRequired})
	}
	if len(country) == 0 {
		return apiError.NewInvalidFieldsError(map[string]string{"country": apiError.ErrorCodeRequired})
	}
	_, err := url.ParseRequestURI(website)
	if err != nil {
		return apiError.NewInvalidFieldsError(map[string]string{"website": apiError.ErrorCodeInvalidValue})
	}
	if !isPhoneValid(phone) {
		return apiError.NewInvalidFieldsError(map[string]string{"phone": apiError.ErrorCodeInvalidValue})
	}
	return nil
}

func isPhoneValid(phone string) bool {
	re := regexp.MustCompile(`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`)
	return re.MatchString(phone)
}
