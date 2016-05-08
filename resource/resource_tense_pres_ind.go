package resource

import (
	// "bytes"
	"errors"
	"log"
	// "fmt"
	"net/http"
	// "regexp"
	// "strconv"

	"github.com/manyminds/api2go"
	"github.com/timrourke/po/constraints"
	"github.com/timrourke/po/model"
	"github.com/timrourke/po/storage"
	"github.com/timrourke/validator"
)

// TensePresIndResource for api2go routes
type TensePresIndResource struct {
	TensePresIndStorage storage.TensePresIndStorage
}

// FindAll to satisfy api2go data source interface
func (s TensePresIndResource) FindAll(r api2go.Request) (api2go.Responder, error) {
	queryConstraints, err := constraints.ApplyPaginatedConstraints(r)
	if err != nil {
		return &Response{}, err
	}

	// Get the results
	tenses, err := s.TensePresIndStorage.GetAllPaginated(queryConstraints)
	if err != nil {
		return &Response{}, err
	}

	return &Response{Res: tenses}, nil
}

// PaginatedFindAll
func (s TensePresIndResource) PaginatedFindAll(r api2go.Request) (uint, api2go.Responder, error) {
	queryConstraints, err := constraints.ApplyPaginatedConstraints(r)
	if err != nil {
		return 0, &Response{}, err
	}

	// Get the results
	tenses, err := s.TensePresIndStorage.GetAllPaginated(queryConstraints)
	if err != nil {
		return 0, &Response{}, err
	}

	return uint(len(tenses)), &Response{Res: tenses}, nil
}

// FindOne to satisfy `api2go.DataSource` interface
// this method should return the tense with the given ID, otherwise an error
func (s TensePresIndResource) FindOne(ID string, r api2go.Request) (api2go.Responder, error) {
	queryConstraints, err := constraints.ApplySingleConstraints(r)
	if err != nil {
		return &Response{}, err
	}

	tense, err := s.TensePresIndStorage.GetOne(ID, queryConstraints)
	if err != nil {
		return &Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusNotFound)
	}
	return &Response{Res: tense}, nil
}

// Create method to satisfy `api2go.DataSource` interface
func (s TensePresIndResource) Create(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	tense, ok := obj.(model.TensePresentIndicative)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}
	err := validator.Validate(tense)
	log.Println(err)

	id, _ := s.TensePresIndStorage.Insert(tense)
	tense.SetID(id)

	return &Response{Res: tense, Code: http.StatusCreated}, nil
}

// Delete to satisfy `api2go.DataSource` interface
func (s TensePresIndResource) Delete(id string, r api2go.Request) (api2go.Responder, error) {
	err := s.TensePresIndStorage.Delete(id)

	if err != nil {
		return &Response{Code: http.StatusNotFound}, api2go.NewHTTPError(err, err.Error(), http.StatusNotFound)
	}

	return &Response{Code: http.StatusNoContent}, err
}

// Update stores all changes on the tense
func (s TensePresIndResource) Update(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	tense, ok := obj.(model.TensePresentIndicative)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	err := s.TensePresIndStorage.Update(tense)
	return &Response{Res: tense, Code: http.StatusNoContent}, err
}
