# golang-pexels-photo-video-api-project

This is a Pexels API project using GoLang, to consume the PEXELS API

## Steps

1. Clone this repo using git clone "https://github.com/clinton-felix/golang-pexels-photo-video-api-project"
2. Create a .env file and add a Pexels API key from pexels.com/api
3. Name the API as "PEXELS_API"
4. Run the comand in your terminal: "go run main.go"

## Functions available in Code

Several methods have been defined in the Code

1. SearchPhotos() function method which returns a list of photos

2. CuratedPhotos() function method which returns a slice of curated photos

3. GetPhoto() function method which gets a single photo by its ID

4. GetRandomPhoto() function method which gets a random photo from the API

5. SearchVideo() function method queries the API and returns Videos

6. PopularVideo() function method which returns the list of popular videos

7. GetRandomVideo() function method which returns a random video
