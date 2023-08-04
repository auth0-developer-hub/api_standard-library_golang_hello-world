package messages

import "github.com/auth0-developer-hub/api_standard-library_golang_hello-world/pkg/models"

func PublicMessage() models.ApiResponse {
	return models.ApiResponse{
		Text: "This is a public message.",
	}
}

func ProtectedMessage() models.ApiResponse {
	return models.ApiResponse{
		Text: "This is a protected message.",
	}
}

func AdminMessage() models.ApiResponse {
	return models.ApiResponse{
		Text: "This is an admin message.",
	}
}
