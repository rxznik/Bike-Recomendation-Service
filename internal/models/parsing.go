package models

type ParsingInput struct {
	Detail string `json:"detail"`
}

type ParsingOutput struct {
	Market string `json:"market"`
}

type ParsingParseResponse struct {
	URL string `json:"url"`
}
