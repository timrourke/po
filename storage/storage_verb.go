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

// NewVerbStorage initializes the storage
func NewVerbStorage(db *sqlx.DB) *VerbStorage {
	return &VerbStorage{db}
}

// VerbStorage stores all verbs
type VerbStorage struct {
	db *sqlx.DB
}

// Get all
func (s VerbStorage) GetAllPaginated(constraints constraints.Constraints) (model.ResultSet, error) {
	sqlString := fmt.Sprintf("SELECT * FROM verb ORDER BY %s LIMIT ?,?", constraints.Sort)
	rows, err := s.db.Queryx(sqlString, constraints.Offset, constraints.Limit)
	defer rows.Close()

	if err != nil {
		log.Println(err)
		return nil, err
	}
	results := make([]model.Verb, 0)
	for rows.Next() {
		var n model.Verb
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
	err := s.db.Get(&verb, `SELECT * FROM verb WHERE ID = ? LIMIT 1`, id)
	if err != nil {
		errMessage := fmt.Sprintf("Noun for id %s not found", id)
		return model.Verb{}, api2go.NewHTTPError(errors.New(errMessage), errMessage, http.StatusNotFound)
	}
	return verb, nil
}

// Insert
func (s *VerbStorage) Insert(c model.Verb) (string, error) {
	verb := `INSERT INTO verb (aux_verb_id, gerund, infinitive, past_participle, reflexive) VALUES (?, ?, ?, ?, ?)`
	result, err := s.db.Exec(verb, c.AuxVerbId, c.Gerund, c.Infinitive, c.PastParticiple, c.Reflexive)
	insertId, err := result.LastInsertId()
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	c.SetID(fmt.Sprintf("%d", insertId))
	return c.GetID(), nil
}

// Delete
func (s *VerbStorage) Delete(id string) error {
	delete := `DELETE FROM verb WHERE ID = ?`
	result, err := s.db.Exec(delete, id)
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
	_, err := s.db.NamedExec("UPDATE verb SET "+
		"aux_verb_id=:aux_verb_id, "+
		"gerund=:gerund, "+
		"infinitive=:infinitive, "+
		"past_participle=:past_participle, "+
		"reflexive=:reflexive "+
		"WHERE id = :id", c)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("Noun with id %s does not exist", c.ID)
	}

	return nil
}
