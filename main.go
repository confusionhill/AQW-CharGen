package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"

	"github.com/labstack/echo/v4"
)

func openBrowser(url string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("cmd", "/c", "start", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	return err
}

func main() {
	// Create a new Echo instance
	e := echo.New()

	// Root endpoint that serves HTML with generated data
	e.GET("/", func(c echo.Context) error {
		// Get all query parameters
		params := c.QueryParams()

		// Generate character data with default values
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
			"weaponType":     "Sword",
			"capeFile":       "items/capes/wings4.swf",
			"capeLink":       "Wings4",
			"capeName":       "Wings of the Vindicator",
			"helmFile":       "items/helms/J6.swf",
			"helmLink":       "J6helm",
			"helmName":       "J6 Helm",
			"petFile":        "items/pets/wyvernpet.swf",
			"petLink":        "Wyvernpet",
			"petName":        "Miniature Wyvern",
			"bgIndex":        "0",
			"heroTitle":      "Good",
		}

		// Update character data with query parameters
		for key, values := range params {
			if len(values) > 0 {
				// Convert numeric values
				if key == "level" {
					if val, err := strconv.Atoi(values[0]); err == nil {
						characterData[key] = val
					}
				} else {
					characterData[key] = values[0]
				}
			}
		}

		// Parse and execute the template
		tmpl, err := template.ParseFiles("public/index.html")
		if err != nil {
			return c.String(http.StatusInternalServerError, "Error parsing template: "+err.Error())
		}

		return tmpl.Execute(c.Response().Writer, characterData)
	})

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

	// Start server in a goroutine
	go func() {
		if err := e.Start(":80"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Open browser
	if err := openBrowser("http://localhost:80"); err != nil {
		e.Logger.Printf("Failed to open browser: %v", err)
	}

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Gracefully shutdown the server
	if err := e.Shutdown(nil); err != nil {
		e.Logger.Fatal(err)
	}
}
