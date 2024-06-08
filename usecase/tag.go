package usecase

import (
	"back/dao"
	"back/model"
	"log"
)

func GetTags() ([]model.Tag, error) {
	tags, err := dao.GetTags()

	if err != nil {
		log.Println("an error occurred at usecase/tweets")
	}
	return tags, err
}
