package utils

import (
	"context"
	"encoding/json"
	"net/http"

	"go.opentelemetry.io/otel"
)

type Validator interface {
	Valid(context context.Context) (err error)
}

func Encode[T any](context context.Context, writer http.ResponseWriter, status int, value T) error {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(context, "Encoding")
	defer span.End()

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)

	err := json.NewEncoder(writer).Encode(value)

	if err != nil {
		TraceError(span, err)
		return err
	}

	return nil
}

func DecodeValid[T Validator](context context.Context, request *http.Request) (T, error) {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(context, "DecodingValid")
	defer span.End()

	var value T
	err := json.NewDecoder(request.Body).Decode(&value)

	if err != nil {
		TraceError(span, err)
		return value, err
	}

	err = value.Valid(context)
	return value, err
}

func Decode[T any](context context.Context, request *http.Request) (T, error) {
	tracer := otel.Tracer("spending-api")
	_, span := tracer.Start(context, "Decoding")
	defer span.End()

	var value T
	err := json.NewDecoder(request.Body).Decode(&value)

	if err != nil {
		TraceError(span, err)
		return value, err
	}

	return value, nil
}
