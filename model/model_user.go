package model

import (
// "github.com/timrourke/po/database"
//"errors"
// "fmt"
// "sort"
// "strconv"
//"github.com/manyminds/api2go/jsonapi"
)

/* MySQL Schema

CREATE TABLE `user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `username` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY (`username`),
  UNIQUE KEY (`email`)
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8;

*/

// Base noun model type definition
type User struct {
	Model
	Username *string `json:"username" db:"username"`
	Email    *string `json:"email" db:"email"`
	Password *string `json:"-" db:"password"`
}

func (u User) TableName() string {
	return "user"
}

type SignupUser struct {
	Model
	Username *string `json:"username" db:"username"`
	Email    *string `json:"email" db:"email"`
	Password *string `json:"password" db:"password"`
}
