package model

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

// Base tense present indicative model type definition
type TensePresentIndicative struct {
	Model
	VerbId uint64 `json:"verb" db:"verb_id"`
	Conjugations
}
