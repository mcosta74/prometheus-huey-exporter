/*
The prometheus-nats-exporter exposes metrics about Huey2 tasks executions
*/
package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"syscall"

	"log/slog"

	"github.com/oklog/run"
	"github.com/redis/go-redis/v9"

	"github.com/mcosta74/prometheus-huey-exporter/exporter"
)

var (
	version = "<< dev >>"
	commit  = ""
	date    = ""
)

var (
	fs               *flag.FlagSet
	logLevel         slog.Level
	logFormat        string
	redisAddr        string
	redisChannel     string
	metricsNamespace string
	webPath          string
	webListenAddr    string
	showVersion      bool
)

func init() {
	fs = flag.NewFlagSet("prometheus-huey-exporter", flag.ExitOnError)
	fs.TextVar(&logLevel, "log.level", getEnvLogLevel("HUEY_EXPORTER_LOG_LEVEL", slog.LevelInfo), "Log level (debug, info, warn, error)")
	fs.StringVar(&logFormat, "log.format", getEnvString("HUEY_EXPORTER_LOG_FORMAT", "text"), "Log format (text, json)")
	fs.StringVar(&redisAddr, "redis.address", getEnvString("HUEY_EXPORTER_REDIS_ADDR", "localhost:6379"), "Address of the Redis instance to connect to")
	fs.StringVar(&redisChannel, "redis.channel", getEnvString("HUEY_EXPORTER_REDIS_CHANNEL", "events"), "Redis channel to subscribe to listen for events")
	fs.StringVar(&metricsNamespace, "metrics.namespace", getEnvString("HUEY_EXPORTER_METRICS_NAMESPACE", ""), "Namespace for metrics")
	fs.StringVar(&webPath, "web.telemetry-path", getEnvString("HUEY_EXPORTER_WEB_PATH", "/metrics"), "Path under which to expose metrics")
	fs.StringVar(&webListenAddr, "web.listen-address", getEnvString("HUEY_EXPORTER_WEB_LISTEN_ADDRESS", ":9234"), "Address to listen on for web interface and telemetry")
	fs.BoolVar(&showVersion, "version", false, "Show version information")

}

func main() {
	fs.Parse(os.Args[1:])

	if showVersion {
		fmt.Println(version)
		os.Exit(1)
	}

	logger := getLogger()

	logger.Info("service started")
	defer logger.Info("service stopped")

	// Establish connection with Redis
	rc, err := connectToRedis(redisAddr)
	if err != nil {
		logger.Error("connection with Redis failed", "addr", redisAddr, "err", err)
		os.Exit(1)
	}
	defer rc.Close()

	logger.Info("connected with Redis", "addr", redisAddr)

	var (
		metrics       = exporter.SetupMetrics(metricsNamespace)
		eventHandler  = exporter.NewEventHandler(metrics, logger)
		eventListener = exporter.NewEventListener(rc, redisChannel, logger.With("component", "event-listener"), eventHandler)
		httpHandler   = exporter.MakeHTTPHandler(webPath)
	)

	// Setup go-routines
	var g run.Group
	{
		// Signal Handler
		g.Add(run.SignalHandler(context.Background(), syscall.SIGTERM, syscall.SIGINT))
	}

	{
		// Prometheus HTTP
		webListener, err := net.Listen("tcp", webListenAddr)
		if err != nil {
			logger.Error("error listening on TCP", "addr", webListenAddr, "err", err)
			os.Exit(1)
		}

		g.Add(func() error {
			logger.Info("starting HTTP metrics server", "addr", fmt.Sprintf("http://%s%s", webListener.Addr(), webPath))
			return http.Serve(webListener, httpHandler)
		}, func(err error) {
			webListener.Close()
		})
	}

	{
		// Redis event listener
		ctx, cancel := context.WithCancel(context.Background())
		g.Add(func() error {
			return eventListener.Run(ctx)
		}, func(err error) {
			cancel()
		})
	}
	logger.Info("service shutdown", "err", g.Run())
}

func getLogger() *slog.Logger {
	hOpts := slog.HandlerOptions{
		AddSource: true,
		Level:     logLevel,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.SourceKey:
				// use only filename for source
				source := a.Value.Any().(*slog.Source)
				source.File = filepath.Base(source.File)

			case slog.TimeKey:
				// Use UTC time for timestamp
				a.Value = slog.TimeValue(a.Value.Time().UTC())
			}
			return a
		},
	}

	var handler slog.Handler = slog.NewTextHandler(os.Stderr, &hOpts)
	if logFormat == "json" {
		handler = slog.NewJSONHandler(os.Stderr, &hOpts)
	}
	return slog.New(handler)
}

func connectToRedis(addr string) (*redis.Client, error) {
	rc := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	if _, err := rc.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}
	return rc, nil
}

func getEnvString(key, defaultVal string) string {
	if val, found := os.LookupEnv(key); found {
		return val
	}
	return defaultVal
}

func getEnvLogLevel(key string, defaultVal slog.Level) slog.Level {
	lvl := defaultVal
	if val, found := os.LookupEnv(key); found {
		if err := lvl.UnmarshalText([]byte(val)); err != nil {
			return defaultVal
		}
	}
	return lvl
}
