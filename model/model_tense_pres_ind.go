package model

import (
	// "github.com/timrourke/po/storage"
	"github.com/jmoiron/sqlx"
	"github.com/manyminds/api2go/jsonapi"
	"gopkg.in/guregu/null.v3"
	"log"
	"strconv"
)

/* MySQL Schema

CREATE TABLE `tense_pres_ind` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `verb_id` int(11) NOT NULL,
  `sing_first` varchar(255),
  `sing_second` varchar(255),
  `sing_third` varchar(255),
  `plural_first` varchar(255),
  `plural_second` varchar(255),
  `plural_third` varchar(255),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8;

*/

func NewTensePresentIndicative(db *sqlx.DB, includes map[string]struct{}) TensePresentIndicative {
	return TensePresentIndicative{db: db, includes: includes}
}

// Base tense present indicative model type definition
type TensePresentIndicative struct {
	db 				*sqlx.DB `json:"-" db:"-"` 
	includes 		map[string]struct{} `json:"-" db:"-"`  
	Model
	VerbId 			null.Int `json:"verb_id" db:"verb_id"`
	Conjugations
}

func (t TensePresentIndicative) GetReferences() []jsonapi.Reference {
	return []jsonapi.Reference{
		{
			Type: "verbs",
			Name: "verbs",
		},
	}
}

// GetReferencedIDs to satisfy the jsonapi.MarshalLinkedRelations interface
func (t TensePresentIndicative) GetReferencedIDs() []jsonapi.ReferenceID {
	result := []jsonapi.ReferenceID{}
	result = append(result, jsonapi.ReferenceID{
		ID: 	strconv.FormatInt(t.VerbId.Int64, 10),
		Type: 	"verbs",
		Name: 	"verbs",
	})
  return result
}

// GetReferencedStructs to satisfy the jsonapi.MarhsalIncludedRelations interface
func (t TensePresentIndicative) GetReferencedStructs() []jsonapi.MarshalIdentifier {
	if len(t.includes) == 0 {
		return nil
	}
	results := []jsonapi.MarshalIdentifier{} 

	_, ok := t.includes["verb"]
	if ok {
		verb := getRelatedVerb(t.db, t.VerbId)
		results = append(results, verb)
	}

	return results
}

func getRelatedVerb(db *sqlx.DB, id null.Int) jsonapi.MarshalIdentifier {
	verb := NewVerb(db)
	err := db.Get(&verb, "SELECT * FROM verb WHERE id = ? LIMIT 1", id)
	if (err != nil) {
		log.Println("Error retrieving referenced struct for TensePresentIndicative: ", id, err)
		return verb
	}
	return verb
}