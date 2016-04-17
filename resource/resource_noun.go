package resource

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/manyminds/api2go"
	"github.com/timrourke/po/model"
	"github.com/timrourke/po/storage"
)

// UserResource for api2go routes
type NounResource struct {
	NounStorage *storage.NounStorage
}

// FindAll to satisfy api2go data source interface
func (s NounResource) FindAll(r api2go.Request) (api2go.Responder, error) {
	var (
		limit 			string
		finalLimit	uint64
		sort 				string
	)

	// Default limit for queries
	finalLimit = 20

	// Apply a query limit if supplied in request
	limitQuery, ok := r.QueryParams["page[limit]"]
	if ok {
		limit = limitQuery[0]
	}
	sortQuery, ok := r.QueryParams["sort"]
	if !ok {
		sort = "id ASC"
	} else {
		var buffer bytes.Buffer
		asc := " ASC"
		desc := " DESC"
		for i := 0; i < len(sortQuery); i++ {
			valid := regexp.MustCompile("^-?[A-Za-z0-9_.]+$")
			if !valid.MatchString(sortQuery[i]) {
			    // invalid column name, do not proceed in order to prevent SQL injection
					return &Response{}, api2go.NewHTTPError(errors.New("Invalid query parameter for 'sort'."), "Invalid query parameter for 'sort'.", http.StatusBadRequest)
			}
			if string(sortQuery[i][0]) == "-" {
				buffer.WriteString(sortQuery[i][1:])
				buffer.WriteString(desc)
			} else {
				buffer.WriteString(sortQuery[i])
				buffer.WriteString(asc)
			}

			// Add trailing comma
			if i < len(sortQuery) - 1 {
				buffer.WriteString(", ")
			}
		}
		sort = buffer.String()
	}
	fmt.Println("sort", sort)

	if limit != "" {
		var err error
		finalLimit, err = strconv.ParseUint(limit, 10, 64)
		if err != nil {
			return &Response{}, err
		}
	}

	// Get the results
	nouns, err := s.NounStorage.GetAllPaginated(0, finalLimit, sort)
	if err != nil {
		return &Response{}, err
	}

	return &Response{Res: nouns}, nil
}

// PaginatedFindAll can be used to load users in chunks
func (s NounResource) PaginatedFindAll(r api2go.Request) (uint, api2go.Responder, error) {
	var (
		finalOffset, finalLimit 					uint64
		number, size, offset, limit, sort string
	)
	
	// Default limit for queries
	finalLimit = 20

	// Apply query constraints if supplied in request
	numberQuery, ok := r.QueryParams["page[number]"]
	if ok {
		number = numberQuery[0]
	}
	sizeQuery, ok := r.QueryParams["page[size]"]
	if ok {
		size = sizeQuery[0]
	}
	offsetQuery, ok := r.QueryParams["page[offset]"]
	if ok {
		offset = offsetQuery[0]
	}
	limitQuery, ok := r.QueryParams["page[limit]"]
	if ok {
		limit = limitQuery[0]
	}
	sortQuery, ok := r.QueryParams["sort"]
	if !ok {
		sort = "id ASC"
	} else {
		var buffer bytes.Buffer
		asc := " ASC"
		desc := " DESC"
		for i := 0; i < len(sortQuery); i++ {
			valid := regexp.MustCompile("^-?[A-Za-z0-9_.]+$")
			if !valid.MatchString(sortQuery[i]) {
			    // invalid column name, do not proceed in order to prevent SQL injection
					return 0, &Response{}, api2go.NewHTTPError(errors.New("Invalid query parameter for 'sort'."), "Invalid query parameter for 'sort'.", http.StatusBadRequest)
			}
			if string(sortQuery[i][0]) == "-" {
				buffer.WriteString(sortQuery[i][1:])
				buffer.WriteString(desc)
			} else {
				buffer.WriteString(sortQuery[i])
				buffer.WriteString(asc)
			}
			if i < len(sortQuery) - 1 {
				buffer.WriteString(", ")
			}
		}
		sort = buffer.String()
	}
	fmt.Println("sort", sort)

	if size != "" {
		sizeI, err := strconv.ParseUint(size, 10, 64)
		if err != nil {
			return 0, &Response{}, err
		}

		numberI, err := strconv.ParseUint(number, 10, 64)
		if err != nil {
			return 0, &Response{}, err
		}

		finalOffset = sizeI * (numberI - 1)
		finalLimit = sizeI
	} else {
		limitI, err := strconv.ParseUint(limit, 10, 64)
		if err != nil {
			return 0, &Response{}, err
		}

		offsetI, err := strconv.ParseUint(offset, 10, 64)
		if err != nil {
			return 0, &Response{}, err
		}

		finalOffset = offsetI
		finalLimit = limitI
	}

	// Get the results
	nouns, err := s.NounStorage.GetAllPaginated(finalOffset, finalLimit, sort)
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

// Update stores all changes on the user
func (s NounResource) Update(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	noun, ok := obj.(model.Noun)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	err := s.NounStorage.Update(noun)
	return &Response{Res: noun, Code: http.StatusNoContent}, err
}
