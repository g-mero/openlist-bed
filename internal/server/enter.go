package server

import (
	"errors"
	"log"
	"openlist-bed/internal/server/middleware"
	"openlist-bed/pkg/R"
	"openlist-bed/pkg/config"
	"openlist-bed/pkg/logger"
	"openlist-bed/pkg/utils"
	"time"

	"os"
	"os/signal"
	"syscall"

	"github.com/bytedance/sonic"
	"github.com/cshum/vipsgen/vips"
	"github.com/gofiber/fiber/v3"
	"github.com/spf13/viper"
)

func Run() {
	// init config
	err := config.InitConfig()
	if err != nil {
		log.Fatal("can not init config: ", err.Error())
	}

	err = logger.Init(logger.Config{
		Level:    viper.GetString("log_level"),
		FilePath: "data/server.log",
	})
	if err != nil {
		log.Fatal("can not init logger: ", err.Error())
	}
	defer logger.Logger.Sync()

	vips.Startup(nil)

	app := fiber.New(fiber.Config{
		AppName: "openlist-bed-server",
		// Override default error handler
		ErrorHandler: func(c fiber.Ctx, err error) error {
			// Status code defaults to 500
			statusCode := fiber.StatusInternalServerError
			code := "INTERNAL_SERVER_ERROR"

			// Retrieve the custom status code if it's a *fiber.Error
			var e *fiber.Error
			if errors.As(err, &e) {
				statusCode = e.Code
			}

			var re *R.RespError
			if errors.As(err, &re) {
				statusCode = re.StatusCode
				code = re.Code
			}

			c.Status(statusCode)

			return c.JSON(utils.H{
				"code":      code,
				"msg":       err.Error(),
				"timestamp": time.Now().Unix(),
			})
		},
		JSONDecoder: sonic.Unmarshal,
		JSONEncoder: sonic.Marshal,
	})

	app.Use(middleware.Cors())
	registerRoutes(app)

	go func() {
		utils.PrintBox("openlist-bed-server is running", "1. server started on port 6001")
		if err = app.Listen(":6001", fiber.ListenConfig{
			DisableStartupMessage: true,
		}); err != nil {
			log.Fatal("Server Listen Error: ", err)
		}
	}()

	// graceful shut down
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	vips.Shutdown()
	log.Println("Shutdown Server ...")
	_ = app.Shutdown()
}
