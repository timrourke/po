package storage

import (
	"errors"
	"fmt"
	// "strconv"
	"net/http"
	// "regexp"
	"github.com/jmoiron/sqlx"
	"github.com/manyminds/api2go"
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

// GetAll returns the noun map (because we need the ID as key too)
// func (s NounStorage) GetAll() (model.ResultSet, error) {
// 	rows, err := s.db.Queryx("SELECT * FROM noun ORDER BY id asc LIMIT 20")
// 	defer rows.Close()

// 	if err != nil {
// 		log.Println(err)
// 		return nil, err
// 	}
// 	results := make([]model.Noun, 0)
// 	for rows.Next() {
// 		var n model.Noun
// 		err = rows.StructScan(&n)
// 		log.Println(err)
// 		results = append(results, n)
// 	}

// 	nounMap := make(model.ResultSet)
// 	for i := range results {
// 		n := results[i]
// 		nounMap[n.ID] = &n
// 	}
// 	return nounMap, nil
// }

func (s NounStorage) GetAllPaginated(offset uint64, limit uint64, sort string) (model.ResultSet, error) {
	sqlString := fmt.Sprintf("SELECT * FROM noun ORDER BY %s LIMIT ?,?", sort)
	rows, err := s.db.Queryx(sqlString, offset, limit)
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
	fmt.Println("%+v", results)

	var nounMap model.ResultSet
	for i := 0; i < len(results); i++ {
		n := results[i]
		nounMap = append(nounMap, &n)
	}
	return nounMap, nil
}

// // GetOne noun
func (s NounStorage) GetOne(id string) (model.Noun, error) {
	noun := model.Noun{}
	err := s.db.Get(&noun, `SELECT * FROM noun WHERE ID = ? LIMIT 1`, id)
	if err != nil {
		errMessage := fmt.Sprintf("Noun for id %s not found", id)
		return model.Noun{}, api2go.NewHTTPError(errors.New(errMessage), errMessage, http.StatusNotFound)
	}
	return noun, nil
}

// Insert a noun
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

// // Delete one :(
func (s *NounStorage) Delete(id string) error {
	delete := `DELETE FROM noun WHERE ID = ?`
	result, err := s.db.Exec(delete, id)
	numRowsDeleted, _ := result.RowsAffected()
	if err != nil {
		return err
	}
	if numRowsDeleted == 0 {
		return fmt.Errorf("Noun with id %s does not exist", id)
	}
	return nil
}

// // Update a user
func (s *NounStorage) Update(c model.Noun) error {
	fmt.Printf("%v", c)
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
