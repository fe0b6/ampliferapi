package ampliferapi

import "encoding/json"

// API - объкт api
type API struct {
	AccessToken string
}

type rqAns struct {
	OK     bool            `json:"ok"`
	Result json.RawMessage `json:"result"`
}

// GetProjectsAns - объект проектов
type GetProjectsAns struct {
	Projects []Project `json:"projects"`
}

// GetProjectPostsAns - Объект постов
type GetProjectPostsAns struct {
	Posts      []Post         `json:"posts"`
	Pagination map[string]int `json:"pagination"`
}

// Project - объект проекта
type Project struct {
	ID             int64     `json:"id"`
	Name           string    `json:"name"`
	SocialAccounts []Account `json:"socialAccounts"`
}

// Account - Объект аккаунта
type Account struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	Avatar      string `json:"avatar"`
	Network     string `json:"network"`
	NetworkAbbr string `json:"networkAbbr"`
	Active      bool   `json:"active"`
	Publishable bool   `json:"publishable"`
}

// Post - объект поста
type Post struct {
	ID             int64             `json:"id"`
	Clicks         int               `json:"clicks"`
	Likes          int               `json:"likes"`
	Shares         int               `json:"shares"`
	Comments       int               `json:"comments"`
	UniqueViews    int               `json:"uniqueViews"`
	FanUniqueViews int               `json:"fanUniqueViews"`
	TotalViews     int               `json:"totalViews"`
	VideoPlays     int               `json:"videoPlays"`
	States         map[string]string `json:"states"`
	Publications   map[string]string `json:"publications"`
}
