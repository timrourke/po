package model

import (
	"github.com/manyminds/api2go/jsonapi"
	"github.com/timrourke/po/database"
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

func NewTensePresentIndicative(includes map[string]struct{}) TensePresentIndicative {
	return TensePresentIndicative{includes: includes}
}

// Base tense present indicative model type definition
type TensePresentIndicative struct {
	Model
	VerbId *null.Int `json:"verb_id" db:"verb_id"`
	Conjugations
	includes map[string]struct{} `json:"-" db:"-"`
}

func (t TensePresentIndicative) TableName() string {
	return "tense_pres_ind"
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
		ID:   strconv.FormatInt(t.VerbId.Int64, 10),
		Type: "verbs",
		Name: "verbs",
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
		verb := getRelatedVerb(*t.VerbId)
		results = append(results, verb)
	}

	return results
}

func getRelatedVerb(id null.Int) jsonapi.MarshalIdentifier {
	verb := Verb{}
	err := database.DB.Get(&verb, "SELECT * FROM verb WHERE id = ? LIMIT 1", id)
	if err != nil {
		log.Println("Error retrieving referenced struct for TensePresentIndicative: ", id, err)
		return verb
	}
	return verb
}
