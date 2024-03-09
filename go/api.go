package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type Result struct {
	ID   int    `json:"id"`
	UUID string `json:"uuid"`
}

func (s *Server) ListUUID(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("")
	ctx, span := tracer.Start(r.Context(), "api.go")
	defer span.End()

	span.SetAttributes(attribute.String("set arbitrary", "string value"))

	results, err := listUUID(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf(`{"message": "%s"}`, err.Error())))
		return
	}

	span.SetAttributes(attribute.Int("result count", len(results)))

	ctx, span = tracer.Start(ctx, "marshal payload")
	defer span.End()

	resp, err := json.Marshal(results)
	if err != nil {
		log.Printf("marshalling into json: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"message": "error with json"}`))
		return
	}

	_, _ = w.Write(resp)
}

func listUUID(ctx context.Context) ([]Result, error) {
	tracer := otel.Tracer("")
	ctx, span := tracer.Start(ctx, "call microservice")
	defer span.End()

	url := "http://java-api:8080/uuid"
	//url := "http://0.0.0.0:8080/uuid"

	httpClient := &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return []Result{}, fmt.Errorf("make request with context: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")

	slog.InfoContext(ctx, "calling java microservice")

	resp, err := httpClient.Do(req)
	if err != nil {
		return []Result{}, fmt.Errorf("perform API call: %w", err)
	}
	defer resp.Body.Close()

	var res []Result

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		slog.InfoContext(ctx, "raw api result", "raw", res)
		slog.ErrorContext(ctx, "error decoding", "raw", err.Error())
		return []Result{}, fmt.Errorf("parse API result to local struct: %w", err)
	}

	slog.InfoContext(ctx, "result stat", "num", len(res))

	return res, nil
}

func (s *Server) AddUUID(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("")
	ctx, span := tracer.Start(r.Context(), "randomAdd")
	defer span.End()

	_, err := s.DB.ExecContext(ctx, "INSERT INTO uuid (uuid) values(uuid_generate_v4()); ")
	if err != nil {
		log.Printf("insert new uuid: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"message": "error with database"}`))
		return
	}
}
