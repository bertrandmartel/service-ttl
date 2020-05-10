package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gopkg.in/go-playground/validator.v9"
)

type Command struct {
	Binary string   `json:"binary"`
	Params []string `json:"params"`
}

type Config struct {
	Version    string    `json:"version"`
	Port       int       `json:"port"`
	ServerPath string    `json:"serverPath"`
	Commands   []Command `json:"commands"`
}

var running = false
var timeoutMinutes = 1
var serviceExp = time.Now().Local().Add(time.Minute * time.Duration(timeoutMinutes))
var cmd *exec.Cmd = nil

func main() {
	portPtr := flag.Int("port", -1, "enter port")
	configFilePathPtr := flag.String("config", "config.json", "enter config file path")
	flag.Parse()
	port := *portPtr
	configFilePath := *configFilePathPtr
	var config Config
	err := ParseConfig(&config, configFilePath)
	if err != nil {
		log.Fatal(err)
		return
	}
	if port != -1 {
		config.Port = port
	}
	go monitorExpiration()
	e := echo.New()
	UseCommonMiddleware(e)
	routes(e, &config)
	e.Logger.Fatal(e.Start(":" + strconv.Itoa(config.Port)))
}

func routes(e *echo.Echo, cfg *Config) {
	e.POST("/start", func(c echo.Context) (err error) {
		serviceExp = time.Now().Local().Add(time.Minute * time.Duration(timeoutMinutes))
		if !running {
			go startService(cfg)
		} else {
			log.Println("service already running")
		}
		return c.NoContent(http.StatusOK)
	})
}

//middleware for validation
type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func UseCommonMiddleware(e *echo.Echo) {
	e.Validator = &CustomValidator{validator: validator.New()}

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${remote_ip} - - ${time_rfc3339_nano} \"${method} ${uri} ${protocol}\" ${status} ${bytes_out} \"${referer}\" \"${user_agent}\"\n",
	}))
	e.Use(middleware.Recover())
}

func monitorExpiration() {
	for {
		if running && time.Now().After(serviceExp) {
			log.Println("stopping process...")
			if cmd != nil {
				if err := cmd.Process.Kill(); err != nil {
					log.Fatal("failed to kill process: ", err)
				}
				log.Println("process killed")
			}
		}
		time.Sleep(10 * time.Second)
	}
}

func startService(config *Config) {
	log.Println("starting service ...")
	running = true
	for i := 0; i < len(config.Commands); i++ {
		cmd = exec.Command(config.Commands[i].Binary, config.Commands[i].Params...)
		err := cmd.Start()
		if err != nil {
			log.Println(err)
			break
		}
		err = cmd.Wait()
		if err != nil {
			log.Println(err)
		}
	}
	running = false
	log.Println("service stopped ...")
}

func ParseConfig(c interface{}, path string) error {
	jsonFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, c)
	return nil
}
