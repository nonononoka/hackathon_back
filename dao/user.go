package dao

import (
	"back/model"
	"firebase.google.com/go/auth"
	"log"
)

func PostMe(token *auth.Token) (model.User, error) {
	var userInfo model.User

	tx, err := db.Begin()
	if err != nil {
		log.Printf("fail: db.Begin, %v\n", err)
		return userInfo, err
	}
	_, err = tx.Exec("INSERT INTO users (email, id) VALUES(? ,?)", token.Claims["email"].(string), token.UID)
	if err != nil {
		log.Printf("fail: tx.Exec, %v\n", err)
		tx.Rollback()
		return userInfo, err
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		log.Printf("fail: tx.Commit, %v\n", err)
		return userInfo, err
	}
	userInfo.ID = token.UID
	userInfo.Email = token.Claims["email"].(string)
	return userInfo, nil
}

func GetMe(token *auth.Token) (model.User, error) {
	var userInfo model.User

	err := db.QueryRow("select id, email from users where id = ?", token.UID).Scan(&userInfo.ID, &userInfo.Email)

	if err != nil {
		return userInfo, err
	}
	log.Println(userInfo)
	return userInfo, err
}
