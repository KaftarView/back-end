package bootstrap

import (
	"fmt"
)

type Constants struct {
	Context       Context
	ErrorField    ErrorField
	ErrorTag      ErrorTag
	Redis         Redis
	ObjectStorage Object
}

type Context struct {
	Translator                   string
	IsLoadedValidationTranslator string
	IsLoadedJWTKeys              string
	AccessToken                  string
	RefreshToken                 string
	UserID                       string
}

type ErrorField struct {
	Username    string
	Password    string
	Email       string
	OTP         string
	Tittle      string
	Role        string
	Location    string
	Event       string
	Ticket      string
	Discount    string
	EventStatus string
	Media       string
	Organizer   string
	Token       string
	User        string
	Post        string
	Comment     string
	Podcast     string
	Episode     string
}

type ErrorTag struct {
	AlreadyExist            string
	MinimumLength           string
	ContainsLowercase       string
	ContainsUppercase       string
	ContainsNumber          string
	ContainsSpecialChar     string
	NotMatchConfirmPAssword string
	AlreadyVerified         string
	ExpiredToken            string
	InvalidToken            string
	LoginFailed             string
	EmailNotExist           string
	Required                string
	LocationAlreadyTaken    string
	AlreadySubscribed       string
}

type Redis struct {
}

type Object struct {
}

func NewConstants() *Constants {
	return &Constants{
		Context: Context{
			Translator:                   "translator",
			IsLoadedValidationTranslator: "isLoadedValidationTranslator",
			IsLoadedJWTKeys:              "isLoadedJWTKeys",
			AccessToken:                  "access_token",
			RefreshToken:                 "refresh_token",
			UserID:                       "userID",
		},
		ErrorField: ErrorField{
			Username:    "username",
			Password:    "password",
			Email:       "email",
			OTP:         "OTP",
			Tittle:      "tittle",
			Role:        "role",
			Location:    "location",
			Event:       "event",
			Ticket:      "ticket",
			Discount:    "discount",
			EventStatus: "event status",
			Media:       "media",
			Organizer:   "organizer",
			Token:       "token",
			User:        "user",
			Post:        "post",
			Comment:     "comment",
			Podcast:     "podcast",
			Episode:     "episode",
		},
		ErrorTag: ErrorTag{
			AlreadyExist:            "alreadyExist",
			MinimumLength:           "minimumLength",
			ContainsLowercase:       "containsLowercase",
			ContainsUppercase:       "containsUppercase",
			ContainsNumber:          "containsNumber",
			ContainsSpecialChar:     "containsSpecialChar",
			NotMatchConfirmPAssword: "notMatchConfirmPAssword",
			AlreadyVerified:         "alreadyVerified",
			ExpiredToken:            "expiredToken",
			InvalidToken:            "invalidToken",
			LoginFailed:             "loginFailed",
			EmailNotExist:           "emailNotExist",
			Required:                "required",
			LocationAlreadyTaken:    "locationAlreadyTaken",
			AlreadySubscribed:       "alreadySubscribed",
		},
		Redis: Redis{},
	}
}

func (r *Redis) GetUserID(userID int) string {
	return fmt.Sprintf("user:%d", userID)
}

func (o *Object) GetObjectKey(objectID int, objectTittle, objectName string) string {
	return fmt.Sprintf("%s/%d/%s", objectTittle, objectID, objectName)
}
