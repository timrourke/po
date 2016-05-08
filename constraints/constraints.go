package constraints

import (
	"bytes"
	"errors"
	// "log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/manyminds/api2go"
)

type SingleConstraints struct {
	Includes map[string]struct{}
}

type PaginatedConstraints struct {
	Offset   uint64
	Limit    uint64
	Sort     string
	Includes map[string]struct{}
}

func ApplySingleConstraints(r api2go.Request) (SingleConstraints, error) {
	var includeStrings []string

	includeQuery, ok := r.QueryParams["include"]
	if ok {
		includeStrings = strings.Split(includeQuery[0], ",")
	} else {
		includeStrings = nil
	}
	includes := make(map[string]struct{}, 0)
	for _, include := range includeStrings {
		includes[include] = struct{}{}
	}

	return SingleConstraints{Includes: includes}, nil
}

func ApplyPaginatedConstraints(r api2go.Request) (PaginatedConstraints, error) {
	var (
		finalOffset, finalLimit           uint64
		number, size, offset, limit, sort string
		includeStrings                    []string
	)
	// Default constraints for queries
	finalOffset = 0
	finalLimit = 20
	sort = "id ASC"

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
	} else {
		offset = "0"
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
				return PaginatedConstraints{}, api2go.NewHTTPError(errors.New("Invalid query parameter for 'sort'."), "Invalid query parameter for 'sort'.", http.StatusBadRequest)
			}
			if string(sortQuery[i][0]) == "-" {
				buffer.WriteString(sortQuery[i][1:])
				buffer.WriteString(desc)
			} else {
				buffer.WriteString(sortQuery[i])
				buffer.WriteString(asc)
			}
			if i < len(sortQuery)-1 {
				buffer.WriteString(", ")
			}
		}
		sort = buffer.String()
	}
	includeQuery, ok := r.QueryParams["include"]
	if ok {
		includeStrings = strings.Split(includeQuery[0], ".")
	} else {
		includeStrings = nil
	}
	includes := make(map[string]struct{}, 0)
	for _, include := range includeStrings {
		includes[include] = struct{}{}
	}

	if size != "" {
		sizeI, err := strconv.ParseUint(size, 10, 64)
		if err != nil {
			return PaginatedConstraints{}, err
		}

		numberI, err := strconv.ParseUint(number, 10, 64)
		if err != nil {
			return PaginatedConstraints{}, err
		}

		finalOffset = sizeI * (numberI - 1)
		finalLimit = sizeI
	} else if limit != "" {
		limitI, err := strconv.ParseUint(limit, 10, 64)
		if err != nil {
			return PaginatedConstraints{}, err
		}

		offsetI, err := strconv.ParseUint(offset, 10, 64)
		if err != nil {
			return PaginatedConstraints{}, err
		}

		finalOffset = offsetI
		finalLimit = limitI
	}

	return PaginatedConstraints{
		Offset:   finalOffset,
		Limit:    finalLimit,
		Sort:     sort,
		Includes: includes}, nil
}
