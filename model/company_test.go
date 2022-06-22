package model

import (
	"testing"
	apiError "xm/error"
)

func Test_validateCompany(t *testing.T) {
	type args struct {
		name    string
		code    string
		country string
		website string
		phone   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr *apiError.ValidationError
	}{
		{
			"+ve:ShouldPassWhenValidDataPassed",
			args{
				name:    "ABC Enterprise",
				code:    "001",
				country: "India",
				website: "https://www.abc.com",
				phone:   "0911094444",
			},
			nil,
		},
		{
			name: "-ve:ShouldFailWhenEmptyNamePassed",
			args: args{
				name:    "",
				code:    "001",
				country: "India",
				website: "https://www.abc.com",
				phone:   "0911094444",
			},
			wantErr: &apiError.ValidationError{ErrorKey: apiError.ErrorCodeInvalidFields, Errors: map[string]string{"name": apiError.ErrorCodeRequired}},
		},
		{
			name: "-ve:ShouldFailWhenEmptyCodePassed",
			args: args{
				name:    "ABC Enterprise",
				code:    "",
				country: "India",
				website: "https://www.abc.com",
				phone:   "0911094444",
			},
			wantErr: &apiError.ValidationError{ErrorKey: apiError.ErrorCodeInvalidFields, Errors: map[string]string{"code": apiError.ErrorCodeRequired}},
		},
		{
			name: "+ve:ShouldFailWhenEmptyCountryPassed",
			args: args{
				name:    "ABC Enterprise",
				code:    "001",
				country: "",
				website: "https://www.abc.com",
				phone:   "0911094444",
			},
			wantErr: &apiError.ValidationError{ErrorKey: apiError.ErrorCodeInvalidFields, Errors: map[string]string{"country": apiError.ErrorCodeRequired}},
		},
		{
			name: "-ve:ShouldFailWhenInvalidWebsitePassed",
			args: args{
				name:    "ABC Enterprise",
				code:    "001",
				country: "India",
				website: "abc.com",
				phone:   "0911094444",
			},
			wantErr: &apiError.ValidationError{ErrorKey: apiError.ErrorCodeInvalidFields, Errors: map[string]string{"website": apiError.ErrorCodeInvalidValue}},
		},
		{
			name: "-ve:ShouldFailWhenInvalidPhonePassed",
			args: args{
				name:    "ABC Enterprise",
				code:    "001",
				country: "India",
				website: "https://www.abc.com",
				phone:   "abc0911094444",
			},
			wantErr: &apiError.ValidationError{ErrorKey: apiError.ErrorCodeInvalidFields, Errors: map[string]string{"phone": apiError.ErrorCodeInvalidValue}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCompany(tt.args.name, tt.args.code, tt.args.country, tt.args.website, tt.args.phone)
			if err != nil && tt.wantErr == nil {
				t.Errorf("validateCompany() got error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr != nil && err == nil {
				t.Errorf("validateCompany() got error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("validateCompany() got error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
