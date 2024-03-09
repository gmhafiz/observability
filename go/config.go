package main

import (
	"os"
	"strconv"
)

type Cfg struct {
	Api
	Database
	OpenTelemetry
}

func Config() Cfg {
	api := configApi()
	db := configDB()
	otel := configOtel()

	return Cfg{
		api,
		db,
		otel,
	}
}

type Api struct {
	Name string
	Host string
	Port int
}

func configApi() Api {
	apiName := os.Getenv("API_NAME")
	if apiName == "" {
		apiName = "Go_API"
	}

	apiHost := os.Getenv("API_HOST")
	apiPort, err := strconv.Atoi(os.Getenv("API_PORT"))
	if err != nil {
		apiPort = 3080
	}

	return Api{
		Name: apiName,
		Host: apiHost,
		Port: apiPort,
	}
}

type OpenTelemetry struct {
	Enable             bool
	OtlpEndpoint       string
	OtlpServiceName    string
	OtlpServiceVersion string
	OtlpMeterName      string
	OtlpSamplerRatio   float64
}

func configOtel() OpenTelemetry {
	otlpEndpoint := os.Getenv("OTEL_GO_ENDPOINT")
	if otlpEndpoint == "" {
		otlpEndpoint = "otel-collector:4317"
	}
	OtlpServiceName := os.Getenv("OTEL_GO_SERVICE_NAME")
	if OtlpServiceName == "" {
		OtlpServiceName = "go_api"
	}
	OtlpServiceVersion := os.Getenv("OTEL_GO_SERVICE_VERSION")
	if OtlpServiceVersion == "" {
		OtlpServiceVersion = "0.1.0"
	}
	OtlpMeterName := os.Getenv("OTEL_GO_METER_NAME")
	if OtlpMeterName == "" {
		OtlpMeterName = "go-meter"
	}
	OtlpSSamplerRatio, _ := strconv.ParseFloat(os.Getenv("OTEL_GO_SAMPLER_RATIO"), 64)
	if OtlpSSamplerRatio == 0 {
		OtlpSSamplerRatio = 0.1
	}

	return OpenTelemetry{
		OtlpEndpoint:       otlpEndpoint,
		OtlpServiceName:    OtlpServiceName,
		OtlpServiceVersion: OtlpServiceVersion,
		OtlpMeterName:      OtlpMeterName,
		OtlpSamplerRatio:   OtlpSSamplerRatio,
	}
}

type Database struct {
	Host    string
	Port    int
	Name    string
	User    string
	Pass    string
	SslMode string
}

func configDB() Database {
	dbHost := os.Getenv("DB_HOST")
	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		dbPort = 5432
	}
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbSSLMode := os.Getenv("DB_SSL_MODE")
	if dbSSLMode == "" {
		dbSSLMode = "disable"
	}

	return Database{
		Host:    dbHost,
		Port:    dbPort,
		Name:    dbName,
		User:    dbUser,
		Pass:    dbPass,
		SslMode: dbSSLMode,
	}
}
