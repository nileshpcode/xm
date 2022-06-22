package controller

import (
	"encoding/json"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"net/http"
	"xm/app"
	"xm/client"
	apiError "xm/error"
	"xm/model"
	"xm/repository"
)

type companyController struct {
	app              *app.App
	ipLocationClient client.IPLocationClient
	repository       repository.Repository
}

func NewCompanyController(app *app.App, ipLocationClient client.IPLocationClient, repository repository.Repository) *companyController {
	return &companyController{
		app:              app,
		ipLocationClient: ipLocationClient,
		repository:       repository,
	}
}

// RegisterRoutes implements interface RouteSpecifier
func (controller *companyController) RegisterRoutes(muxRouter *mux.Router) {
	router := muxRouter.PathPrefix("/api/companies").Subrouter()

	router.HandleFunc("", protect(controller.ipLocationClient, controller.add)).Methods(http.MethodPost)
	router.HandleFunc("", controller.getAll).Methods(http.MethodGet)
	router.HandleFunc("/{id}", controller.get).Methods(http.MethodGet)
	router.HandleFunc("/{id}", controller.update).Methods(http.MethodPut)
	router.HandleFunc("/{id}", protect(controller.ipLocationClient, controller.delete)).Methods(http.MethodDelete)
}

func (controller *companyController) add(w http.ResponseWriter, r *http.Request) {
	uow := repository.NewUnitOfWork(controller.app.DB, false)
	defer uow.Complete()

	reqDTO := companyDTO{}
	if err := unmarshalJSON(r, &reqDTO); err != nil {
		controller.app.Logger.Err(err).Msg("unable to marshal request body")
		respondError(w, err)
		return
	}

	company, err := model.NewCompany(reqDTO.Name, reqDTO.Code, reqDTO.Country, reqDTO.Website, reqDTO.Phone)
	if err != nil {
		controller.app.Logger.Err(err).Msg("unable add company")
		respondError(w, err)
		return
	}

	if err := controller.repository.Add(uow, company); err != nil {
		controller.app.Logger.Err(err).Msg("unable add company to db")
		respondError(w, err)
		return
	}

	uow.Commit()

	respondJSON(w, http.StatusCreated, toCompanyDTO(company))
	return
}

func (controller *companyController) getAll(w http.ResponseWriter, r *http.Request) {
	uow := repository.NewUnitOfWork(controller.app.DB, true)
	defer uow.Complete()

	var queryProcessors []repository.QueryProcessor

	if name := r.FormValue("name"); len(name) > 0 {
		queryProcessors = append(queryProcessors, repository.Filter("name = ?", name))
	}

	if code := r.FormValue("code"); len(code) > 0 {
		queryProcessors = append(queryProcessors, repository.Filter("code = ?", code))
	}

	if country := r.FormValue("country"); len(country) > 0 {
		queryProcessors = append(queryProcessors, repository.Filter("country = ?", country))
	}

	if website := r.FormValue("website"); len(website) > 0 {
		queryProcessors = append(queryProcessors, repository.Filter("website = ?", website))
	}

	if phone := r.FormValue("phone"); len(phone) > 0 {
		queryProcessors = append(queryProcessors, repository.Filter("phone = ?", phone))
	}

	var companies []model.Company
	if err := controller.repository.GetAll(uow, &companies, queryProcessors); err != nil {
		controller.app.Logger.Err(err).Msg("unable to get companies from db")
		respondError(w, err)
		return
	}

	responseDTO := make([]companyDTO, len(companies))
	for index, company := range companies {
		responseDTO[index] = toCompanyDTO(&company)
	}

	respondJSON(w, http.StatusOK, responseDTO)
	return
}

func (controller *companyController) get(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	uow := repository.NewUnitOfWork(controller.app.DB, true)
	defer uow.Complete()

	company := &model.Company{}
	if err := controller.repository.Get(uow, company, uuid.FromStringOrNil(id)); err != nil {
		if err.IsRecordNotFoundError() {
			respondJSON(w, http.StatusNotFound, nil)
			return
		}
		controller.app.Logger.Err(err).Msg("unable get company from db")
		respondError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, toCompanyDTO(company))
	return
}

func (controller *companyController) update(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	uow := repository.NewUnitOfWork(controller.app.DB, false)
	defer uow.Complete()

	company := &model.Company{}
	if err := controller.repository.Get(uow, company, uuid.FromStringOrNil(id)); err != nil {
		if err.IsRecordNotFoundError() {
			respondJSON(w, http.StatusNotFound, nil)
			return
		}
		controller.app.Logger.Err(err).Msg("unable get company from db")
		respondError(w, err)
		return
	}

	reqDTO := companyDTO{}
	if err := unmarshalJSON(r, &reqDTO); err != nil {
		controller.app.Logger.Err(err).Msg("unable to marshal request body")
		respondError(w, err)
		return
	}

	err := company.Update(reqDTO.Name, reqDTO.Code, reqDTO.Country, reqDTO.Website, reqDTO.Phone)
	if err != nil {
		controller.app.Logger.Err(err).Msg("unable update company")
		respondError(w, err)
		return
	}

	if err := controller.repository.Update(uow, company); err != nil {
		controller.app.Logger.Err(err).Msg("unable update company to db")
		respondError(w, err)
		return
	}

	uow.Commit()

	respondJSON(w, http.StatusOK, toCompanyDTO(company))
	return
}

func (controller *companyController) delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	uow := repository.NewUnitOfWork(controller.app.DB, false)
	defer uow.Complete()

	company := &model.Company{}
	if err := controller.repository.Get(uow, company, uuid.FromStringOrNil(id)); err != nil {
		if err.IsRecordNotFoundError() {
			respondJSON(w, http.StatusNotFound, nil)
			return
		}
		controller.app.Logger.Err(err).Msg("unable get company from db")
		respondError(w, err)
		return
	}

	if err := controller.repository.Delete(uow, company); err != nil {
		controller.app.Logger.Err(err).Msg("unable delete company from db")
		respondError(w, err)
		return
	}

	uow.Commit()

	respondJSON(w, http.StatusOK, nil)
	return
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type companyDTO struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Code    string `json:"code"`
	Country string `json:"country"`
	Website string `json:"website"`
	Phone   string `json:"phone"`
}

func toCompanyDTO(company *model.Company) companyDTO {
	dto := companyDTO{
		ID:      company.ID.String(),
		Name:    company.Name,
		Code:    company.Code,
		Country: company.Country,
		Website: company.Website,
		Phone:   company.Phone,
	}
	return dto
}

// unmarshalJSON checks for empty body and then parses JSON into the target
func unmarshalJSON(r *http.Request, target interface{}) error {
	if r.Body == nil {
		return apiError.NewInvalidRequestPayloadError(apiError.ErrorCodeEmptyRequestBody)
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return apiError.NewDataReadWriteError(err)
	}

	if len(body) == 0 {
		return apiError.NewInvalidRequestPayloadError(apiError.ErrorCodeEmptyRequestBody)
	}

	err = json.Unmarshal(body, target)
	if err != nil {
		return apiError.NewInvalidRequestPayloadError(apiError.ErrorCodeInvalidJSON)
	}
	return nil
}

// respondJSON makes the response with payload as json format
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

// respondError returns a validation error else
func respondError(w http.ResponseWriter, err error) {
	switch err.(type) {
	case apiError.ValidationError:
		respondJSON(w, http.StatusBadRequest, err)
	default:
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": apiError.ErrorCodeInternalError})
	}
}
