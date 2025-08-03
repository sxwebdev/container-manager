package logger

import (
	"fmt"
	"log/slog"
	"os"
)

// Logger интерфейс для логирования
type Logger interface {
	Debug(args ...any)
	Debugf(template string, args ...any)
	Debugw(msg string, keysAndValues ...any)

	Info(args ...any)
	Infof(template string, args ...any)
	Infow(msg string, keysAndValues ...any)

	Warn(args ...any)
	Warnf(template string, args ...any)
	Warnw(msg string, keysAndValues ...any)

	Error(args ...any)
	Errorf(template string, args ...any)
	Errorw(msg string, keysAndValues ...any)

	Fatal(args ...any)
	Fatalf(template string, args ...any)
	Fatalw(msg string, keysAndValues ...any)
}

// SlogLogger структура, реализующая интерфейс Logger
type SlogLogger struct {
	logger *slog.Logger
}

// Config конфигурация для логгера
type Config struct {
	Level      slog.Level
	Format     string // "json" или "text"
	ExtraAttrs map[string]any
}

// New создает новый экземпляр логгера
func New(cfg Config) Logger {
	var handler slog.Handler

	opts := &slog.HandlerOptions{
		Level: cfg.Level,
	}

	// Выбираем формат вывода
	switch cfg.Format {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, opts)
	default:
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	// Создаем базовый логгер
	logger := slog.New(handler)

	// Добавляем дополнительные атрибуты, если они заданы
	if len(cfg.ExtraAttrs) > 0 {
		var attrs []any
		for key, value := range cfg.ExtraAttrs {
			attrs = append(attrs, slog.Any(key, value))
		}
		logger = logger.With(attrs...)
	}

	return &SlogLogger{
		logger: logger,
	}
}

// Debug методы
func (l *SlogLogger) Debug(args ...any) {
	l.logger.Debug(fmt.Sprint(args...))
}

func (l *SlogLogger) Debugf(template string, args ...any) {
	l.logger.Debug(fmt.Sprintf(template, args...))
}

func (l *SlogLogger) Debugw(msg string, keysAndValues ...any) {
	l.logger.Debug(msg, l.convertToAttrs(keysAndValues...)...)
}

// Info методы
func (l *SlogLogger) Info(args ...any) {
	l.logger.Info(fmt.Sprint(args...))
}

func (l *SlogLogger) Infof(template string, args ...any) {
	l.logger.Info(fmt.Sprintf(template, args...))
}

func (l *SlogLogger) Infow(msg string, keysAndValues ...any) {
	l.logger.Info(msg, l.convertToAttrs(keysAndValues...)...)
}

// Warn методы
func (l *SlogLogger) Warn(args ...any) {
	l.logger.Warn(fmt.Sprint(args...))
}

func (l *SlogLogger) Warnf(template string, args ...any) {
	l.logger.Warn(fmt.Sprintf(template, args...))
}

func (l *SlogLogger) Warnw(msg string, keysAndValues ...any) {
	l.logger.Warn(msg, l.convertToAttrs(keysAndValues...)...)
}

// Error методы
func (l *SlogLogger) Error(args ...any) {
	l.logger.Error(fmt.Sprint(args...))
}

func (l *SlogLogger) Errorf(template string, args ...any) {
	l.logger.Error(fmt.Sprintf(template, args...))
}

func (l *SlogLogger) Errorw(msg string, keysAndValues ...any) {
	l.logger.Error(msg, l.convertToAttrs(keysAndValues...)...)
}

// Fatal методы
func (l *SlogLogger) Fatal(args ...any) {
	l.logger.Error(fmt.Sprint(args...))
	os.Exit(1)
}

func (l *SlogLogger) Fatalf(template string, args ...any) {
	l.logger.Error(fmt.Sprintf(template, args...))
	os.Exit(1)
}

func (l *SlogLogger) Fatalw(msg string, keysAndValues ...any) {
	l.logger.Error(msg, l.convertToAttrs(keysAndValues...)...)
	os.Exit(1)
}

// convertToAttrs преобразует пары ключ-значение в slog.Attr
func (l *SlogLogger) convertToAttrs(keysAndValues ...any) []any {
	if len(keysAndValues)%2 != 0 {
		// Если нечетное количество аргументов, добавляем последний как "EXTRA"
		keysAndValues = append(keysAndValues, "EXTRA")
	}

	attrs := make([]any, 0, len(keysAndValues))
	for i := 0; i < len(keysAndValues); i += 2 {
		key, ok := keysAndValues[i].(string)
		if !ok {
			key = fmt.Sprintf("%v", keysAndValues[i])
		}
		attrs = append(attrs, slog.Any(key, keysAndValues[i+1]))
	}

	return attrs
}
