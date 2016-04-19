package storage

import (
	"errors"
	"fmt"
	// "strconv"
	"net/http"
	// "regexp"
	"github.com/jmoiron/sqlx"
	"github.com/manyminds/api2go"
	"github.com/timrourke/po/constraints"
	"github.com/timrourke/po/model"
	"log"
)

// NewNounStorage initializes the storage
func NewNounStorage(db *sqlx.DB) *NounStorage {
	return &NounStorage{db}
}

// NounStorage stores all nouns
type NounStorage struct {
	db *sqlx.DB
}

// Get all
func (s NounStorage) GetAllPaginated(constraints constraints.Constraints) (model.ResultSet, error) {
	fmt.Println("we have a constraints struct", constraints)
	fmt.Sprintf("%+v", constraints)
	fmt.Println("Sort val should be: ", constraints.Sort)
	sqlString := fmt.Sprintf("SELECT * FROM noun ORDER BY %s LIMIT ?,?", constraints.Sort)
	fmt.Println("sqlString", sqlString)
	rows, err := s.db.Queryx(sqlString, constraints.Offset, constraints.Limit)
	defer rows.Close()

	if err != nil {
		log.Println(err)
		return nil, err
	}
	results := make([]model.Noun, 0)
	for rows.Next() {
		var n model.Noun
		err = rows.StructScan(&n)
		if err != nil {
			return nil, err
		}
		results = append(results, n)
	}

	var nounMap model.ResultSet
	for i := 0; i < len(results); i++ {
		n := results[i]
		nounMap = append(nounMap, &n)
	}
	return nounMap, nil
}

// Get one
func (s NounStorage) GetOne(id string) (model.Noun, error) {
	noun := model.Noun{}
	err := s.db.Get(&noun, `SELECT * FROM noun WHERE ID = ? LIMIT 1`, id)
	if err != nil {
		errMessage := fmt.Sprintf("Noun for id %s not found", id)
		return model.Noun{}, api2go.NewHTTPError(errors.New(errMessage), errMessage, http.StatusNotFound)
	}
	return noun, nil
}

// Insert
func (s *NounStorage) Insert(c model.Noun) (string, error) {
	noun := `INSERT INTO noun (singular, plural) VALUES (?, ?)`
	result, err := s.db.Exec(noun, c.Singular, c.Plural)
	insertId, err := result.LastInsertId()
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	c.SetID(fmt.Sprintf("%d", insertId))
	return c.GetID(), nil
}

// Delete
func (s *NounStorage) Delete(id string) error {
	delete := `DELETE FROM noun WHERE ID = ?`
	result, err := s.db.Exec(delete, id)
	numRowsDeleted, _ := result.RowsAffected()
	if err != nil {
		fmt.Println(err)
		return err
	}
	if numRowsDeleted == 0 {
		return fmt.Errorf("Noun with id %s does not exist", id)
	}
	return nil
}

// Update
func (s *NounStorage) Update(c model.Noun) error {
	_, err := s.db.NamedExec("UPDATE noun SET "+
		"singular=:singular, "+
		"plural=:plural "+
		"WHERE id = :id", c)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("Noun with id %s does not exist", c.ID)
	}

	return nil
}
