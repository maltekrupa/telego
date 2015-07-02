package telego

import (
	"encoding/json"
	"errors"
	"net/http"
)

var (
	APIURL = "https://api.telegram.org/bot"
)

type telego struct {
	token  string
	apiUrl string
}

func NewTelego(token string) *telego {
	t := new(telego)
	t.token = token
	t.apiUrl = APIURL
	return t
}

func (t *telego) ChangeUrl(url string) {
	t.apiUrl = url + "/bot"
}

type Update struct {
	Id      int64   `json:"update_id"`
	Message Message `json:"message"`
}

type ResponseUpdate struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Me struct {
	Id         int64  `json:"id"`
	First_name string `json:"first_name"`
	Username   string `json:"username"`
}

type ResponseMe struct {
	Ok     bool `json:"ok"`
	Result Me   `json:"result"`
}

type User struct {
	Id         int64  `json:"id"`
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
	Username   string `json:"username"`
}

type GroupChat struct {
	Id    int64  `json:"id"`
	Title string `json:"title"`
}

type PhotoSize struct {
	File_id   string `json:"file_id"`
	Width     int64  `json:"width"`
	Height    int64  `json:"height"`
	File_size int64  `json:"file_size"`
}

type Audio struct {
	File_id   string `json:"file_id"`
	Duration  int64  `json:"duration"`
	Mime_type string `json:"mime_type"`
	File_size int64  `json:"file_size"`
}

type Document struct {
	File_id   string    `json:"file_id"`
	Thumb     PhotoSize `json:"thumb"`
	File_name string    `json:"file_name"`
	Mime_type string    `json:"mime_type"`
	File_size int64     `json:"file_size"`
}

type Sticker struct {
	File_id   string    `json:"file_id"`
	Width     int64     `json:"width"`
	Height    int64     `json:"height"`
	Thumb     PhotoSize `json:"thumb"`
	File_size int64     `json:"file_size"`
}

type Video struct {
	File_id   string    `json:"file_id"`
	Width     int64     `json:"width"`
	Height    int64     `json:"height"`
	Duration  int64     `json:"duration"`
	Thumb     PhotoSize `json:"thumb"`
	Mime_type string    `json:"mime_type"`
	File_size int64     `json:"file_size"`
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
	Total_count int64         `json:"total_count"`
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
	Message_id int64  `json:"message_id"`
	From       User   `json:"from"`
	Date       int64  `json:"date"`
	Chat       User   `json:"chat"`
	Text       string `json:"text"`
}

// Gets the update stream for the bot.
// Currently this only gets
func (t telego) GetUpdates() (ResponseUpdate, error) {
	var response ResponseUpdate

	updateUrl := t.apiUrl + t.token + "/getUpdates"

	resp, err := http.Get(updateUrl)
	if err != nil {
		return ResponseUpdate{}, err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	dec.Decode(&response)
	return response, nil
}

func (t telego) GetMe() (ResponseMe, error) {
	var response ResponseMe

	meUrl := t.apiUrl + t.token + "/getMe"

	resp, err := http.Get(meUrl)
	if err != nil {
		return ResponseMe{}, err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	dec.Decode(&response)
	return response, nil
}

func (t telego) GetMessageFromId(id int64) (Message, error) {
	updates, err := t.GetUpdates()
	if err != nil {
		return Message{}, err
	}

	for _, v := range updates.Result {
		if v.Message.Message_id == id {
			return v.Message, nil
		}
	}
	return Message{}, errors.New("No message found")
}

func (t telego) GetLastMessage() (Message, error) {
	updates, err := t.GetUpdates()
	if err != nil {
		return Message{}, err
	}
	return updates.Result[len(updates.Result)-1].Message, nil
}
