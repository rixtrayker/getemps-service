package logger

import (
	"context"
	"encoding/json"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type CustomLogger struct {
	db     *sqlx.DB
	logger *logrus.Logger
	toDB   bool
}

func NewCustomLogger(db *sqlx.DB, logLevel string, toDB bool) *CustomLogger {
	logger := logrus.New()
	
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	return &CustomLogger{
		db:     db,
		logger: logger,
		toDB:   toDB,
	}
}

func (l *CustomLogger) Info(ctx context.Context, message string, fields map[string]interface{}) {
	l.logger.WithFields(fields).Info(message)
	
	if l.toDB {
		l.logToDB("INFO", message, fields, ctx)
	}
}

func (l *CustomLogger) Error(ctx context.Context, message string, err error, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	if err != nil {
		fields["error"] = err.Error()
	}
	
	l.logger.WithFields(fields).Error(message)
	
	if l.toDB {
		l.logToDB("ERROR", message, fields, ctx)
	}
}

func (l *CustomLogger) Warn(ctx context.Context, message string, fields map[string]interface{}) {
	l.logger.WithFields(fields).Warn(message)
	
	if l.toDB {
		l.logToDB("WARN", message, fields, ctx)
	}
}

func (l *CustomLogger) Debug(ctx context.Context, message string, fields map[string]interface{}) {
	l.logger.WithFields(fields).Debug(message)
	
	if l.toDB {
		l.logToDB("DEBUG", message, fields, ctx)
	}
}

func (l *CustomLogger) logToDB(level, message string, fields map[string]interface{}, ctx context.Context) {
	if l.db == nil {
		return
	}

	var contextJSON []byte
	if fields != nil {
		contextJSON, _ = json.Marshal(fields)
	}

	// Extract endpoint and request ID from context if available
	var endpoint, requestID string
	if ctx != nil {
		if ep := ctx.Value("endpoint"); ep != nil {
			endpoint = ep.(string)
		}
		if rid := ctx.Value("request_id"); rid != nil {
			requestID = rid.(string)
		}
	}

	query := `
		INSERT INTO logs (level, message, context, endpoint, request_id, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
	`

	// Execute asynchronously to avoid blocking the main request
	go func() {
		_, err := l.db.Exec(query, level, message, contextJSON, endpoint, requestID)
		if err != nil {
			l.logger.WithError(err).Error("Failed to write log to database")
		}
	}()
}