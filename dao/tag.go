package dao

import (
	"back/model"
	"log"
)

func GetTags() ([]model.Tag, error) {
	var tags []model.Tag

	rows, err := db.Query("select * from tag;")

	if err != nil {
		log.Printf(err.Error())
		return tags, err
	}

	for rows.Next() {
		var t model.Tag
		if err := rows.Scan(&t.ID, &t.Tag); err != nil {
			log.Printf(err.Error())
			return tags, err
		}

		tags = append(tags, t)
	}
	if err := rows.Close(); err != nil {
		log.Printf(err.Error())
		return tags, err
	}
	return tags, nil
}
