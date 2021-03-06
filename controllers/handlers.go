package controllers

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gyan1230/asat/config"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mikespook/gorbac.v2"
)

//Information :
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

//PowerDataStr :
type PowerDataStr struct {
	IndexName    string    `json:"index_name"` // `json:"index_name,omitempty" bson:"index_name,omitempty"`
	Title        string    `json:"title"`
	Desc         string    `json:"desc"`
	Created      int       `json:"created"`
	Updated      int       `json:"updated"`
	CreatedDate  time.Time `json:"created_date"`
	UpdatedDate  time.Time `json:"updated_date"`
	Active       string    `json:"active"`
	Visualizable string    `json:"visualizable"`
	CatalogUUID  string    `json:"catalog_uuid"`
	Source       string    `json:"source"`
	OrgType      string    `json:"org_type"`
	Org          []string  `json:"org"`
	Sector       []string  `json:"sector"`
	Field        []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"field"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Total   int    `json:"total"`
	Count   int    `json:"count"`
	Limit   string `json:"limit"`
	Offset  string `json:"offset"`
	Records []struct {
		StateSystemRegion                    string `json:"state_system_region"`
		October2018PeakDemandMw              string `json:"october_2018_peak_demand__mw_"`
		October2018PeakMetMw                 string `json:"october_2018_peak_met__mw_"`
		October2018SurplusDeficitMw          string `json:"october_2018_surplus_deficit_____mw_"`
		October2018SurplusDeficit            string `json:"october_2018_surplus_deficit_____"`
		April2018October2018PeakDemandMw     string `json:"april_2018_october_2018_peak_demand__mw_"`
		April2018October2018PeakMetMw        string `json:"april_2018_october_2018_peak_met__mw_"`
		April2018October2018SurplusDeficitMw string `json:"april_2018_october_2018_surplus_deficit_____mw_"`
		April2018October2018SurplusDeficit   string `json:"april_2018_october_2018_surplus_deficit_____"`
	} `json:"records"`
	Version string `json:"version"`
}

var jwtKey = []byte("my_secret_key")

//ShowAll :
func ShowAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	usr, err := Allusers(r)
	if err != nil {
		http.Error(w, http.StatusText(500)+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usr)
	return
}

//Register ...
func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("content-type", "application/json")
	var person User
	_ = json.NewDecoder(r.Body).Decode(&person)
	p, err := GetUser(r.Context(), person.Email)
	if err != nil {
		bs, _ := bcrypt.GenerateFromPassword([]byte(person.Password), bcrypt.MinCost)
		person.Password = string(bs)
		collection := config.Client.Database("userDb").Collection("user")
		result, _ := collection.InsertOne(r.Context(), person)
		json.NewEncoder(w).Encode(result)
		w.WriteHeader(http.StatusCreated)
		return
	}
	if p != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("User already exist"))
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("User could not be created"))
	return
}

//Login ...
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	if AlreadyLoggedIn(w, r) {
		log.Println("Already login return to home...")
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode("Already login return to home.")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// is there a Email?
	var person Credential
	_ = json.NewDecoder(r.Body).Decode(&person)
	u, err := GetUser(r.Context(), person.Email) // return user (if present), nil in u,err OR nil,err (if not present user) in u,err
	if u == nil {
		log.Println("Email not exists")
		http.Error(w, "Email not exists", http.StatusForbidden)
		return
	}
	// does the entered password match the stored password?
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(person.Password))
	if err != nil {
		log.Println("Email and/or password do not match")
		http.Error(w, "Email and/or password do not match", http.StatusForbidden)
		return
	}

	log.Println("sucessfully login")

	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	expirationTime := time.Now().Add(15 * time.Second)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Username: u.Email,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}
	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string >> get the complete signed token
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// create session
	c := &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	}
	// c.MaxAge = sessionLength

	http.SetCookie(w, c)
	fmt.Println("login cookie set ::::")

	tmp := struct {
		Email    string `json:"email"`
		Fullname string `json:"fullname"`
	}{
		Email:    u.Email,
		Fullname: u.Fullname,
	}

	res := Resp{Data: tmp}

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(res)
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}

//Logout ...
func Logout(w http.ResponseWriter, r *http.Request) {
	if !AlreadyLoggedIn(w, r) {
		log.Println("Return to index :::")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	c, _ := r.Cookie("token")
	// delete the session
	// delete(dbSessions, c.Value)
	// remove the cookie
	c = &http.Cookie{
		Name:   "token",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, c)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

//GetTweetData :
func GetTweetData(w http.ResponseWriter, req *http.Request) {
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
	for _, v := range info {
		count += v.FavoriteCount
	}
	w.Write([]byte(`{ "Tweets like count":"` + string(count) + `" }`))

	log.Println("User screen name", info[0].User.ScreenName)
	log.Println("Tweets like count:", count)
}

//DisplayAllPowerData :
func DisplayAllPowerData(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	url := "https://api.data.gov.in/resource/1f023275-e8a1-4dc9-a014-4eed45403154?api-key=579b464db66ec23bdd000001cdd3946e44ce4aad7209ff7b23ac571b&format=json&offset=0&limit=10"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("cache-control", "no-cache")
	//	req.Header.Add("Postman-Token", "064f3b39-db1d-4c70-9f9e-3512529a6d9c")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Response error:", err)
	}

	defer res.Body.Close()
	//body, _ := ioutil.ReadAll(res.Body)

	//	fmt.Println(res)
	//	fmt.Println(string(body))

	var data PowerDataStr
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		log.Printf("error decoding response: %v", err)
		if e, ok := err.(*json.SyntaxError); ok {
			log.Printf("syntax error at byte offset %d", e.Offset)
		}

	}
	log.Println("Data displayed:::", data.Title)
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(data)
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return

}

//StoreEnergyData :
func StoreEnergyData(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	url := "https://api.data.gov.in/resource/1f023275-e8a1-4dc9-a014-4eed45403154?api-key=579b464db66ec23bdd000001cdd3946e44ce4aad7209ff7b23ac571b&format=json&offset=0&limit=10"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("cache-control", "no-cache")
	//	req.Header.Add("Postman-Token", "064f3b39-db1d-4c70-9f9e-3512529a6d9c")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Response error:", err)
	}

	defer res.Body.Close()

	var allData PowerDataStr

	err = json.NewDecoder(res.Body).Decode(&allData)

	if err != nil {
		log.Printf("error decoding response: %v", err)
		if e, ok := err.(*json.SyntaxError); ok {
			log.Printf("syntax error at byte offset %d", e.Offset)
		}

	}

	collection := config.Client.Database("power").Collection("data")
	insert, err := collection.InsertOne(r.Context(), allData)
	if err != nil {
		log.Println("Error in inserting power data ::::", err)
	}
	log.Println("Inserted document:::", insert)

}

//Role :
func Role(w http.ResponseWriter, r *http.Request) {
	rbac := gorbac.New()

	r1 := newMyRole("role-1")
	r2 := newMyRole("role-2")

	if err := rbac.Add(r1); err != nil {
		fmt.Printf("Add Error: %s", err)
		return
	}

	if err := rbac.Add(r2); err != nil {
		fmt.Printf("Add Error: %s", err)
		return
	}

	if err := rbac.SetParents("role-1", []string{"role-2"}); err != nil {
		fmt.Printf("SetParents Error: %s", err)
		return
	}

	role, parents, err := rbac.Get("role-1")
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	if myRole, ok := role.(*myRole); ok {
		fmt.Printf("Name:\t%s\nLabel:\t%s\nDesc:\t%s\nParents:\t%s\n",
			myRole.ID(), myRole.Label, myRole.Description, parents)
	}
}

type myRole struct {
	*gorbac.StdRole
	Label       string
	Description string
}

func loadByName(name string) (label, description string) {
	// loading data from storages or somewhere
	return name + " for testing", "This is the description for " + name
}

func newMyRole(name string) gorbac.Role {
	// loading extra properties by `name`.
	label, desc := loadByName(name)
	role := &myRole{
		Label:       label,
		Description: desc,
	}
	role.StdRole = gorbac.NewStdRole(name)
	return role
}
