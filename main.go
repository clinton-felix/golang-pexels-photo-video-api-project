package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	PhotoApi = "https://api.pexels.com/v1"
	VideoApi = "https://api.pexels.com/videos"
)

// creating a client struct
type Client struct{
	Token string
	hc http.Client
	RemainingTimes int32		// remaining times the API can be called
}

// initializing a NewClient function
func NewClient(token string) *Client {
	c := http.Client{}
	return &Client{Token: token, hc: c}
}


/* 
	Creating a struct to filter the search results returned from the API
	so that only relevant fields of the API are siloed for use
*/
type SearchResults struct{
	Page 			int32 		`json:"page"`
	PerPage 		int32 		`json:"per_page"`
	TotalResults 	int32 		`json:"total_results"`
	NextPage 		string 		`json:"next_page"`
	Photos 			[]Photo		`json:"photos"`		// Photos is a slice of the struct Photo
}

// creating the Photo struct
type Photo struct{
	Id				int32		`json:"id"`
	Width			int32		`json:"width"`
	Height			int32		`json:"height"`
	Url				string		`json:"url"`
	Photographer	string		`json:"photographer"`
	PhotographerUrl	string		`json:"photographer_url"`
	Src				PhotoSource	`json:"src"`
}


// creating the PhotoSource struct
type PhotoSource struct {
	Original		string		`json:"original"`
	Large			string		`json:"large"`
	Large2x			string		`json:"large2x"`
	Medium			string		`json:"medium"`
	Small			string		`json:"small"`
	Potrait			string		`json:"potrait"`
	Square			string		`json:"square"`
	Landscape		string		`json:"landscape"`
	Tiny			string		`json:"tiny"`
}

// creating the curatedresults struct
type CuratedResult struct {
	Page			int32		`json:"page"`
	PerPage 		int32		`json:"per_page"`
	NextPage 		string 		`json:"next_page"`
	Photos			[]Photo 	`json:"photos"`
}

// creating the video struct
type Video struct{
	Id				int32			`json:"id"`
	Width			int32			`json:"width"`
	Height			int32			`json:"height"`
	Url				string			`json:"url"`
	Image			string			`json:"image"`
	FullRes			interface{}		`json:"full_res"`
	Duration		float64			`json:"duration"`
	VideoFiles		[]VideoFiles	`json:"video_files"`
	VideoPictures	[]VideoPictures	`json:"video_pictures"`
}

// creatng the videoSearch results
type VideoSearchResult struct{
	Page			int32		`json:"page"`
	PerPage			int32		`json:"per_page"`
	TotalResults	int32		`json:"total_results"`
	NextPage		string		`json:"next_page"`
	Videos 			[]Video 	`json:"videos"`
}

// creating the popular video struct
type PopularVideos struct{
	Page			int32		`json:"page"`
	PerPage 		int32		`json:"per_page"`
	TotalResults	int32		`json:"total_results"`
	Url				string		`json:"url"`
	Videos			[]Video		`json:"videos"`
}

// creating the Video file struct
type VideoFiles struct {
	Id 				int32		`json:"id"`
	Quality			string		`json:"quality"`
	FileType		string		`json:"file_type"`
	Width			int32		`json:"width"`
	Height			int32		`json:"height"`
	Link			string		`json:"link"`
}

// creating the Video Pictures struct
type VideoPictures struct {
	Id 				int32 		`json:"id"`
	Picture 		string		`json:"picture"`
	Nr				int32		`json:"nr"`
}

/* ============ End of structs ========== */

/* 
	Creating SearchPhotos method for Client struct. Takes query, perpage and page
	as input params and returns *SearchResults and error
*/
func (c *Client) SearchPhotos(query string, perPage, page int32) (*SearchResults, error) {
	// formating our url request string
	url := fmt.Sprintf(PhotoApi+"/search?query=%s&per_page=%d&page=%d", query, perPage, page)
	// get the response using the method, and handle errors
	resp, err := c.requestDoWithAuth("GET", url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// assign the response from response body using ioutil.ReadAll to data recieved
	data, err := ioutil.ReadAll(resp.Body )
	if err != nil {
		return nil, err		// since the parent function returns two values
	}
	var result SearchResults
	err = json.Unmarshal(data, &result)
	return &result, err
}


// Creating the curatedPhotos method
func (c *Client) CuratedPhotos(perPage, page int32) (*CuratedResult, error)  {
	url := fmt.Sprintf(PhotoApi+"/curated?per_page=%d&page=%d", perPage, page)
	// getting the response
	resp, err := c.requestDoWithAuth("GET", url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// assigning body of response to data using ioutils.ReadAll
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result CuratedResult
	err = json.Unmarshal(data, &result)
	return &result, err
}

// creating the getPhoto by Id method
func (c *Client) GetPhoto(id int32) (*Photo, error)  {
	url := fmt.Sprintf(PhotoApi+"/photos/%d", id)
	resp, err := c.requestDoWithAuth("GET", url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// assigning the body of the response to "data", using ioutils.ReadAll
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result Photo
	err = json.Unmarshal(data, &result)
	return &result, err
}

// Creating a getRandomPhoto method
func (c *Client) GetRandomPhoto() (*Photo, error)  {
	// to get a random photo, we need to first create a random number
	rand.Seed(time.Now().Unix())
	randNum := rand.Intn(1001)
	result, err := c.CuratedPhotos(1, int32(randNum))
	if err == nil && len(result.Photos) == 1 {
		return &result.Photos[0], nil
	}
	return nil, err
}


/* 
	creating the requestDoWithAuth client method which takes method and url as params
	and returns a http response and and error.
	
	*NB: auth represents the token
*/
func (c *Client) requestDoWithAuth(method, url string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)
	// handling error first
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", c.Token)
	resp, err := c.hc.Do(req)
	if err != nil {
		return resp, err
	}

	// setting the remaining times API can be called by client param
	times, err := strconv.Atoi(resp.Header.Get("X-Ratelimit-Remaining"))
	if err != nil {
		return nil, err
	} else {
		c.RemainingTimes = int32(times)
	}
	return resp, err
}

// Creating the SearchVideo Method
func (c *Client) SearchVideo(query string, perPage, page int32) (*VideoSearchResult, error) {
	url := fmt.Sprintf(VideoApi+"/search?query=%s&per_page=%d&page=%d", query, perPage, page)
	resp, err := c.requestDoWithAuth("GET", url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var  result VideoSearchResult
	err = json.Unmarshal(data, &result)
	return &result, err
}

// Creating the Popular video Method
func (c *Client) PopularVideo(perPage, page int32) (*PopularVideos, error) {
	url := fmt.Sprintf(VideoApi+"/popular?per_page=%d&page=%d", perPage,page)
	resp, err := c.requestDoWithAuth("GET", url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result PopularVideos
	err = json.Unmarshal(data, &result)
	return &result, err
}


// Creating the GetRandomVideo Method
func (c *Client) GetRandomVideo() (*Video, error) {
	rand.Seed(time.Now().Unix())
	randNum := rand.Intn(1001)
	result, err := c.PopularVideo(1, int32(randNum))
	if err == nil && len(result.Videos) == 1 {
		return &result.Videos[0], nil
	}
	return nil, err
}

// Creating a method to get the Remaining Api Call request in the month
func (c *Client) GetRemainingRequestInThisMonth() int32 {
	return c.RemainingTimes
}

/* ======= End of method definitions */


func main() {
	os.Setenv("PexelsToken", "563492ad6f91700001000001ab0efaa0bab643439c5fbcf65d13602c")
	Token := os.Getenv("PexelsToken")

	// defining a function for a client based on the Token
	var c = NewClient(Token)

	// calling the SearchPhotos method on the client, to the pexels API for results..
	result, err := c.SearchPhotos("waves", 15, 1)
	if err != nil {
		fmt.Errorf("Search error: %v", err)
	}
	if result.Page == 0 {
		fmt.Errorf("Search result is wrong..")
	}
	// return results if everything works well
	fmt.Println(result)
}


// API: 563492ad6f91700001000001ab0efaa0bab643439c5fbcf65d13602c