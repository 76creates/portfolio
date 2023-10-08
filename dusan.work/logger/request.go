package logger

import (
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"time"
)

type RequestLogger struct {
	*Logger
}

// NewRequest - create new request logger from logger
// request logger will track request duration and will log it on request end
func (l *Logger) NewRequest(id uuid.UUID) *RequestLogger {
	request := RequestLogger{Logger: l.Child()}
	request.extraField = append(
		request.extraField,
		&RequestTimeField{time.Now()},
		&RequestIdField{id},
	)
	return &request
}

type RequestTimeField struct {
	startTime time.Time
}

func (r *RequestTimeField) Field(e *zerolog.Event) *zerolog.Event {
	return e.Str("t", time.Since(r.startTime).String())
}

type RequestIdField struct {
	id uuid.UUID
}

func (r *RequestIdField) Field(e *zerolog.Event) *zerolog.Event {
	return e.Str("id", r.id.String())
}
