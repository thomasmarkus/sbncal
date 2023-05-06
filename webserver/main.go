package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Port string `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"server"`
	Database struct {
		Username string `yaml:"user"`
		Password string `yaml:"pass"`
	} `yaml:"database"`
}

type Wedstrijd struct {
	Datum         string   `json:"Datum"`
	Plaats        string   `json:"Plaats"`
	LSR           bool     `json:"LSR"`
	MSR           bool     `json:"MSR"`
	KSR           bool     `json:"KSR"`
	Jeugd         bool     `json:"Jeugd"`
	BSR           bool     `json:"BSR"`
	Kwalificatie  bool     `json:"Kwalificatie"`
	Afstanden     []string `json:"Afstanden"`
	MinLeeftijd   string   `json:"min.leeftijd"`
	Organisator   []string `json:"Organisator"`
	Inschrijflink []string `json:"Inschrijflink"`
	Uitslag       []string `json:"Uitslag"`
	ONK           bool     `json:"ONK"`
	Koppel        bool     `json:"Koppel"`
	Afgelast      bool     `json:"Afgelast"`
	Coords        []string `json:"coords"`
}

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
	jsonFile, err := os.Open(filename)
	if err != nil {
		return wedstrijden, err
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &wedstrijden)

	return wedstrijden, nil
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
