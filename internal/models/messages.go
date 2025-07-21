package models

type AnalyticsMessage struct {
	Detail      string                `json:"detail"`
	TimeToCrash string                `json:"time_to_crash"`
	Payload     AnalyticsPayloadField `json:"payload"`
}

type AnalyticsPayloadField struct {
	UserID    string                 `json:"user_id"`
	UserEmail string                 `json:"user_email"`
	Location  AnalyticsLocationField `json:"location"`
}

type AnalyticsLocationField struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type RecommendationsMessage struct {
	AnalyticsMessage
	Recomendations RecommendationsField `json:"recomendations"`
}

type RecommendationsField struct {
	Market     string `json:"market"`
	Nearest_TO string `json:"nearest_TO"`
}
