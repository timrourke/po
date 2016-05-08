package storage

import (
	"errors"
	"fmt"
	"github.com/manyminds/api2go"
	"github.com/timrourke/po/constraints"
	"github.com/timrourke/po/database"
	"github.com/timrourke/po/model"
	"github.com/timrourke/sqlx-helpers"
	"log"
	"net/http"
)

// VerbStorage stores all verbs
type VerbStorage struct{}

// Get all
func (s VerbStorage) GetAllPaginated(constraints constraints.PaginatedConstraints) (model.ResultSet, error) {
	sqlString := fmt.Sprintf("SELECT * FROM verb ORDER BY %s LIMIT ?,?", constraints.Sort)
	rows, err := database.DB.Queryx(sqlString, constraints.Offset, constraints.Limit)
	defer rows.Close()

	if err != nil {
		log.Println(err)
		return nil, err
	}

	results := make([]model.Verb, 0)
	for rows.Next() {
		n := model.Verb{}
		err = rows.StructScan(&n)
		if err != nil {
			return nil, err
		}
		results = append(results, n)
	}

	var verbMap model.ResultSet
	for i := 0; i < len(results); i++ {
		n := results[i]
		verbMap = append(verbMap, &n)
	}
	return verbMap, nil
}

// Get one
func (s VerbStorage) GetOne(id string) (model.Verb, error) {
	verb := model.Verb{}
	err := database.DB.Get(&verb, `SELECT * FROM verb WHERE id = ? LIMIT 1`, id)
	if err != nil {
		errMessage := fmt.Sprintf("Verb for id %s not found", id)
		return model.Verb{}, api2go.NewHTTPError(errors.New(errMessage), errMessage, http.StatusNotFound)
	}
	return verb, nil
}

// Insert
func (s *VerbStorage) Insert(c model.Verb) (string, error) {
	sql, values := helper.CreateInsert(c, c.TableName())
	result, err := database.DB.Exec(sql, values...)

	if err != nil {
		log.Println(err)
		return "", err
	}

	insertId, err := result.LastInsertId()
	if err != nil {
		log.Println(err)
		return "", err
	}

	c.SetID(fmt.Sprintf("%d", insertId))
	return c.GetID(), nil
}

// Delete
func (s *VerbStorage) Delete(id string) error {
	delete := `DELETE FROM verb WHERE id = ?`
	result, err := database.DB.Exec(delete, id)
	numRowsDeleted, _ := result.RowsAffected()

	if err != nil {
		fmt.Println(err)
		return err
	}

	if numRowsDeleted == 0 {
		return fmt.Errorf("Verb with id %s does not exist", id)
	}

	return nil
}

// Update
func (s *VerbStorage) Update(c model.Verb) error {
	sql, values := helper.CreateUpdateOne(c, c.TableName(), c.GetID())
	_, err := database.DB.Exec(sql, values...)

	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("Verb with id %s does not exist", c.ID)
	}

	return nil
}
