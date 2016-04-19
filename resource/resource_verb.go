package resource

import (
	// "bytes"
	"errors"
	// "fmt"
	"net/http"
	// "regexp"
	// "strconv"

	"github.com/manyminds/api2go"
	"github.com/timrourke/po/constraints"
	"github.com/timrourke/po/model"
	"github.com/timrourke/po/storage"
)

// VerbResource for api2go routes
type VerbResource struct {
	VerbStorage *storage.VerbStorage
}

// FindAll to satisfy api2go data source interface
func (s VerbResource) FindAll(r api2go.Request) (api2go.Responder, error) {
	queryConstraints, err := constraints.ApplyPaginatedConstraints(r)
	if err != nil {
		return &Response{}, err
	}

	// Get the results
	verbs, err := s.VerbStorage.GetAllPaginated(queryConstraints)
	if err != nil {
		return &Response{}, err
	}

	return &Response{Res: verbs}, nil
}

// PaginatedFindAll
func (s VerbResource) PaginatedFindAll(r api2go.Request) (uint, api2go.Responder, error) {
	queryConstraints, err := constraints.ApplyPaginatedConstraints(r)
	if err != nil {
		return 0, &Response{}, err
	}

	// Get the results
	verbs, err := s.VerbStorage.GetAllPaginated(queryConstraints)
	if err != nil {
		return 0, &Response{}, err
	}

	return uint(len(verbs)), &Response{Res: verbs}, nil
}

// FindOne to satisfy `api2go.DataSource` interface
// this method should return the verb with the given ID, otherwise an error
func (s VerbResource) FindOne(ID string, r api2go.Request) (api2go.Responder, error) {
	verb, err := s.VerbStorage.GetOne(ID)
	if err != nil {
		return &Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusNotFound)
	}
	return &Response{Res: verb}, nil
}

// Create method to satisfy `api2go.DataSource` interface
func (s VerbResource) Create(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	verb, ok := obj.(model.Verb)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	id, _ := s.VerbStorage.Insert(verb)
	verb.SetID(id)

	return &Response{Res: verb, Code: http.StatusCreated}, nil
}

// Delete to satisfy `api2go.DataSource` interface
func (s VerbResource) Delete(id string, r api2go.Request) (api2go.Responder, error) {
	err := s.VerbStorage.Delete(id)

	if err != nil {
		return &Response{Code: http.StatusNotFound}, api2go.NewHTTPError(err, err.Error(), http.StatusNotFound)
	}

	return &Response{Code: http.StatusNoContent}, err
}

// Update stores all changes on the verb
func (s VerbResource) Update(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	verb, ok := obj.(model.Verb)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	err := s.VerbStorage.Update(verb)
	return &Response{Res: verb, Code: http.StatusNoContent}, err
}
