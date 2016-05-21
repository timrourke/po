package storage

import (
	"errors"
	"fmt"
	"github.com/manyminds/api2go"
	"github.com/timrourke/po/constraints"
	"github.com/timrourke/po/database"
	"github.com/timrourke/po/model"
	"log"
	"net/http"
)

// NounStorage stores all nouns
type NounStorage struct{}

// Get all
func (s NounStorage) GetAllPaginated(constraints constraints.PaginatedConstraints) (model.ResultSet, error) {
	results := make([]model.Noun, 0)
	sqlString := fmt.Sprintf("SELECT * FROM noun ORDER BY %s LIMIT ?,?", constraints.Sort)
	err := database.DB.Select(&results, sqlString, constraints.Offset, constraints.Limit)

	if err != nil {
		log.Println(err)
		return nil, err
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
	err := database.DB.Get(&noun, `SELECT * FROM noun WHERE ID = ? LIMIT 1`, id)
	if err != nil {
		errMessage := fmt.Sprintf("Noun for id %s not found", id)
		return model.Noun{}, api2go.NewHTTPError(errors.New(errMessage), errMessage, http.StatusNotFound)
	}
	return noun, nil
}

// Insert
func (s *NounStorage) Insert(c model.Noun) (string, error) {
	noun := `INSERT INTO noun (singular, plural) VALUES (?, ?)`
	result, err := database.DB.Exec(noun, c.Singular, c.Plural)
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
	result, err := database.DB.Exec(delete, id)
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
	_, err := database.DB.NamedExec("UPDATE noun SET "+
		"singular=:singular, "+
		"plural=:plural "+
		"WHERE id = :id", c)

	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("Noun with id %s does not exist", c.ID)
	}

	return nil
}
