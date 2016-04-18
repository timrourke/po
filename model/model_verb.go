package model

import (
  //"errors"
  // "fmt"
  // "sort"
  // "strconv"
  "database/sql"
  "github.com/manyminds/api2go/jsonapi"
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
  AuxVerbId       sql.NullInt64 `json:"auxiliaryVerb" db:"aux_verb_id"`
  Gerund          string        `json:"gerund" db:"gerund"`
  Infinitive      string        `json:"infinitive" db:"infinitive"`
  PastParticiple  string        `json:"pastParticiple" db:"past_participle"`
  Reflexive       bool          `json:"reflexive" db:"reflexive"`
}

// Anonymous conjugations type definition for all verb tenses
type Conjugations struct {
  FirstPersonSingular   string `json:"firstPersonSingular" db:"sing_first"`
  SecondPersonSingular  string `json:"secondPersonSingular" db:"sing_second"`
  ThirdPersonSingular   string `json:"thirdPersonSingular" db:"sing_third"`
  FirstPersonPlural     string `json:"firstPersonPlural" db:"plural_first"`
  SecondPersonPlural    string `json:"secondPersonPlural" db:"plural_second"`
  ThirdPersonPlural     string `json:"thirdPersonPlural" db:"plural_third"`
}

func (v Verb) GetReferences() []jsonapi.Reference {
  return []jsonapi.Reference{
    {
      Type: "verb",
      Name: "auxiliaryVerb",
    },

  }
}