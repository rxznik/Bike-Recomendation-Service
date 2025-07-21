package mockproducer

import "github.com/devprod-tech/webike_recomendations-Vitalya/internal/models"

var mockMessages = []models.AnalyticsMessage{
	{
		Detail:      "педали",
		TimeToCrash: "2026-01-01 00:00:00",
		Payload: models.AnalyticsPayloadField{
			UserID:    "iuwdh98217uuhdqo",
			UserEmail: "pF5kS@example.com",
			Location:  models.AnalyticsLocationField{Latitude: 55.638248, Longitude: 37.612877},
		},
	},
	{
		Detail:      "камера",
		TimeToCrash: "2026-02-02 00:00:00",
		Payload: models.AnalyticsPayloadField{
			UserID:    "bkjtrpobktrfksdlf",
			UserEmail: "0fGtQ@example.com",
			Location:  models.AnalyticsLocationField{Latitude: 55.598285, Longitude: 37.041622},
		},
	},
	{
		Detail:      "цепь",
		TimeToCrash: "2026-03-03 00:00:00",
		Payload: models.AnalyticsPayloadField{
			UserID:    "bgfoppbgkfmlmlkml",
			UserEmail: "bHc8o@example.com",
			Location:  models.AnalyticsLocationField{Latitude: 55.636674, Longitude: 37.215899},
		},
	},
	{
		Detail:      "каретка",
		TimeToCrash: "2026-04-04 00:00:00",
		Payload: models.AnalyticsPayloadField{
			UserID:    "vr98e0840ghfjnvml",
			UserEmail: "z9TtQ@example.com",
			Location:  models.AnalyticsLocationField{Latitude: 55.669959, Longitude: 37.283061},
		},
	},
	{
		Detail:      "рама",
		TimeToCrash: "2026-05-05 00:00:00",
		Payload: models.AnalyticsPayloadField{
			UserID:    "923fujnmvmojkmlkml",
			UserEmail: "gVn8t@example.com",
			Location:  models.AnalyticsLocationField{Latitude: 55.834225, Longitude: 37.353184},
		},
	},
}
