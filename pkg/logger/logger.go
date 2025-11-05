package logger

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

// Log é o logger principal do pacote (zerolog.Logger)
var Log zerolog.Logger

// Init inicializa o zerolog com o nível informado. Usa ConsoleWriter para
// tornar a saída legível em desenvolvimento.
func Init(level string) {
	lvl := zerolog.InfoLevel
	switch strings.ToLower(level) {
	case "debug":
		lvl = zerolog.DebugLevel
	case "info":
		lvl = zerolog.InfoLevel
	case "warn", "warning":
		lvl = zerolog.WarnLevel
	case "error":
		lvl = zerolog.ErrorLevel
	}

	// Console writer torna a saída amigável para humanos
	writer := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	Log = zerolog.New(writer).With().Timestamp().Logger().Level(lvl)

	// Também configurar o logger global do pacote zerolog/log
	zlog.Logger = Log
	zerolog.SetGlobalLevel(lvl)
}
