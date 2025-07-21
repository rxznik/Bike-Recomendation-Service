package velostrana

import (
	"context"
	"errors"
	"net/http"
)

const VelostranaBaseURL = "https://www.velostrana.ru"

var detailNameToPath = map[string]string{
	"задние переключатели":   "pereklu4enie/perekluchateli-zadnie",
	"передние переключатели": "pereklu4enie/perekluchateli-perednie",
	"велопетух":              "pereklu4enie/petuhi-na-velosiped",
	"педали":                 "pedali",
	"летняя резина":          "pokryshki-i-kamery/letnyaya-rezina",
	"зимняя резина":          "pokryshki-i-kamery/zimnie-pokryshki",
	"камера":                 "pokryshki-i-kamery/kameri",
	"трещотки":               "privod/treschotki",
	"цепь":                   "privod/cepi",
	"замок цепи":             "privod/zamki",
	"система шатунов":        "privod/sistema",
	"каретка":                "privod/karetki",
	"бонки":                  "privod/bonki",
	"натяжитель цепи":        "privod/natyagiteli-cepi",
	"звезда":                 "privod/zvezdi",
	"кассета":                "kasseti-treshotki",
	"рама":                   "rami",
}

func convertFromDetailNameToPath(detailName string) (string, error) {
	if path, ok := detailNameToPath[detailName]; ok {
		return path, nil
	}
	return "", errors.New("detail name not found")
}

func CreateURLToListDetails(detailName string) (string, error) {
	path, err := convertFromDetailNameToPath(detailName)
	if err != nil {
		return "", err
	}
	return VelostranaBaseURL + "/velozapchasti/" + path, nil
}

func CreateVelostranaRequest(ctx context.Context, detailName string) (*http.Request, error) {
	url, err := CreateURLToListDetails(detailName)
	if err != nil {
		return nil, ErrFailedToCreateURLToListDetails
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)

	if err != nil {
		return nil, ErrCreateVelostranaRequest
	}
	return req, nil
}

func AddBaseURL(path string) string {
	if path[0] != '/' {
		path = "/" + path
	}
	return VelostranaBaseURL + path
}
