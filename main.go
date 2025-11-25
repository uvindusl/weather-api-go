package main

import (
	"encoding/json"
	"log"
	"os"
	"io/ioutil"
	"net/http"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

type Response struct {
	Address	string `json:"address"`
	Timezone string `json:"timezone"`
	CurrentConditions Weather `json:"currentConditions"`
}

type Weather struct {
	Temp float64 `json:"temp"`
	Humidity float64 `json:"humidity"`
	Wind float64 `json:"windspeed"`
	Conditions string `json:"conditions"`
}

func getWeather(c *fiber.Ctx) error {
	apiKey := loadEnv("APIKEY")
	location := c.Params("location")
	date := c.Params("date")
	url := fmt.Sprintf("https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline/%s/%s?key=%s", location, date, apiKey)
    response, err := http.Get(url)

    if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch weather data",
			"details": err.Error(),
		})
    }

    responseData, err := ioutil.ReadAll(response.Body)
    if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch weather data",
			"details": err.Error(),
		})
    }
	
	var responseObject Response
	err = json.Unmarshal(responseData, &responseObject)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse weather data",
			"details": err.Error(),
		})
	}

	return c.JSON(responseObject)
}

func loadEnv(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Print("Error loading .env file")
	}
	return os.Getenv(key)
}

func main() {
	app := fiber.New()

	app.Get("/weather/:location/:date" , getWeather)

	log.Fatal(app.Listen(":3000"))
}
