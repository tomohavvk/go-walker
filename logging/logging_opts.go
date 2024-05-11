package logging

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
)

type LogHandler struct {
	slog.Handler
	l *log.Logger
}

const whiteColor = "\x1b[37m%s\x1b[0m"
const magentaColor = "\x1b[35m%s\x1b[0m"
const blueColor = "\x1b[34m%s\x1b[0m"
const yellowColor = "\x1b[33m%s\x1b[0m"
const redColor = "\x1b[31m%s\x1b[0m"
const cyanColor = "\x1b[36m%s\x1b[0m"
const greenColor = "\x1b[32m%s\x1b[0m"

func (h *LogHandler) Handle(ctx context.Context, r slog.Record) error {
	level := r.Level.String() + ":"

	switch r.Level {
	case slog.LevelDebug:
		level = fmt.Sprintf(magentaColor, level)
	case slog.LevelInfo:
		level = fmt.Sprintf(blueColor, level)
	case slog.LevelWarn:
		level = fmt.Sprintf(yellowColor, level)
	case slog.LevelError:
		level = fmt.Sprintf(redColor, level)
	}

	fields := make(map[string]interface{}, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()

		return true
	})

	jsonBytes, err := json.Marshal(fields)
	if err != nil {
		return err
	}

	jsonStr := string(jsonBytes)

	if jsonStr == "{}" {
		jsonStr = ""
	}

	timeStr := r.Time.Format("[2006-01-02 15:04:05]")
	msg := fmt.Sprintf(cyanColor, r.Message)

	h.l.Println(fmt.Sprintf(greenColor, timeStr), level, msg, fmt.Sprintf(whiteColor, jsonStr))

	return nil
}

func NewLogger(out io.Writer, level slog.Level) slog.Logger {
	handler := &LogHandler{
		Handler: slog.NewJSONHandler(out, &slog.HandlerOptions{Level: level, AddSource: true}),
		l:       log.New(out, "", 0)}

	return *slog.New(handler)
}
