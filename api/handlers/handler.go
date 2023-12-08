package handlers

import (
	"cncgo/api/models"
	"cncgo/pkg/connection"
	mockConnection "cncgo/pkg/connection/mock"
	"cncgo/pkg/connection/usb"
	"cncgo/pkg/machine"
	"cncgo/pkg/machine/grbl"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CncHandler struct {
	cnc        machine.Service
	connection connection.Service
}

func (ch *CncHandler) hello(c *gin.Context) {
	c.String(http.StatusOK, "Hello depuis le serveur Gin! mode : %s connexion : %s", ch.cnc.GetName(), ch.connection.GetName())
}

func (ch *CncHandler) state(c *gin.Context) {
	if s := ch.cnc.GetState(); s != nil {
		c.JSON(http.StatusAccepted, s)
	}
}

func (ch *CncHandler) configConnexion(c *gin.Context) {
	var conf models.Config
	// Vérifie si le corps de la requête peut être décodé en JSON
	if err := c.ShouldBindJSON(&conf); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Initialise le service
	err := ch.initConnection(&conf)
	if err != nil {
		c.JSON(400, gin.H{"error_connection": err.Error()})
	}
	err = ch.initMachine(&conf)
	if err != nil {
		c.JSON(400, gin.H{"error_machine": err.Error()})
	} else {
		c.JSON(200, gin.H{"message": "Configuration done!"})
	}
}
func (ch *CncHandler) initConnection(conf *models.Config) (err error) {
	switch conf.Connection {
	case "usb":
		ch.connection, err = usb.New(conf)
	case "test":
		ch.connection = &mockConnection.MockService{}
	default:
		//todo add error
	}
	return
}

func (ch *CncHandler) initMachine(conf *models.Config) (err error) {
	switch conf.Machine {
	case "grbl":
		if ch.cnc, err = grbl.New(ch.connection); err != nil {
			return
		}
	default:
		//todo add error
	}
	return
}

func (ch *CncHandler) handelCommand(c *gin.Context) {
	var com models.Command
	if err := c.ShouldBindJSON(&com); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if ch.cnc != nil {
		resp, err := ch.cnc.SendCommand(com)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"message": string(resp)})
		}
	} else {
		ch.errorConfig(c)
	}

}

func (ch *CncHandler) handelFile(c *gin.Context) {
	var file models.File
	if err := c.ShouldBindJSON(&file); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if ch.cnc != nil {
		c.JSON(http.StatusOK, gin.H{"message": "SendFile() started"})
		err := ch.cnc.SendFile(file.Path)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "SendFile stop"})
		}
	} else {
		ch.errorConfig(c)
	}

}

func (ch *CncHandler) Close() {
	ch.connection.Close()
}

func (ch *CncHandler) errorConfig(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{"message": "Need to setup your config, /config"})
}
