package http

import "github.com/labstack/echo"

func SetRoutes(e *echo.Echo, h *handler) {
	e.GET("/peer", h.webSocketEndpoint)
}
