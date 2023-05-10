package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

var wedstrijden []Wedstrijd

func main() {
	//The config file is gitignored, copy from config_defaults.yaml and edit
	config, err := decodeConfig("../config.yaml")
	if err != nil {
		fmt.Print(err)
		return
	}

	//For now, import wedstrijden once from a json file
	wedstrijden, err = importWedstrijden("../wedstrijden.json")
	if err != nil {
		fmt.Print(err)
		return
	}

	router := setupRouter()
	router.Run(fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port))
}

func setupRouter() *gin.Engine {
	router := gin.Default()

	//REST API
	router.GET("/api/wedstrijden", getWedstrijden)

	//Serve just a plain html file as index
	router.LoadHTMLGlob("templates/*.html")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	return router
}

func importWedstrijden(filename string) ([]Wedstrijd, error) {
	var wedstrijden []Wedstrijd
	byteValue, err := os.ReadFile(filename)
	if err != nil {
		return wedstrijden, err
	}

	err = json.Unmarshal(byteValue, &wedstrijden)
	return wedstrijden, err
}

func decodeConfig(filename string) (Config, error) {
	var cfg Config
	f, err := os.Open(filename)
	if err != nil {
		return cfg, err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}

func getWedstrijden(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, wedstrijden)
}
