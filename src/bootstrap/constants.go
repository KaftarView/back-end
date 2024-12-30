package bootstrap

import (
	"fmt"
)

type Constants struct {
	Context     Context
	ErrorField  ErrorField
	ErrorTag    ErrorTag
	Redis       Redis
	S3Service   S3Service
	JWTKeysPath string
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
	Permission  string
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
	News        string
	Journal     string
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
	NotSubscribe            string
}

type Redis struct {
}

type S3Service struct {
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
			Permission:  "permission",
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
			News:        "news",
			Journal:     "journal",
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
			NotSubscribe:            "notSubscribed",
		},
		Redis:       Redis{},
		JWTKeysPath: "./src/jwtKeys",
	}
}

func (r *Redis) GetUserID(userID int) string {
	return fmt.Sprintf("user:%d", userID)
}

func (s *S3Service) GetEventBannerKey(eventID uint, bannerFilename string) string {
	return fmt.Sprintf("events/%d/banner/%s", eventID, bannerFilename)
}

func (s *S3Service) GetEventSessionKey(eventID, sessionID uint, sessionFilename string) string {
	return fmt.Sprintf("events/%d/sessions/%d/%s", eventID, sessionID, sessionFilename)
}

func (s *S3Service) GetPodcastBannerKey(podcastID uint, bannerFilename string) string {
	return fmt.Sprintf("podcasts/%d/banner/%s", podcastID, bannerFilename)
}

func (s *S3Service) GetPodcastEpisodeBannerKey(podcastID, episodeID uint, bannerFileName string) string {
	return fmt.Sprintf("podcasts/%d/episodes/%d/banner/%s", podcastID, episodeID, bannerFileName)
}

func (s *S3Service) GetPodcastEpisodeKey(podcastID, episodeID uint, filename string) string {
	return fmt.Sprintf("podcasts/%d/episodes/%d/content/%s", podcastID, episodeID, filename)
}

func (s *S3Service) GetNewsBannerKey(newsID uint, bannerFilename string) string {
	return fmt.Sprintf("news/%d/banners/%s", newsID, bannerFilename)
}

func (s *S3Service) GetJournalBannerKey(journalID uint, bannerFilename string) string {
	return fmt.Sprintf("journals/%d/banner/%s", journalID, bannerFilename)
}

func (s *S3Service) GetJournalFileKey(journalID uint, pdfFilename string) string {
	return fmt.Sprintf("journals/%d/content/%s", journalID, pdfFilename)
}

func (s *S3Service) GetOrganizerProfileKey(organizerID uint, profileName string) string {
	return fmt.Sprintf("organizers/%d/profile/%s", organizerID, profileName)
}

func (s *S3Service) GetCouncilorProfileKey(councilorID uint, profileName string) string {
	return fmt.Sprintf("councilors/%d/profile/%s", councilorID, profileName)
}
