package handler

import (
	"bytes"
	"encoding/json"
	"image"
	"image/jpeg"
	"log"
	"net/http"

	"github.com/GregChrisnaDev/Amartha-Sol-3/internal/model"
	"github.com/GregChrisnaDev/Amartha-Sol-3/internal/usecase"
)

func writeJSON(w http.ResponseWriter, statusCode int, message string, errResponse string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := APIResponse{
		Code:    statusCode,
		Message: message,
		Data:    data,
		Error:   errResponse,
	}
	json.NewEncoder(w).Encode(resp)
}

func validateUserAuth(r *http.Request, userUC usecase.UserUsecase, expectedRole int) *model.User {
	username, password, ok := r.BasicAuth()
	if !ok {
		return nil
	}

	user := userUC.ValidateUser(r.Context(), usecase.ValidateUserReq{Email: username, Password: password})
	if user == nil {
		return nil
	}

	if expectedRole != 0 && expectedRole != user.Role {
		return nil
	}

	return user

}

func convertImageToBuffer(r *http.Request, paramName string) (*bytes.Buffer, error) {
	file, _, err := r.FormFile(paramName)
	if err != nil {
		log.Println("[convertImageToBuffer] error get image", err.Error())
		return nil, err
	}
	defer file.Close()

	// Decode the image
	img, _, err := image.Decode(file)
	if err != nil {
		log.Println("[convertImageToBuffer] error decode", err.Error())
		return nil, err
	}

	// Convert to JPEG in memory
	var buf bytes.Buffer
	err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 85})
	if err != nil {
		log.Println("[convertImageToBuffer] error encode", err.Error())
		return nil, err
	}

	return &buf, nil
}
