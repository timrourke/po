package model

import (
	"github.com/manyminds/api2go/jsonapi"
	"github.com/timrourke/po/database"
	"gopkg.in/guregu/null.v3"
	"log"
	"strconv"
)

/* MySQL Schema

CREATE TABLE `verb` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `aux_verb_id` int(11) DEFAULT NULL,
  `gerund` varchar(255) DEFAULT '',
  `infinitive` varchar(255) NOT NULL,
  `past_participle` varchar(255) DEFAULT '',
  `reflexive` bool DEFAULT 0,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8;

CREATE INDEX infinitive_index ON verb (infinitive) USING BTREE;

*/

// Base verb model type definition
type Verb struct {
	Model
	AuxVerbId      *null.Int `json:"auxiliaryVerb" db:"aux_verb_id"`
	Gerund         *string   `json:"gerund" db:"gerund"`
	Infinitive     *string   `json:"infinitive" db:"infinitive"`
	PastParticiple *string   `json:"pastParticiple" db:"past_participle"`
	Reflexive      *bool     `json:"reflexive" db:"reflexive"`
}

func (v Verb) TableName() string {
	return "verb"
}

// Anonymous conjugations type definition for all verb tenses
type Conjugations struct {
	FirstPersonSingular  *string `json:"firstPersonSingular" db:"sing_first"`
	SecondPersonSingular *string `json:"secondPersonSingular" db:"sing_second"`
	ThirdPersonSingular  *string `json:"thirdPersonSingular" db:"sing_third"`
	FirstPersonPlural    *string `json:"firstPersonPlural" db:"plural_first"`
	SecondPersonPlural   *string `json:"secondPersonPlural" db:"plural_second"`
	ThirdPersonPlural    *string `json:"thirdPersonPlural" db:"plural_third"`
}

func (v Verb) GetReferences() []jsonapi.Reference {
	return []jsonapi.Reference{
		{
			Type: "verbs",
			Name: "auxiliaryVerbs",
		},
		{
			Type: "tensePresentIndicatives",
			Name: "tensePresentIndicatives",
		},
	}
}

// GetReferencedIDs to satisfy the jsonapi.MarshalLinkedRelations interface
func (v Verb) GetReferencedIDs() []jsonapi.ReferenceID {
	result := []jsonapi.ReferenceID{}

	// Get Auxiliary Verb ID
	if v.AuxVerbId != nil {
		result = append(result, jsonapi.ReferenceID{
			ID:   strconv.FormatInt(v.AuxVerbId.Int64, 10),
			Type: "verbs",
			Name: "auxiliaryVerbs",
		})
	}

	// Get Tense Present Indicative ID
	rows, err := database.DB.Queryx("SELECT id FROM tense_pres_ind WHERE verb_id = ?", v.GetID())
	defer rows.Close()
	if err != nil {
		log.Println("Error retrieving referenced ID for Verb: ", v.GetID(), err)
		return []jsonapi.ReferenceID{}
	}
	for rows.Next() {
		var id int64
		err := rows.Scan(&id)
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, jsonapi.ReferenceID{
			ID:   strconv.FormatInt(id, 10),
			Type: "tensePresentIndicatives",
			Name: "tensePresentIndicatives",
		})
	}
	return result
}

func (v Verb) GetReferencedStructs() []jsonapi.MarshalIdentifier {
	results := []jsonapi.MarshalIdentifier{}

	// Get Auxiliary Verb
	auxVerb := Verb{}
	if v.AuxVerbId != nil {
		err := database.DB.Get(&auxVerb, "SELECT * FROM verb WHERE id = ? LIMIT 1", v.AuxVerbId)
		if err != nil {
			log.Println("Error retrieving referenced Auxiliary Verb for Verb: ", v.GetID(), err)
			return []jsonapi.MarshalIdentifier{}
		}
		results = append(results, auxVerb)
	}

	// Get Tense Present Indicative
	rows, err := database.DB.Queryx("SELECT * FROM tense_pres_ind WHERE verb_id = ?", v.GetID())
	defer rows.Close()
	if err != nil {
		log.Println("Error retrieving referenced Tense Present Indicative for Verb: ", v.GetID(), err)
		return []jsonapi.MarshalIdentifier{}
	}
	for rows.Next() {
		t := TensePresentIndicative{}
		err = rows.StructScan(&t)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, t)
	}

	return results
}
