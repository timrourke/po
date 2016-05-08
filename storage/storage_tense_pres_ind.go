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

// TensePresIndStorage stores all tenses
type TensePresIndStorage struct{}

// Get all
func (t TensePresIndStorage) GetAllPaginated(constraints constraints.PaginatedConstraints) (model.ResultSet, error) {
	sqlString := fmt.Sprintf("SELECT * FROM tense_pres_ind ORDER BY %s LIMIT ?,?", constraints.Sort)
	rows, err := database.DB.Queryx(sqlString, constraints.Offset, constraints.Limit)
	defer rows.Close()

	if err != nil {
		log.Println(err)
		return nil, err
	}

	results := make([]model.TensePresentIndicative, 0)
	for rows.Next() {
		n := model.NewTensePresentIndicative(constraints.Includes)
		err = rows.StructScan(&n)
		if err != nil {
			return nil, err
		}
		results = append(results, n)
	}

	var tenseMap model.ResultSet
	for i := 0; i < len(results); i++ {
		n := results[i]
		tenseMap = append(tenseMap, &n)
	}
	return tenseMap, nil
}

// Get one
func (t TensePresIndStorage) GetOne(id string, constraints constraints.SingleConstraints) (model.TensePresentIndicative, error) {
	tense := model.NewTensePresentIndicative(constraints.Includes)
	err := database.DB.Get(&tense, `SELECT * FROM tense_pres_ind WHERE ID = ? LIMIT 1`, id)
	if err != nil {
		errMessage := fmt.Sprintf("Tense Present Indicative for id %s not found", id)
		return model.TensePresentIndicative{}, api2go.NewHTTPError(errors.New(errMessage), errMessage, http.StatusNotFound)
	}
	return tense, nil
}

// Insert
func (t *TensePresIndStorage) Insert(c model.TensePresentIndicative) (string, error) {
	sql, values := helper.CreateInsert(c, c.TableName())
	result, err := database.DB.Exec(sql, values...)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	insertId, err := result.LastInsertId()
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	c.SetID(fmt.Sprintf("%d", insertId))
	return c.GetID(), nil
}

// Delete
func (t *TensePresIndStorage) Delete(id string) error {
	delete := `DELETE FROM tense_pres_ind WHERE ID = ?`
	result, err := database.DB.Exec(delete, id)
	numRowsDeleted, _ := result.RowsAffected()
	if err != nil {
		fmt.Println(err)
		return err
	}
	if numRowsDeleted == 0 {
		return fmt.Errorf("Tense Present Indicative with id %s does not exist", id)
	}
	return nil
}

// Update
func (t *TensePresIndStorage) Update(c model.TensePresentIndicative) error {
	sql, values := helper.CreateUpdateOne(c, c.TableName(), c.GetID())
	_, err := database.DB.Exec(sql, values...)

	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("Tense Present Indicative with id %s does not exist", c.ID)
	}

	return nil
}
