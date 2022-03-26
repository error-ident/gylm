package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"log"
	"math/rand"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/joho/godotenv"
	"github.com/shkh/lastfm-go/lastfm"
)

var (
	trackName  string
	artistName string
	imageData  string
	userUrl    string
	imageUrl   string
	theme      string
)

type TrackData struct {
	TrackName  string
	ArtistName string
	Image      string
	UserUrl    string
	ThemeName  string
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("No .env file found")
	}

	if os.Getenv("THEME") == "" {
		theme = "default"
	} else {
		theme = os.Getenv("THEME")
	}
}

func main() {
	key := os.Getenv("KEY")
	secret := os.Getenv("SECRET")
	user := os.Getenv("USERNAME")
	limit := os.Getenv("LIMIT")
	userUrl = os.Getenv("YANDEX_URL")

	http.HandleFunc("/yandex", func(w http.ResponseWriter, r *http.Request) {
		api := lastfm.New(key, secret)

		result, err := api.User.GetRecentTracks(lastfm.P{"limit": limit,
			"user": user})
		if err != nil {
			return
		}

		for _, track := range result.Tracks {
			if track.NowPlaying == "true" {
				trackName = track.Name
				artistName = track.Artist.Name
				imageUrl = track.Images[3].Url

				break
			} else {
				rand.Seed(time.Now().UnixNano())
				track := result.Tracks[rand.Intn(len(result.Tracks))]

				trackName = track.Name
				artistName = track.Artist.Name
				imageUrl = track.Images[3].Url

				fmt.Println(track.Artist.Name, track.Name)
				break
			}
		}

		imageData, err = trackImageToBase64(imageUrl)
		if err != nil {
			log.Printf("image64 error: %v", err)
		}

		data := TrackData{
			TrackName:  trackName,
			ArtistName: artistName,
			Image:      imageData,
			UserUrl:    userUrl,
			ThemeName:  theme,
		}

		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			fmt.Printf("template execution: %s", err)
		}

		w.Header().Set("Content-Type", "image/svg+xml")
		w.Header().Set("Cache-Control", "max-age=0")
		if err := tmpl.Execute(w, data); err != nil {
			log.Printf("Error executing template: %v", err)
		}

	})

	fs := http.FileServer(http.Dir("./themes"))
	http.Handle("/themes/", http.StripPrefix("/themes/", fs))

	//server log
	err := http.ListenAndServe(":1984", nil)
	if err != nil {
		fmt.Println("ListenAndServe: ", err)
	}

}

func trackImageToBase64(url string) (string, error) {
	responseImage, err := http.Get(url)
	if err != nil || responseImage.StatusCode != 200 {
		log.Printf("image not found")
	}
	defer responseImage.Body.Close()
	imageBody, _, err := image.Decode(responseImage.Body)
	if err != nil {
		log.Printf("image decode is failed, imageUrl: %s, %v", url, err)
	}

	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, imageBody, nil)
	if err != nil {
		log.Printf("jpeg encode is failed")
	}

	imageByte := buf.Bytes()
	imageData = base64.StdEncoding.EncodeToString(imageByte)

	return imageData, nil
}
