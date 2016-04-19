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

// NounResource for api2go routes
type NounResource struct {
	NounStorage *storage.NounStorage
}

// FindAll to satisfy api2go data source interface
func (s NounResource) FindAll(r api2go.Request) (api2go.Responder, error) {
	queryConstraints, err := constraints.ApplyPaginatedConstraints(r)
	if err != nil {
		return &Response{}, err
	}

	// Get the results
	nouns, err := s.NounStorage.GetAllPaginated(queryConstraints)
	if err != nil {
		return &Response{}, err
	}

	return &Response{Res: nouns}, nil
}

// PaginatedFindAll can be used to load nouns in chunks
func (s NounResource) PaginatedFindAll(r api2go.Request) (uint, api2go.Responder, error) {
	queryConstraints, err := constraints.ApplyPaginatedConstraints(r)
	if err != nil {
		return 0, &Response{}, err
	}

	// Get the results
	nouns, err := s.NounStorage.GetAllPaginated(queryConstraints)
	if err != nil {
		return 0, &Response{}, err
	}

	return uint(len(nouns)), &Response{Res: nouns}, nil
}

// FindOne to satisfy `api2go.DataSource` interface
// this method should return the noun with the given ID, otherwise an error
func (s NounResource) FindOne(ID string, r api2go.Request) (api2go.Responder, error) {
	noun, err := s.NounStorage.GetOne(ID)
	if err != nil {
		return &Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusNotFound)
	}
	return &Response{Res: noun}, nil
}

// Create method to satisfy `api2go.DataSource` interface
func (s NounResource) Create(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	noun, ok := obj.(model.Noun)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	id, _ := s.NounStorage.Insert(noun)
	noun.SetID(id)

	return &Response{Res: noun, Code: http.StatusCreated}, nil
}

// Delete to satisfy `api2go.DataSource` interface
func (s NounResource) Delete(id string, r api2go.Request) (api2go.Responder, error) {
	err := s.NounStorage.Delete(id)

	if err != nil {
		return &Response{Code: http.StatusNotFound}, api2go.NewHTTPError(err, err.Error(), http.StatusNotFound)
	}

	return &Response{Code: http.StatusNoContent}, err
}

// Update stores all changes on the noun
func (s NounResource) Update(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	noun, ok := obj.(model.Noun)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	err := s.NounStorage.Update(noun)
	return &Response{Res: noun, Code: http.StatusNoContent}, err
}
