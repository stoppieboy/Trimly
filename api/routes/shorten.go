package routes

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type Request struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"short"`
	Expiry      time.Duration `json:"expiry"`
}

type Response struct {
	URL                 string        `json:"url"`
	CustomShort         string        `json:"short"`
	Expiry              time.Duration `json:"expiry"`
	XRateLimitRemaining int           `json:"rate_limit"`
	XRateLimitReset     time.Duration `json:"rate_limit_reset"`
}

func ShortenURL(c *fiber.Ctx) Response {
	return Response{}
}
