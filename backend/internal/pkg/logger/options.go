package logger

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	// DefaultContainerLogPath 为容器内默认日志文件路径。
	DefaultContainerLogPath = "/app/data/logs/sub2api.log"
	defaultLogFilename      = "sub2api.log"
)

type InitOptions struct {
	Level           string
	Format          string
	ServiceName     string
	Environment     string
	Caller          bool
	StacktraceLevel string
	Output          OutputOptions
	Rotation        RotationOptions
	Sampling        SamplingOptions
}

type OutputOptions struct {
	ToStdout bool
	ToFile   bool
	FilePath string
}

type RotationOptions struct {
	MaxSizeMB  int
	MaxBackups int
	MaxAgeDays int
	Compress   bool
	LocalTime  bool
}

type SamplingOptions struct {
	Enabled    bool
	Initial    int
	Thereafter int
}

func (o InitOptions) normalized() InitOptions {
	out := o
	out.Level = strings.ToLower(strings.TrimSpace(out.Level))
	if out.Level == "" {
		out.Level = "info"
	}
	out.Format = strings.ToLower(strings.TrimSpace(out.Format))
	if out.Format == "" {
		out.Format = "console"
	}
	out.ServiceName = strings.TrimSpace(out.ServiceName)
	if out.ServiceName == "" {
		out.ServiceName = "sub2api"
	}
	out.Environment = strings.TrimSpace(out.Environment)
	if out.Environment == "" {
		out.Environment = "production"
	}
	out.StacktraceLevel = strings.ToLower(strings.TrimSpace(out.StacktraceLevel))
	if out.StacktraceLevel == "" {
		out.StacktraceLevel = "error"
	}
	if !out.Output.ToStdout && !out.Output.ToFile {
		out.Output.ToStdout = true
	}
	out.Output.FilePath = resolveLogFilePath(out.Output.FilePath)
	if out.Rotation.MaxSizeMB <= 0 {
		out.Rotation.MaxSizeMB = 100
	}
	if out.Rotation.MaxBackups < 0 {
		out.Rotation.MaxBackups = 10
	}
	if out.Rotation.MaxAgeDays < 0 {
		out.Rotation.MaxAgeDays = 7
	}
	if out.Sampling.Enabled {
		if out.Sampling.Initial <= 0 {
			out.Sampling.Initial = 100
		}
		if out.Sampling.Thereafter <= 0 {
			out.Sampling.Thereafter = 100
		}
	}
	return out
}

func resolveLogFilePath(explicit string) string {
	explicit = strings.TrimSpace(explicit)
	if explicit != "" {
		return explicit
	}
	dataDir := strings.TrimSpace(os.Getenv("DATA_DIR"))
	if dataDir != "" {
		return filepath.Join(dataDir, "logs", defaultLogFilename)
	}
	return DefaultContainerLogPath
}

func bootstrapOptions() InitOptions {
	quiet := envBool("LOG_QUIET", false)
	toStdout := envBool("LOG_OUTPUT_TO_STDOUT", true)
	toFile := envBool("LOG_OUTPUT_TO_FILE", false)
	filePath := strings.TrimSpace(os.Getenv("LOG_OUTPUT_FILE_PATH"))
	if quiet {
		toStdout = false
		toFile = true
		if filePath == "" {
			filePath = os.DevNull
		}
	}
	return InitOptions{
		Level:           envString("LOG_LEVEL", "info"),
		Format:          envString("LOG_FORMAT", "console"),
		ServiceName:     envString("LOG_SERVICE_NAME", "sub2api"),
		Environment:     envString("LOG_ENV", "bootstrap"),
		StacktraceLevel: envString("LOG_STACKTRACE_LEVEL", "error"),
		Output: OutputOptions{
			ToStdout: toStdout,
			ToFile:   toFile,
			FilePath: filePath,
		},
		Rotation: RotationOptions{
			MaxSizeMB:  100,
			MaxBackups: 10,
			MaxAgeDays: 7,
			Compress:   true,
			LocalTime:  true,
		},
		Sampling: SamplingOptions{
			Enabled:    false,
			Initial:    100,
			Thereafter: 100,
		},
	}
}

func envString(key, defaultValue string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return defaultValue
}

func envBool(key string, defaultValue bool) bool {
	value := strings.TrimSpace(strings.ToLower(os.Getenv(key)))
	switch value {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		return defaultValue
	}
}

func parseLevel(level string) (Level, bool) {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		return LevelDebug, true
	case "info":
		return LevelInfo, true
	case "warn":
		return LevelWarn, true
	case "error":
		return LevelError, true
	default:
		return LevelInfo, false
	}
}

func parseStacktraceLevel(level string) (Level, bool) {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "none":
		return LevelFatal + 1, true
	case "error":
		return LevelError, true
	case "fatal":
		return LevelFatal, true
	default:
		return LevelError, false
	}
}

func samplingTick() time.Duration {
	return time.Second
}
