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

// NewTensePresIndStorage initializes the storage
func NewTensePresIndStorage(db *sqlx.DB) *TensePresIndStorage {
	return &TensePresIndStorage{db}
}

// TensePresIndStorage stores all tenses
type TensePresIndStorage struct {
	db *sqlx.DB
}

// Get all
func (t TensePresIndStorage) GetAllPaginated(constraints constraints.PaginatedConstraints) (model.ResultSet, error) {
	sqlString := fmt.Sprintf("SELECT * FROM tense_pres_ind ORDER BY %s LIMIT ?,?", constraints.Sort)
	rows, err := t.db.Queryx(sqlString, constraints.Offset, constraints.Limit)
	defer rows.Close()

	if err != nil {
		log.Println(err)
		return nil, err
	}
	results := make([]model.TensePresentIndicative, 0)
	for rows.Next() {
		n := model.NewTensePresentIndicative(t.db, constraints.Includes)
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
	tense := model.NewTensePresentIndicative(t.db, constraints.Includes)
	err := t.db.Get(&tense, `SELECT * FROM tense_pres_ind WHERE ID = ? LIMIT 1`, id)
	if err != nil {
		errMessage := fmt.Sprintf("Tense Present Indicative for id %s not found", id)
		return model.TensePresentIndicative{}, api2go.NewHTTPError(errors.New(errMessage), errMessage, http.StatusNotFound)
	}
	return tense, nil
}

// Insert
func (t *TensePresIndStorage) Insert(c model.TensePresentIndicative) (string, error) {
	tense := `INSERT INTO tense_pres_ind (created_at, verb_id, sing_first, sing_second, sing_third, plural_first, plural_second, plural_third) VALUES (NOW(), ?, ?, ?, ?, ?, ?, ?)`
	result, err := t.db.Exec(tense, c.VerbId, c.FirstPersonSingular, c.SecondPersonSingular, c.ThirdPersonSingular, c.FirstPersonPlural, c.SecondPersonPlural, c.ThirdPersonPlural)
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
	result, err := t.db.Exec(delete, id)
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
	_, err := t.db.NamedExec("UPDATE tense_pres_ind SET "+
		"updated_at=NOW(), "+
		"verb_id=:verb_id, "+
		"sing_first=:sing_first, "+
		"sing_second=:sing_second, "+
		"sing_third=:sing_third, "+
		"plural_first=:plural_first, "+
		"plural_second=:plural_second, "+
		"plural_third=:plural_third, "+
		"WHERE id = :id", c)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("Tense Present Indicative with id %s does not exist", c.ID)
	}

	return nil
}
