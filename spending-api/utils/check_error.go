package utils

import (
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/rs/zerolog/log"
)

func CheckError(err error) {
	if err != nil {
		log.Error().Msg(err.Error())
		panic(err)
	}
}

func TraceError(span trace.Span, err error) {
	if err != nil {
		log.Error().Msg(err.Error())
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
}
