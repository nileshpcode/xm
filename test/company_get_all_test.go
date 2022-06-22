package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"xm/model"
)

func TestGetAll(t *testing.T) {
	testApplication.PrepareEmptyTables()

	testDataCompany1 := addCompanyToDB(t, "ABC Enterprise", "001", "India", "https://www.abc.com", "990100000")
	testDataCompany2 := addCompanyToDB(t, "XYZ Enterprise", "002", "US", "https://www.xyz.com", "800100009")
	testDataCompany3 := addCompanyToDB(t, "123 Enterprise", "003", "India", "https://www.123.com", "980100010")

	type queryParams struct {
		name    string
		code    string
		country string
		website string
		phone   string
	}

	tests := []struct {
		name              string
		qp                queryParams
		expectedCompanies []*model.Company
	}{
		{"+ve:ShouldGetAllCompany", queryParams{}, []*model.Company{testDataCompany1, testDataCompany2, testDataCompany3}},
		{"+ve:ShouldGetCompaniesWithNameQueryParam", queryParams{name: "123 Enterprise"}, []*model.Company{testDataCompany3}},
		{"+ve:ShouldGetCompaniesWithCodeQueryParam", queryParams{code: "001"}, []*model.Company{testDataCompany1}},
		{"+ve:ShouldGetCompaniesWithCountryQueryParam", queryParams{country: "India"}, []*model.Company{testDataCompany1, testDataCompany3}},
		{"+ve:ShouldGetCompaniesWithWebsiteQueryParam", queryParams{website: "https://www.xyz.com"}, []*model.Company{testDataCompany2}},
		{"+ve:ShouldGetCompaniesWithPhoneQueryParam", queryParams{phone: "980100010"}, []*model.Company{testDataCompany3}},
		{"-ve:ShouldNotGetWithMatchingNameButNotMatchingCountryQP", queryParams{name: "XYZ Enterprise", country: "India"}, []*model.Company{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			apiURL := fmt.Sprintf("/api/companies?name=%s&code=%s&country=%s&website=%s&phone=%s", tt.qp.name, tt.qp.code, tt.qp.country, tt.qp.website, tt.qp.phone)
			response := callAPI(http.MethodGet, apiURL, nil)

			checkResponseCode(t, http.StatusOK, response.Code)

			var responseDTOs []companyDTO
			if err := json.Unmarshal(response.Body.Bytes(), &responseDTOs); err != nil {
				t.Errorf("unable to parse response: %v", err)
				return
			}

			if len(tt.expectedCompanies) != len(responseDTOs) {
				t.Errorf("expected count of companies %d, got %v", len(tt.expectedCompanies), len(responseDTOs))
				return
			}

			for index, responseDto := range responseDTOs {
				expectedData := tt.expectedCompanies[index]

				if expectedData.ID.String() != responseDto.ID {
					t.Errorf("expected id %v\nGot %v", expectedData.ID, responseDto.ID)
					return
				}

				if expectedData.Name != responseDto.Name {
					t.Errorf("expected name %v\nGot %v", expectedData.Name, responseDto.Name)
					return
				}

				if expectedData.Code != responseDto.Code {
					t.Errorf("expected code %v\nGot %v", expectedData.Code, responseDto.Code)
					return
				}

				if expectedData.Country != responseDto.Country {
					t.Errorf("expected country %v\nGot %v", expectedData.Country, responseDto.Country)
					return
				}

				if expectedData.Website != responseDto.Website {
					t.Errorf("expected website %v\nGot %v", expectedData.Website, responseDto.Website)
					return
				}

				if expectedData.Phone != responseDto.Phone {
					t.Errorf("expected phone %v\nGot %v", expectedData.Phone, responseDto.Phone)
					return
				}
			}
		})
	}
}
