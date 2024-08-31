package main

import (
	"fmt"
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

		releaseVersion := os.Getenv("RELEASE_VERSION")
		if releaseVersion == "" {
			releaseVersion = "null"
			slog.Info("RELEASE_VERSION environment is not set")
		}

		manifestVersion := os.Getenv("MANIFEST_VERSION_LABEL")
		if manifestVersion == "" {
			manifestVersion = "null"
			slog.Info("MANIFEST_VERSION_LABEL environment is not set")
		}

		appEnvironment := os.Getenv("ENVIRONMENT_LABEL")
		if appEnvironment == "" {
			appEnvironment = "null"
			slog.Info("ENVIRONMENT_LABEL environment is not set")
		}

		const templateStr = `
		{
		  "private_ip": "{{.privateIP}}",
		  "hostname": "{{.hostname}}",
		  "go_version": "{{.goVersion}}",
		  "release_version": "{{.releaseVersion}}",
		  "manifest_version": "{{.manifestVersion}}",
		  "app_environment": "{{.appEnvironment}}"
		}
		`

		// Return the information as JSON in the response
		return c.JSON(http.StatusOK, struct {
			PrivateIP       string `json:"private_ip"`
			Hostname        string `json:"hostname"`
			GoVersion       string `json:"go_version"`
			ReleaseVersion  string `json:"release_version"`
			ManifestVersion string `json:"manifest_version"`
			AppEnvironment  string `json:"app_environment"`
		}{
			PrivateIP:       privateIP,
			Hostname:        hostname,
			GoVersion:       goVersion,
			ReleaseVersion:  releaseVersion,
			ManifestVersion: manifestVersion,
			AppEnvironment:  appEnvironment,
		})
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
