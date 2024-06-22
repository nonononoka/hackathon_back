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
	_, err = tx.Exec("INSERT INTO user (id, name, email, bio, image) VALUES(? ,?, ?, ?, ?)", token.UID, token.Claims["name"].(string), token.Claims["email"].(string), nil, nil)
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
	err = db.QueryRow("select id, name, email, bio, image from user where id = ?", token.UID).Scan(&userInfo.ID, &userInfo.Name, &userInfo.Email, &userInfo.Bio, &userInfo.Image)
	if err != nil {
		log.Printf("fail: db.QueryRow, %v\n", err)
		return userInfo, err
	}
	return userInfo, nil
}

func PutMe(token *auth.Token, name string, bio string, image string) (model.User, error) {
	var userInfo model.User

	tx, err := db.Begin()
	if err != nil {
		log.Printf("fail: db.Begin, %v\n", err)
		return userInfo, err
	}
	_, err = tx.Exec("UPDATE user SET name = ?, bio = ?, image = ? WHERE id = ?;", name, bio, image, token.UID)
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
	err = db.QueryRow("select id, name, email, bio, image from user where id = ?", token.UID).Scan(&userInfo.ID, &userInfo.Name, &userInfo.Email, &userInfo.Bio, &userInfo.Image)
	if err != nil {
		log.Printf("fail: db.QueryRow, %v\n", err)
		return userInfo, err
	}
	return userInfo, nil
}

func GetMe(token *auth.Token) (model.User, error) {
	var userInfo model.User

	err := db.QueryRow("select id, name, email, bio, image from user where id = ?", token.UID).Scan(&userInfo.ID, &userInfo.Name, &userInfo.Email, &userInfo.Bio, &userInfo.Image)

	if err != nil {
		log.Println(err.Error())
		return userInfo, err
	}
	return userInfo, err
}

func GetUsers(token *auth.Token) ([]model.User, error) {
	userInfos := make([]model.User, 0)

	rows, err := db.Query("select id, name, email, bio, image from user;")
	if err != nil {
		log.Printf(err.Error())
		return userInfos, err
	}

	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Bio, &u.Image); err != nil {
			log.Printf(err.Error())
			return userInfos, err
		}
		following, followed, err := checkFollow(token.UID, u.ID)
		if err != nil {
			return userInfos, err
		}
		// uにフォローされているか
		u.IsFollowed = followed
		// uをフォローしているか
		u.IsFollowing = following
		if token.UID != u.ID {
			userInfos = append(userInfos, u)
		}
	}
	if err := rows.Close(); err != nil {
		log.Printf(err.Error())
		return userInfos, err
	}
	return userInfos, nil
}
