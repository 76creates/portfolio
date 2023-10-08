package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/76creates/portfolio/dsn.onl/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/google/uuid"
	"reflect"
	"runtime"
	"time"
)

type Conf struct {
	IP           string        `mapstructure:"ip" yaml:"ip"`
	Port         string        `mapstructure:"port" yaml:"port"`
	Concurrency  int           `mapstructure:"concurency" yaml:"concurency"`
	BodyLimit    int           `mapstructure:"bodyLimit" yaml:"bodyLimit"`
	ReadTimeout  time.Duration `mapstructure:"readTimeout" yaml:"readTimeout"`
	WriteTimeout time.Duration `mapstructure:"writeTimeout" yaml:"writeTimeout"`
	IdleTimeout  time.Duration `mapstructure:"idleTimeout" yaml:"idleTimeout"`
	Logger       LoggerConf    `mapstructure:"logger" yaml:"logger"`
}

type LoggerConf struct {
	File              string            `mapstructure:"file" yaml:"file"`
	Level             string            `mapstructure:"level" yaml:"level"`
	LogCaller         bool              `mapstructure:"logCaller" yaml:"logCaller"`
	LogToStdout       bool              `mapstructure:"logToStdout" yaml:"logToStdout"`
	StaticExtraFields map[string]string `mapstructure:"staticExtraFields" yaml:"staticExtraFields"`
}

type server struct {
	ctx    context.Context
	logger *logger.Logger
}

func (c Conf) GetApp(ctx context.Context) *fiber.App {
	srv := server{
		ctx: ctx,
	}
	lgr := logger.NewLogger()
	if c.Logger.File != "" {
		lgr = lgr.WithFile(c.Logger.File)
	}
	lgr.WithLevel(c.Logger.Level)
	if c.Logger.LogCaller {
		lgr = lgr.WithCaller()
	}
	if c.Logger.LogToStdout {
		lgr = lgr.WithSilent(true)
	}
	for k, v := range c.Logger.StaticExtraFields {
		lgr = lgr.WithStr(k, v)
	}
	srv.logger = lgr

	serverSettings := fiber.Config{
		ServerHeader:          "go/Fiber",
		BodyLimit:             c.BodyLimit,
		Concurrency:           c.Concurrency,
		DisableKeepalive:      true,
		DisableStartupMessage: true,
		ReadTimeout:           c.ReadTimeout,
		WriteTimeout:          c.WriteTimeout,
		IdleTimeout:           c.IdleTimeout,
	}
	app := fiber.New(serverSettings)

	v1 := app.Group("v1")
	v1.Use(
		recover.New(
			recover.Config{
				EnableStackTrace:  true,
				StackTraceHandler: StackTraceHandler(lgr),
			},
		),
	)
	addRootRoutes(app)
	return app
}

func (s *server) NewRequest(c *fiber.Ctx) *logger.RequestLogger {
	reqId := uuid.New()
	c.Set("Req-Id", reqId.String())
	c.Locals("id", reqId)
	return s.logger.NewRequest(reqId)
}

func StackTraceHandler(panicLogger *logger.Logger) func(c *fiber.Ctx, e interface{}) {
	return func(c *fiber.Ctx, e interface{}) {
		_id := c.Locals("id")
		var l logger.Lgr
		if _id != nil {
			if reflect.TypeOf(_id).ConvertibleTo(reflect.TypeOf(uuid.Nil)) {
				l = panicLogger.WithStr("id", _id.(uuid.UUID).String())
			} else {
				l = panicLogger
			}
		}

		buf := make([]byte, 4096)
		buf = buf[:runtime.Stack(buf, false)]
		msg := fmt.Sprintf("panic: %v\n%s\n", e, buf)
		l.Error(errors.New(msg))
	}
}
