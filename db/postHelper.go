package db

import (
	"achyuBackend/yolp"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func InsertPost(postId, email, ip, displayName, content, imageURL, userWhere string, latitude, longitude float64, dataset *yolp.DataSet) error {
	db, err := sql.Open("sqlite3", "./database.sqlite")
	if err != nil {
		return err
	}
	defer db.Close()

	cmd := "INSERT INTO posts (postId,email,ip,displayName,content,imageURL,userWhere,latitude,longitude,address,constructionName,roadName) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)"
	_, err = db.Exec(cmd, postId, email, ip, displayName, content, imageURL, userWhere, latitude, longitude, dataset.Address, dataset.PlaceName, dataset.RoadName)
	if err != nil {
		return err
	}

	return nil
}

func MovePost(email, postId string) error {
	db, err := sql.Open("sqlite3", "./database.sqlite")
	if err != nil {
		return err
	}
	defer db.Close()

	cmd := "SELECT ip,displayName,content,imageURL,userWhere,latitude,longitude,address,constructionName,roadName FROM posts WHERE email=? AND postId=?"
	row := db.QueryRow(cmd, email, postId)
	var ip string
	var displayName string
	var content string
	var imageURL string
	var userWhere string
	var latitude float64
	var longitude float64
	var address string
	var constructionName string
	var roadName string
	err = row.Scan(&ip, &displayName, &content, &imageURL, &userWhere, &latitude, &longitude, &address, &constructionName, &roadName)
	if err != nil {
		return err
	}

	cmd = "DELETE FROM posts WHERE email=? AND postId=?"
	_, err = db.Exec(cmd, email, postId)
	if err != nil {
		return err
	}
	cmd = "INSERT INTO deleted (postId,email,ip,displayName,content,imageURL,userWhere,latitude,longitude,address,constructionName,roadName) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)"
	_, err = db.Exec(cmd, postId, email, ip, displayName, content, imageURL, userWhere, latitude, longitude, address, constructionName, roadName)
	if err != nil {
		return err
	}

	return nil
}

type ResponseGet struct {
	PostId           string  `json:"postId"`
	DisplayName      string  `json:"displayName"`
	Content          string  `json:"content"`
	ImageURL         string  `json:"imageUrl"`
	UserWhere        string  `json:"userWhere"`
	Latitude         float64 `json:"latitude"`
	Longitude        float64 `json:"longitude"`
	Address          string  `json:"address"`
	ConstructionName string  `json:"constructionName"`
	RoadName         string  `json:"roadName"`
}

func GetPost(latitude, longitude float64) ([]ResponseGet, error) {
	db, err := sql.Open("sqlite3", "./database.sqlite")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	cmd := "SELECT postId,displayName,content,imageURL,userWhere,latitude,longitude,address,constructionName,roadName FROM posts WHERE latitude BETWEEN ? - 0.4 AND ? + 0.4 AND longitude BETWEEN ? - 0.4 AND ? + 0.4"
	row, err := db.Query(cmd, latitude, latitude, longitude, longitude)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	var response []ResponseGet
	for row.Next() {
		var res ResponseGet
		err := row.Scan(&res.PostId, &res.DisplayName, &res.Content, &res.ImageURL, &res.UserWhere, &res.Latitude, &res.Longitude, &res.Address, &res.ConstructionName, &res.RoadName)
		if err != nil {
			return nil, err
		}
		response = append(response, res)
	}

	if err := row.Err(); err != nil {
		return nil, err
	}

	return response, nil
}
