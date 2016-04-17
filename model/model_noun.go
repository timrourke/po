package model

import (
  //"errors"
  // "fmt"
  // "sort"
  // "strconv"

  //"github.com/manyminds/api2go/jsonapi"
)

/* MySQL Schema

CREATE TABLE `noun` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `singular` varchar(255) NOT NULL,
  `plural` varchar(255) DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8;

*/

// Base noun model type definition
type Noun struct {
  Model
  Singular  string `json:"singular" db:"singular"`
  Plural    string `json:"plural" db:"plural"`
}