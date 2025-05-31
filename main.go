package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	// Create a new Echo instance
	e := echo.New()
	// Proxy handler
	e.GET("/bgs/*", func(c echo.Context) error {
		// Extract the path from the original request
		path := c.Param("*")

		// Construct the remote URL
		targetURL := "https://game.aq.com/flash/chardetail/bgs/" + path

		// Make the request to the remote server
		resp, err := http.Get(targetURL)
		if err != nil {
			return c.String(http.StatusBadGateway, "Error reaching target: "+err.Error())
		}
		defer resp.Body.Close()

		// Set the same content-type
		c.Response().Header().Set(echo.HeaderContentType, resp.Header.Get(echo.HeaderContentType))

		// Copy the body to the response
		c.Response().WriteHeader(resp.StatusCode)
		_, err = io.Copy(c.Response(), resp.Body)
		fmt.Println(path)
		return err
	})
	e.GET("/game/*", func(c echo.Context) error {

		// Extract the path from the original request
		path := c.Param("*")

		// Construct the remote URL
		targetURL := "https://game.aq.com/game/" + path

		// Make the request to the remote server
		resp, err := http.Get(targetURL)
		if err != nil {
			return c.String(http.StatusBadGateway, "Error reaching target: "+err.Error())
		}
		defer resp.Body.Close()

		// Set the same content-type
		c.Response().Header().Set(echo.HeaderContentType, resp.Header.Get(echo.HeaderContentType))

		// Copy the body to the response
		c.Response().WriteHeader(resp.StatusCode)
		_, err = io.Copy(c.Response(), resp.Body)
		return err
	})

	// Serve static files from the public directory
	e.Static("/", "public")

	e.POST("/make", func(c echo.Context) error {
		return c.String(200, "makan siang")
	})

	// Add endpoint to generate character data
	e.GET("/generate", func(c echo.Context) error {
		// Get background index from query param
		bgIndex := c.QueryParam("bgIndex")
		if bgIndex == "" {
			bgIndex = "0" // default to Battleon
		}

		// Generate character data
		characterData := map[string]interface{}{
			"username":       "Defo not Artix",
			"level":          99,
			"class":          "Natlan Revenger",
			"accountType":    "Upholder",
			"gold":           "1,000,000",
			"hairColor":      "0x000000",
			"skinColor":      "0xFFCC99",
			"eyeColor":       "0x0000FF",
			"trimColor":      "0xFF0000",
			"baseColor":      "0x0000FF",
			"accessoryColor": "0xFFFF00",
			"gender":         "M",
			"hairFile":       "hair1",
			"hairName":       "Spiky Hair",
			"classFile":      "peasant2_skin.swf",
			"classLink":      "Peasant2",
			"className":      "Not a Staff Class",
			"weaponFile":     "items/swords/sword01.swf",
			"weaponLink":     "Sword01",
			"weaponName":     "Default Sword",
			"capeFile":       "items/capes/wings4.swf",
			"capeLink":       "Wings4",
			"capeName":       "Wings of the Vindicator",
			"helmFile":       "items/helms/J6.swf",
			"helmLink":       "J6helm",
			"helmName":       "J6 Helm",
			"petFile":        "items/pets/wyvernpet.swf",
			"petLink":        "Wyvernpet",
			"petName":        "Miniature Wyvern",
			"bgIndex":        bgIndex,
			"heroTitle":      "Good",
		}

		return c.JSON(200, characterData)
	})

	// Start server
	e.Logger.Fatal(e.Start(":80"))
}
