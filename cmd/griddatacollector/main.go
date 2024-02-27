// griddatacollector : connect to terna api via http and downloads energy generation data
// in json format

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"time"
)

// API Token response struct
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func main() {

	// command line flags parsing

	apiHost := flag.String("host", "https://api.terna.it", "terna api host")
	clientID := flag.String("clientid", "q8wdvqt4xgn7xzxfbve7ev8e", "terna api client id")
	clientSecret := flag.String("secret", "se3CkKvfnb", "terna api client secret")

	flag.Parse()

	// generate access token

	// curl -H "Content-Type: application/x-www-form-urlencoded" "https://api.terna.it/transparency/oauth/accessToken" -X POST -d "client_id=5d4rscEcpyxywu4jdoiWerhsl" -d "client_secret=Iy4c6tuErp" -d "grant_type=client_credentials"

	params := fmt.Sprintf("client_id=%s&client_secret=%s&grant_type=client_credentials", *clientID, *clientSecret)
	data := []byte(params)

	req, _ := http.NewRequest("POST", *apiHost+"/transparency/oauth/accessToken", bytes.NewBuffer(data))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	// Parse JSON response
	body, _ := io.ReadAll(resp.Body)
	var tr TokenResponse
	json.Unmarshal(body, &tr)

	// Access token
	accessToken := tr.AccessToken

	fmt.Printf("response: %v\n", tr)
	resp.Body.Close()

	// using access token to download data

	// https://api.terna.it/generation/v2.0/actual-generation?dateFrom=2/3/2019&dateTo=20/3/2019&type=Wind&type=Hydro
	// curl example
	// curl --location --request GET 'https://api.terna.it/generation/v2.0/actual-generation?dateFrom=23/04/2022&dateTo=23/04/2022' --header 'Authorization: <token>'

	// today := time.Now().Format("02/01/2006")
	// calculate the day before from today
	yesterday := time.Now().AddDate(0, 0, -1).Format("02/01/2006")
	// the day before yesterday
	beforeYesterday := time.Now().AddDate(0, 0, -2).Format("02/01/2006")

	url := *apiHost + "/generation/v2.0/actual-generation?dateFrom=" + beforeYesterday + "&dateTo=" + yesterday + "&type=Wind&type=Hydro"
	// url := *apiHost + "/transparency/v1.0/getactualgeneration?dateFrom=" + beforeYesterday + "&dateTo=" + yesterday + "&type=Wind&type=Hydro"

	fmt.Printf("url = %s\n", url)

	req, _ = http.NewRequest("GET", url, nil)

	req.Header.Add("Authorization", accessToken)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ = io.ReadAll(res.Body)

	fmt.Println(string(body))

	resp.Body.Close()

}
