package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gyan1230/asat/controllers"
)

type Information struct {
	CreatedAt string `json:"created_at"`
	ID        int64  `json:"id"`
	IDStr     string `json:"id_str"`
	Text      string `json:"text"`
	Truncated bool   `json:"truncated"`
	Entities  struct {
		Hashtags     []interface{} `json:"hashtags"`
		Symbols      []interface{} `json:"symbols"`
		UserMentions []interface{} `json:"user_mentions"`
		Urls         []interface{} `json:"urls"`
	} `json:"entities"`
	Source               string      `json:"source"`
	InReplyToStatusID    interface{} `json:"in_reply_to_status_id"`
	InReplyToStatusIDStr interface{} `json:"in_reply_to_status_id_str"`
	InReplyToUserID      interface{} `json:"in_reply_to_user_id"`
	InReplyToUserIDStr   interface{} `json:"in_reply_to_user_id_str"`
	InReplyToScreenName  interface{} `json:"in_reply_to_screen_name"`
	User                 struct {
		ID          int64       `json:"id"`
		IDStr       string      `json:"id_str"`
		Name        string      `json:"name"`
		ScreenName  string      `json:"screen_name"`
		Location    string      `json:"location"`
		Description string      `json:"description"`
		URL         interface{} `json:"url"`
		Entities    struct {
			Description struct {
				Urls []interface{} `json:"urls"`
			} `json:"description"`
		} `json:"entities"`
		Protected                      bool        `json:"protected"`
		FollowersCount                 int         `json:"followers_count"`
		FriendsCount                   int         `json:"friends_count"`
		ListedCount                    int         `json:"listed_count"`
		CreatedAt                      string      `json:"created_at"`
		FavouritesCount                int         `json:"favourites_count"`
		UtcOffset                      interface{} `json:"utc_offset"`
		TimeZone                       interface{} `json:"time_zone"`
		GeoEnabled                     bool        `json:"geo_enabled"`
		Verified                       bool        `json:"verified"`
		StatusesCount                  int         `json:"statuses_count"`
		Lang                           interface{} `json:"lang"`
		ContributorsEnabled            bool        `json:"contributors_enabled"`
		IsTranslator                   bool        `json:"is_translator"`
		IsTranslationEnabled           bool        `json:"is_translation_enabled"`
		ProfileBackgroundColor         string      `json:"profile_background_color"`
		ProfileBackgroundImageURL      string      `json:"profile_background_image_url"`
		ProfileBackgroundImageURLHTTPS string      `json:"profile_background_image_url_https"`
		ProfileBackgroundTile          bool        `json:"profile_background_tile"`
		ProfileImageURL                string      `json:"profile_image_url"`
		ProfileImageURLHTTPS           string      `json:"profile_image_url_https"`
		ProfileLinkColor               string      `json:"profile_link_color"`
		ProfileSidebarBorderColor      string      `json:"profile_sidebar_border_color"`
		ProfileSidebarFillColor        string      `json:"profile_sidebar_fill_color"`
		ProfileTextColor               string      `json:"profile_text_color"`
		ProfileUseBackgroundImage      bool        `json:"profile_use_background_image"`
		HasExtendedProfile             bool        `json:"has_extended_profile"`
		DefaultProfile                 bool        `json:"default_profile"`
		DefaultProfileImage            bool        `json:"default_profile_image"`
		Following                      interface{} `json:"following"`
		FollowRequestSent              interface{} `json:"follow_request_sent"`
		Notifications                  interface{} `json:"notifications"`
		TranslatorType                 string      `json:"translator_type"`
	} `json:"user"`
	Geo           interface{} `json:"geo"`
	Coordinates   interface{} `json:"coordinates"`
	Place         interface{} `json:"place"`
	Contributors  interface{} `json:"contributors"`
	IsQuoteStatus bool        `json:"is_quote_status"`
	RetweetCount  int         `json:"retweet_count"`
	FavoriteCount int         `json:"favorite_count"`
	Favorited     bool        `json:"favorited"`
	Retweeted     bool        `json:"retweeted"`
	Lang          string      `json:"lang"`
}

func main() {
	log.Println("Starting application :::::::::")
	http.HandleFunc("/", index)
	http.HandleFunc("/register", controllers.Register)
	http.HandleFunc("/admin", controllers.ShowAll)
	http.HandleFunc("/login", controllers.Login)
	http.HandleFunc("/logout", controllers.Logout)
	http.HandleFunc("/tweetData", getTweetData)

	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Println("Server listen error:", err)
	}

}

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	log.Println("In index:::::::::::::::")
	json.NewEncoder(w).Encode("In Index:::")
}

func getTweetData(w http.ResponseWriter, req *http.Request) {
	url := "https://api.twitter.com/1.1/statuses/user_timeline.json?user_id=3120243180&screen_name=vikramsparamesh"

	req, _ = http.NewRequest("GET", url, nil)

	req.Header.Add("Authorization", "Bearer AAAAAAAAAAAAAAAAAAAAABlA%2BgAAAAAAMWkicp39DkIQkGSe0nQrCmOJMNg%3Dm9bcjqeL4jlC2PVL9K0qz8gGr3jNXjnUidbxbCKuej9j0hDDtD")
	req.Header.Add("User-Agent", "PostmanRuntime/7.13.0")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Cache-Control", "no-cache")
	//	req.Header.Add("Postman-Token", "e55120bb-292b-4cc4-8b01-5e5181230ce2,51e598a2-07ae-41b7-acbd-5302fec2fb1c")
	req.Header.Add("Host", "api.twitter.com")
	req.Header.Add("cookie", "personalization_id=v1_tf1dDIphNvn4tBs3u8LUkA==")
	req.Header.Add("cookie", "guest_id=v1%3A155854116672032548")
	req.Header.Add("accept-encoding", "gzip, deflate")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("cache-control", "no-cache")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln("error in response fetch", err)
	}
	defer res.Body.Close()

	// Check that the server actually sent compressed data
	var reader io.ReadCloser
	switch res.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(res.Body)
		defer reader.Close()
	default:
		reader = res.Body
	}

	var info []Information
	w.Header().Set("content-type", "application/json")
	err = json.NewDecoder(reader).Decode(&info)
	if err != nil {
		log.Printf("error decoding response: %v", err)
		if e, ok := err.(*json.SyntaxError); ok {
			log.Printf("syntax error at byte offset %d", e.Offset)
		}

	}

	//	w.Write([]byte(`{ "message": "` + err.Error() + `" }`))

	//w.Write([]byte(`{ "User details":"` + info[0].User.ScreenName + `" }`))
	count := 0
	for i := 0; i < 4; i++ {
		count += info[i].FavoriteCount
	}
	w.Write([]byte(`{ "Tweets like count":"` + string(count) + `" }`))

	fmt.Println("User screen name", info[0].User.ScreenName)
	fmt.Println("count", count)
	//	fmt.Println("info is", info)
}
