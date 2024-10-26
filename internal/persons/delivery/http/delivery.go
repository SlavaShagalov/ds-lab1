package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	pPersons "github.com/SlavaShagalov/ds-lab1/internal/persons"
	"github.com/SlavaShagalov/ds-lab1/internal/pkg/constants"
	pErrors "github.com/SlavaShagalov/ds-lab1/internal/pkg/errors"
	pHTTP "github.com/SlavaShagalov/ds-lab1/internal/pkg/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

const (
	personsPrefix = "/persons"

	personsPath = constants.ApiPrefix + personsPrefix
	personPath  = personsPath + "/{id}"
)

type delivery struct {
	repo pPersons.Repository
	log  *zap.Logger
}

func RegisterHandlers(mux *mux.Router, repo pPersons.Repository, log *zap.Logger) {
	del := delivery{
		repo: repo,
		log:  log,
	}

	mux.HandleFunc(personsPath, del.create).Methods(http.MethodPost)
	mux.HandleFunc(personPath, del.get).Methods(http.MethodGet)
	mux.HandleFunc(personsPath, del.list).Methods(http.MethodGet)
	mux.HandleFunc(personPath, del.partialUpdate).Methods(http.MethodPatch)
	mux.HandleFunc(personPath, del.delete).Methods(http.MethodDelete)
}

// create godoc
//
//	@Summary		Create a new person
//	@Description	Create a new person
//	@Tags			workspaces
//	@Accept			json
//	@Produce		json
//	@Param			id				path		int				true	"Workspace ID"
//	@Param			PersonCreateData	body		createRequest	true	"Person create data"
//	@Success		200				{object}	createResponse	"Created person data."
//	@Failure		400				{object}	http.JSONError
//	@Failure		401				{object}	http.JSONError
//	@Failure		405
//	@Failure		500
//	@Router			/workspaces/{id}/persons [post]
func (del *delivery) create(w http.ResponseWriter, r *http.Request) {
	body, err := pHTTP.ReadBody(r, del.log)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	var request createRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		pHTTP.HandleError(w, r, pErrors.ErrReadBody)
		return
	}

	params := pPersons.CreateParams{
		Name:    request.Name,
		Age:     request.Age,
		Address: request.Address,
		Work:    request.Work,
	}

	person, err := del.repo.Create(r.Context(), &params)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	w.Header().Add("Location", fmt.Sprintf(personsPath+"/%d", person.ID))
	w.WriteHeader(http.StatusCreated)
}

// get godoc
//
//	@Summary		Returns person by id
//	@Description	Returns person by id
//	@Tags			persons
//	@Produce		json
//	@Param			id	path		int			true	"Person ID"
//	@Success		200	{object}	getResponse	"Person data"
//	@Failure		400	{object}	http.JSONError
//	@Failure		401	{object}	http.JSONError
//	@Failure		404	{object}	http.JSONError
//	@Failure		405
//	@Failure		500
//	@Router			/persons/{id} [get]
func (del *delivery) get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	personID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	person, err := del.repo.Get(r.Context(), personID)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	response := newGetResponse(person)
	pHTTP.SendJSON(w, r, http.StatusOK, response)
}

// list godoc
//
//	@Summary		Returns persons by workspace id
//	@Description	Returns persons by workspace id
//	@Tags			persons
//	@Produce		json
//	@Param			title	query		string			true	"Title filter"
//	@Success		200		{object}	listResponse	"Persons data"
//	@Failure		400		{object}	http.JSONError
//	@Failure		401		{object}	http.JSONError
//	@Failure		405
//	@Failure		500
//	@Router			/persons [get]
//
//	@Security		cookieAuth
func (del *delivery) list(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	var err error
	var limit int64 = 0
	if queryParams.Get("limit") != "" {
		limit, err = strconv.ParseInt(queryParams.Get("limit"), 10, 64)
		if err != nil {
			pHTTP.HandleError(w, r, err)
			return
		}
	}
	var offset int64 = 0
	if queryParams.Get("offset") != "" {
		offset, err = strconv.ParseInt(queryParams.Get("offset"), 10, 64)
		if err != nil {
			pHTTP.HandleError(w, r, err)
			return
		}
	}

	persons, err := del.repo.List(r.Context(), offset, limit)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	//response := newListResponse(persons)
	pHTTP.SendJSON(w, r, http.StatusOK, persons)
}

// partialUpdate godoc
//
//	@Summary		Partial update of person
//	@Description	Partial update of person
//	@Tags			persons
//	@Accept			json
//	@Produce		json
//	@Param			id				path		int						true	"Person ID"
//	@Param			PersonUpdateData	body		partialUpdateRequest	true	"Person data to update"
//	@Success		200				{object}	getResponse				"Updated person data."
//	@Failure		400				{object}	http.JSONError
//	@Failure		401				{object}	http.JSONError
//	@Failure		405
//	@Failure		500
//	@Router			/persons/{id}  [patch]
//
//	@Security		cookieAuth
func (del *delivery) partialUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	personID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	body, err := pHTTP.ReadBody(r, del.log)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	var request partialUpdateRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		pHTTP.HandleError(w, r, pErrors.ErrReadBody)
		return
	}

	params := pPersons.PartialUpdateParams{
		ID:      personID,
		Name:    request.Name,
		Age:     request.Age,
		Address: request.Address,
		Work:    request.Work,
	}

	person, err := del.repo.PartialUpdate(r.Context(), &params)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	response := newGetResponse(person)
	pHTTP.SendJSON(w, r, http.StatusOK, response)
}

// delete godoc
//
//	@Summary		Delete person by id
//	@Description	Delete person by id
//	@Tags			persons
//	@Produce		json
//	@Param			id	path	int	true	"Person ID"
//	@Success		204	"Person deleted successfully"
//	@Failure		400	{object}	http.JSONError
//	@Failure		401	{object}	http.JSONError
//	@Failure		404	{object}	http.JSONError
//	@Failure		405
//	@Failure		500
//	@Router			/persons/{id} [delete]
//
//	@Security		cookieAuth
func (del *delivery) delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	personID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	err = del.repo.Delete(r.Context(), personID)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
