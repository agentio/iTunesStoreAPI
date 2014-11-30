package iTunesStoreAPI

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Connection struct {
	Country string
	Client  *http.Client
}

var Categories = map[string]string{
	"topfreeapplications":         "Top Free Applications",
	"toppaidapplications":         "Top Paid Applications",
	"topgrossingapplications":     "Top Grossing Applications",
	"topfreeipadapplications":     "Top Free iPad Applications",
	"toppaidipadapplications":     "Top Paid iPad Applications",
	"topgrossingipadapplications": "Top Grossing iPad Applications",
	"newapplications":             "New Applications",
	"newfreeapplications":         "New Free Applications",
	"newpaidapplications":         "New Paid Applications",
}

var Genres = map[int]string{
	0:    "All",
	6000: "Business",
	6001: "Weather",
	6002: "Utilities",
	6003: "Travel",
	6004: "Sports",
	6005: "Social Networking",
	6006: "Reference",
	6007: "Productivity",
	6008: "Photo & Video",
	6009: "News",
	6010: "Navigation",
	6011: "Music",
	6012: "Lifestyle",
	6013: "Health & Fitness",
	6014: "Games",
	6015: "Finance",
	6016: "Entertainment",
	6017: "Education",
	6018: "Books",
	6020: "Medical",
	6021: "Newsstand",
	6022: "Catalogs",
}

type Link struct {
	XMLName xml.Name `xml:"link"`
	Rel     string   `xml:"rel,attr"`
	Type    string   `xml:"type,attr"`
	Href    string   `xml:"href,attr"`
}

type ContentType struct {
	XMLName xml.Name `xml:"contentType"`
	Term    string   `xml:"term,attr"`
	Label   string   `xml:"label,attr"`
}

type Category struct {
	XMLName xml.Name `xml:"category"`
	Id      string   `xml:"id,attr"`
	Term    string   `xml:"term,attr"`
	Scheme  string   `xml:"scheme,attr"`
	Label   string   `xml:"label,attr"`
}

type Artist struct {
	XMLName xml.Name `xml:"artist"`
	Name    string   `xml:",innerxml"`
	Href    string   `xml:"href,attr"`
}

type Price struct {
	XMLName  xml.Name `xml:"price"`
	Amount   string   `xml:"amount,attr"`
	Currency string   `xml:"currency,attr"`
	Price    string   `xml:",innerxml"`
}

type Author struct {
	XMLName xml.Name `xml:"author"`
	Name    string   `xml:"name"`
	Uri     string   `xml:"uri"`
}

type ItemId struct {
	XMLName  xml.Name `xml:"id"`
	Number   string   `xml:"http://itunes.apple.com/rss id,attr"`
	BundleId string   `xml:"http://itunes.apple.com/rss bundleId,attr"`
	Url      string   `xml:",innerxml"`
}

type Image struct {
	XMLName xml.Name `xml:"image"`
	URL     string   `xml:",innerxml"`
	Height  uint16   `xml:"height,attr"`
}

type ReleaseDate struct {
	XMLName xml.Name `xml:"releaseDate"`
	Date    string   `xml:",innerxml"`
	Label   string   `xml:"label,attr"`
}

type Entry struct {
	XMLName     xml.Name    `xml:"entry"`
	Updated     string      `xml:"updated"`
	Id          ItemId      `xml:"id"`
	Title       string      `xml:"title"`
	Summary     string      `xml:"summary"`
	Name        string      `xml:"name"`
	Link        Link        `xml:"link"`
	ContentType ContentType `xml:"contentType"`
	Categories  []Category  `xml:"category"`
	Artists     []Artist    `xml:"artist"`
	Price       Price       `xml:"price"`
	Images      []Image     `xml:"image"`
	Rights      string      `xml:"rights"`
	ReleaseDate ReleaseDate `xml:"releaseDate"`
	Content     string      `xml:"content"`
}

type Feed struct {
	XMLName xml.Name `xml:"feed"`
	Id      string   `xml:"id"`
	Title   string   `xml:"title"`
	Updated string   `xml:"updated"`
	Links   []Link   `xml:"link"`
	Icon    string   `xml:"icon"`
	Author  Author   `xml:"author"`
	Rights  string   `xml:"rights"`
	Entries []Entry  `xml:"entry"`
}

func (connection *Connection) performRequest(req *http.Request) (contents []byte, err error) {
	if connection.Client == nil {
		connection.Client = &http.Client{}
	}
	response, err := connection.Client.Do(req)
	if err != nil {
		return
	}
	contents, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	defer response.Body.Close()

	fmt.Printf("status code %v\n", response.StatusCode)
	fmt.Println(string(contents))

	return contents, err
}

func (self *Connection) requestForArguments(requestURL string) (req http.Request, err error) {
	newreq, err := http.NewRequest("GET", requestURL, nil)
	req = *newreq
	return req, err
}

func (connection *Connection) FetchAppList(category string, genre, limit int) (feed *Feed, err error) {
	link := "https://itunes.apple.com/us/rss/" + category + "/limit=" + strconv.Itoa(limit) + "/genre=" + strconv.Itoa(genre) + "/xml"

	httpRequest, err := http.NewRequest("GET", link, nil)
	contents, err := connection.performRequest(httpRequest)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("%+v\n", string(contents))
	feed = &Feed{}
	err = xml.Unmarshal(contents, &feed)
	return feed, err
}

type Result struct {
	Kind                               string   `json:"kind"`
	Features                           []string `json:"features"`
	SupportedDevices                   []string `json:"supportedDevices"`
	IsGameCenterEnabled                bool     `json:"isGameCenterEnabled"`
	ArtistViewURL                      string   `json:"artistViewURL"`
	ArtworkUrl60                       string   `json:"artworkUrl60"`
	ArtworkUrl100                      string   `json:"artworkUrl100"`
	ArtworkUrl512                      string   `json:"artworkUrl512"`
	ScreenshotURLs                     []string `json:"screenshotURLs"`
	IPadScreenshotURLs                 []string `json:"ipadScreenshotUrls"`
	ArtistId                           int      `json:"artistId"`
	ArtistName                         string   `json:"artistName"`
	Price                              float32  `json:"price"`
	Version                            string   `json:"version"`
	Description                        string   `json:"description"`
	Currency                           string   `json:"currency"`
	Genres                             []string `json:"genres"`
	GenreIds                           []string `json:"genreIds"`
	ReleaseDate                        string   `json:"releaseDate"`
	SellerName                         string   `json:"sellerName"`
	BundleId                           string   `json:"bundleId"`
	TrackId                            int      `json:"trackId"`
	TrackName                          string   `json:"trackName"`
	PrimaryGenreName                   string   `json:"primaryGenreName"`
	PrimaryGenreId                     int      `json:"primaryGenreId"`
	MinimumOsVersion                   string   `json:"minimumOsVersion"`
	FormattedPrice                     string   `json:"formattedPrice"`
	WrapperType                        string   `json:"wrapperType"`
	TrackCensoredName                  string   `json:"trackCensoredName"`
	TrackViewUrl                       string   `json:"trackViewUrl"`
	ContentAdvisoryRating              string   `json:"contentAdvisoryRating"`
	LanguageCodesISO2A                 []string `json:"languageCodesISO2A"`
	FileSizeBytes                      string   `json:"fileSizeBytes"`
	SellerUrl                          string   `json:"sellerUrl"`
	AverageUserRatingForCurrentVersion float32  `json:"averageUserRatingForCurrentVersion"`
	UserRatingCountForCurrentVersion   int      `json:"userRatingCountForCurrentVersion"`
	TrackContentRating                 string   `json:"trackContentRating"`
	AverageUserRating                  float32  `json:"averageUserRating"`
	UserRatingCount                    int      `json:"userRatingCount"`
}

type ResultSet struct {
	ResultCount int      `json:"resultCount"`
	Results     []Result `json:"results"`
}

func (connection *Connection) LookupItemWithId(appid string) (results *ResultSet, err error) {
	link := "http://itunes.apple.com/lookup?id=" + appid

	httpRequest, err := http.NewRequest("GET", link, nil)
	contents, err := connection.performRequest(httpRequest)
	if err != nil {
		return nil, err
	}
	info := &ResultSet{}
	err = json.Unmarshal(contents, &info)
	if err != nil {
		return nil, err
	}
	return info, err
}
