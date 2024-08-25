package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"runtime"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func getPrivateIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		// Check if the address is a valid IP address and is not a loopback
		ipNet, ok := addr.(*net.IPNet)
		if ok && !ipNet.IP.IsLoopback() {
			// Check for IPv4 addresses only
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}

	return "", fmt.Errorf("no private IP address found")
}

func main() {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		// Get the private IP address
		privateIP, err := getPrivateIP()
		if err != nil {
			return c.HTML(http.StatusInternalServerError, "Failed to retrieve private IP address")
		}

		// Get the hostname
		hostname, err := os.Hostname()
		if err != nil {
			return c.HTML(http.StatusInternalServerError, "Failed to retrieve hostname")
		}

		// Get the Go version
		goVersion := runtime.Version()

		// Return the information in the response
		response := fmt.Sprintf("Private IP: %s<br>Hostname: %s<br>Go Version: %s", privateIP, hostname, goVersion)
		return c.HTML(http.StatusOK, response)
	})

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct{ Status string }{Status: "OK"})
	})

	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	e.Logger.Fatal(e.Start(":" + httpPort))
}

// Simple implementation of an integer minimum
// Adapted from: https://gobyexample.com/testing-and-benchmarking
func IntMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}
