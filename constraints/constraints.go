package constraints

import (
  "bytes"
  "errors"
  "fmt"
  "net/http"
  // "fmt"
  "regexp"
  "strconv"

  "github.com/manyminds/api2go"
)

type Constraints struct {
  Offset  uint64
  Limit   uint64
  Sort    string
}

func ApplyPaginatedConstraints(r api2go.Request) (Constraints, error) {
  var (
    finalOffset, finalLimit           uint64
    number, size, offset, limit, sort string
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
          return Constraints{}, api2go.NewHTTPError(errors.New("Invalid query parameter for 'sort'."), "Invalid query parameter for 'sort'.", http.StatusBadRequest)
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

  if size != "" {
    sizeI, err := strconv.ParseUint(size, 10, 64)
    if err != nil {
      return Constraints{}, err
    }

    numberI, err := strconv.ParseUint(number, 10, 64)
    if err != nil {
      return Constraints{}, err
    }

    finalOffset = sizeI * (numberI - 1)
    finalLimit = sizeI
  } else if limit != "" {
    limitI, err := strconv.ParseUint(limit, 10, 64)
    if err != nil {
      return Constraints{}, err
    }

    offsetI, err := strconv.ParseUint(offset, 10, 64)
    if err != nil {
      return Constraints{}, err
    }

    finalOffset = offsetI
    finalLimit = limitI
  }

  return Constraints{
    Offset: finalOffset, 
    Limit: finalLimit,
    Sort: sort}, nil
}