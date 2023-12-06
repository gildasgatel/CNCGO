package grbl_test

import (
	"testing"

	connection "cncgo/backend/pkg/connection/mock"
	"cncgo/backend/pkg/machine/grbl"
)

func TestGrblIntegration(t *testing.T) {
	// Initialisez votre connexion factice (mock)
	mockConnection := &connection.MockService{}

	// Créez une nouvelle instance de Grbl avec votre connexion factice
	_, err := grbl.New(mockConnection)
	if err != nil {
		t.Fatalf("Erreur lors de la création de Grbl: %v", err)
	}

	// Utilisez maintenant les méthodes de Grbl et vérifiez le comportement en interaction avec la connexion.
	// Par exemple, appelez SendCommand ou SendFile et vérifiez si les appels aux méthodes de la connexion se font comme prévu.
}
