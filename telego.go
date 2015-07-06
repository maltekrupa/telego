package telego

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

type telego struct {
	token          string
	apiUrl         string
	url            string
	offset         int
	lastSeenOffset int
}

type Update struct {
	Id      int     `json:"update_id"`
	Message Message `json:"message"`
}

type ResponseUpdate struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Me struct {
	Id         int    `json:"id"`
	First_name string `json:"first_name"`
	Username   string `json:"username"`
}

type ResponseMe struct {
	Ok     bool `json:"ok"`
	Result Me   `json:"result"`
}

type User struct {
	Id         int    `json:"id"`
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
	Username   string `json:"username"`
	Title      string `json:"title"`
}

type GroupChat struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
}

type PhotoSize struct {
	File_id   string `json:"file_id"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	File_size int    `json:"file_size"`
}

type Audio struct {
	File_id   string `json:"file_id"`
	Duration  int    `json:"duration"`
	Mime_type string `json:"mime_type"`
	File_size int    `json:"file_size"`
}

type Document struct {
	File_id   string    `json:"file_id"`
	Thumb     PhotoSize `json:"thumb"`
	File_name string    `json:"file_name"`
	Mime_type string    `json:"mime_type"`
	File_size int       `json:"file_size"`
}

type Sticker struct {
	File_id   string    `json:"file_id"`
	Width     int       `json:"width"`
	Height    int       `json:"height"`
	Thumb     PhotoSize `json:"thumb"`
	File_size int       `json:"file_size"`
}

type Video struct {
	File_id   string    `json:"file_id"`
	Width     int       `json:"width"`
	Height    int       `json:"height"`
	Duration  int       `json:"duration"`
	Thumb     PhotoSize `json:"thumb"`
	Mime_type string    `json:"mime_type"`
	File_size int       `json:"file_size"`
	Caption   string    `json:"caption"`
}

type Contact struct {
	Phone_number string `json:"phone_number"`
	First_name   string `json:"first_name"`
	Last_name    string `json:"last_name"`
	User_id      string `json:"user_id"`
}

type Location struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type UserProfilePhotos struct {
	Total_count int           `json:"total_count"`
	Photos      [][]PhotoSize `json:"photos"`
}

type ReplyKeyboardMarkup struct {
	Keyboard          [][]string `json:"keyboard"`
	Resize_keyboard   bool       `json:"resize_keyboard"`
	One_time_keyboard bool       `json:"one_time_keyboard"`
	Selective         bool       `json:"selective"`
}

type ReplyKeyboardHide struct {
	Hide_keyboard bool `json:"hide_keyboard"`
	Selective     bool `json:"selective"`
}

type ForceReply struct {
	Force_reply bool `json:"force_reply"`
	Selective   bool `json:"selective"`
}

type Message struct {
	Message_id int    `json:"message_id"`
	From       User   `json:"from"`
	Date       int    `json:"date"`
	Chat       User   `json:"chat"`
	Text       string `json:"text"`
}

type ResponseSendMessage struct {
	Ok     bool    `json:"ok"`
	Result Message `json:"result"`
}

var (
	APIURL = "https://api.telegram.org/bot"
)

func NewTelego(token string) *telego {
	t := new(telego)
	t.token = token
	t.apiUrl = APIURL
	t.url = t.apiUrl + t.token
	return t
}

func (t *telego) ChangeUrl(url string) {
	t.apiUrl = url + "/bot"
	t.url = t.apiUrl + t.token
}

func (t telego) SendMessage(id int, text string) (rsm *ResponseSendMessage, err error) {
	url := t.url + "/sendMessage?chat_id=" + strconv.Itoa(id) + "&text=" + text
	url2 := strings.Replace(url, " ", "%20", -1)

	resp, err := http.Get(url2)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	dec.Decode(&rsm)
	return rsm, nil
}

// Gets the update stream for the bot.
// TODO: Enhance with parameters.
func (t *telego) GetUpdates() (ResponseUpdate, error) {
	var response ResponseUpdate
	var url string

	if t.offset != 0 {
		url = t.url + "/getUpdates?offset=" + strconv.Itoa(t.offset+1)
	} else {
		url = t.url + "/getUpdates"
	}

	resp, err := http.Get(url)
	if err != nil {
		return ResponseUpdate{}, err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	dec.Decode(&response)

	if len(response.Result) >= 1 {
		t.lastSeenOffset = response.Result[len(response.Result)-1].Id
	}
	return response, nil
}

func (t telego) GetMe() (ResponseMe, error) {
	var response ResponseMe

	meUrl := t.url + "/getMe"

	resp, err := http.Get(meUrl)
	if err != nil {
		return ResponseMe{}, err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	dec.Decode(&response)
	return response, nil
}

func (t telego) GetLastMessage() (Message, error) {
	updates, err := t.GetUpdates()
	if err != nil {
		return Message{}, err
	}
	if len(updates.Result) < 1 {
		return Message{}, errors.New("No updates available")
	}
	return updates.Result[len(updates.Result)-1].Message, nil
}

func (t telego) GetOffset() int {
	return t.offset
}

func (t *telego) UpdateOffset() {
	t.offset = t.lastSeenOffset
}
