package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log/slog"
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
		privateIP, err := getPrivateIP()
		if err != nil {
			return c.HTML(http.StatusInternalServerError, "Failed to retrieve private IP address")
		}

		hostname, err := os.Hostname()
		if err != nil {
			return c.HTML(http.StatusInternalServerError, "Failed to retrieve hostname")
		}

		goVersion := runtime.Version()

		appVersion := os.Getenv("VERSION_LABEL")
		if appVersion == "" {
			appVersion = "null"
			slog.Info("VERSION_LABEL environment is not set")
		}

		appEnvironment := os.Getenv("ENVIRONMENT_LABEL")
		if appEnvironment == "" {
			appEnvironment = "null"
			slog.Info("ENVIRONMENT_LABEL environment is not set")
		}

		const templateStr = `
		Private IP: {{.privateIP}}<br>
		Hostname: {{.hostname}}<br>
		Go Version: {{.goVersion}}<br>
		App Version: {{.appVersion}}<br>
		App Environment: {{.appEnvironment}}<br>
		`

		tmpl, err := template.New("info").Parse(templateStr)
		if err != nil {
			return err
		}

		data := map[string]string{
			"privateIP":      privateIP,
			"hostname":       hostname,
			"goVersion":      goVersion,
			"appVersion":     appVersion,
			"appEnvironment": appEnvironment,
		}

		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, data); err != nil {
			return err
		}

		// Return the information in the response
		return c.HTML(http.StatusOK, buf.String())
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
