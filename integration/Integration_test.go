package integration

import (
	"bytes"
	"cncgo/api/handlers"
	"cncgo/api/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAPIIntegration(t *testing.T) {
	router := gin.Default()

	var hand handlers.CncHandler
	hand.SetupRouter(router)

	requestBody, _ := json.Marshal(models.Config{
		Machine:    "grbl",
		Connection: "usb",
		BaudRate:   115200,
		PortName:   "/dev/ttyACM0",
	})

	req, err := http.NewRequest("POST", "/config", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatalf("Erreur lors de la création de la requête : %v", err)
	}

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusBadRequest {
		t.Errorf("Code de statut incorrect: got %v mais attendu %v",
			status, http.StatusCreated)
	}

	expectedResponse := `{"message": "User created successfully"}`
	if recorder.Body.String() != expectedResponse {
		t.Errorf("Réponse incorrecte: got %v mais attendu %v",
			recorder.Body.String(), expectedResponse)
	}
}
