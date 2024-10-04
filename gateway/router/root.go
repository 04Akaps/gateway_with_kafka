package router

import (
	"fmt"
	"github.com/04Akaps/gateway_with_kafka.git/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"strings"
	"time"
)

var allowHeaders = []string{
	"ORIGIN",
	"Content-Length",
	"Content-Type",
	"Access-Control-Allow-Headers",
	"Access-Control-Allow-Origin",
	"Authorization",
	"user-id-key",
}

type Router struct {
	engine *fiber.App
	port   string
}

func NewRouter(
	cfg config.GatewayCfg,
) {
	r := &Router{
		port: fmt.Sprintf(":%s", cfg.ServiceInfo.Port),
	}

	r.engine = fiber.New(fiber.Config{})
	r.engine.Use(r.traceMiddleware)

	r.engine.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     strings.Join([]string{"GET", "POST", "PUT", "DELETE", "PATCH"}, ", "),
		AllowHeaders:     strings.Join(allowHeaders, ", "),
		ExposeHeaders:    strings.Join(allowHeaders, ", "),
		AllowCredentials: false,
		AllowOriginsFunc: func(origin string) bool { return true },
		MaxAge:           12 * int(time.Hour.Seconds()),
	}))

}

func (r *Router) traceMiddleware(c *fiber.Ctx) error {
	// TODO -> kafka go routing을 활용한 send 및 header Key 검증
	// 1. header key 검증은 -> redis or localcache 사용 --> Azure같은 부분에서 부분적으로 메인터넌스를 걸어서 완전 관리형으로 동작하지 못하는 케이스가 있기 떄문에
	start := time.Now() // 요청 시작 시간 기록

	// 요청 처리
	err := c.Next()

	// 요청 정보 추출
	latency := time.Since(start)
	method := c.Method()                    // HTTP 메소드 (GET, POST 등)
	url := c.BaseURL()                      // 요청 URL
	headerValue := c.Get("TODO-Key")        // 헤더의 특정 키 값 추출 ("Your-Key" 부분은 추적할 실제 키 값으로 변경)
	statusCode := c.Response().StatusCode() // HTTP 상태 코드

	fmt.Println(latency, method, url, headerValue, statusCode)

	return err
}
