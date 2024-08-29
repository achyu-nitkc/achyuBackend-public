package yolp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type DataSet struct {
	Address   string
	PlaceName string
	RoadName  string
}

type ResultItem struct {
	Label *string `json:"Label"`
}

type JsonField struct {
	Address  []string     `json:"Address"`
	RoadName *string      `json:"Roadname"`
	Result   []ResultItem `json:"Result"`
}

type ResSet struct {
	ResultSet JsonField `json:"ResultSet"`
}

func GetAddressData(latitude, longitude float64) (*DataSet, error) {
	strLatitude := fmt.Sprintf("lat=%f", latitude)
	strLongitude := fmt.Sprintf("lon=%f", longitude)
	clientId := yolpSecret()

	yolpUrl := "https://map.yahooapis.jp/placeinfo/V1/get?" + strLatitude + "&" + strLongitude + "&" + clientId + "&output=json"
	resp, err := http.Get(yolpUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var resultSet ResSet
	err = json.Unmarshal(body, &resultSet)
	if err != nil {
		return nil, err
	}
	var placeName string
	if len(resultSet.ResultSet.Result) == 0 {
		placeName = ""
	} else if resultSet.ResultSet.Result[0].Label != nil {
		placeName = *resultSet.ResultSet.Result[0].Label
	} else {
		placeName = ""
	}
	address := strings.Join(resultSet.ResultSet.Address, " ")
	var roadName string
	if resultSet.ResultSet.RoadName != nil {
		roadName = *resultSet.ResultSet.RoadName
	} else {
		roadName = ""
	}
	return &DataSet{
		Address:   address,
		PlaceName: placeName,
		RoadName:  roadName,
	}, nil
}
