package apiserver

import (
	"fmt"
	"log"
	"net/http"

	"p2p/internal/config"
	"p2p/internal/db"
	"p2p/internal/errors"
	"p2p/internal/format"
	"p2p/pkg/peer"
	peerHTTP "p2p/pkg/peer/delivery/http"
	peerService "p2p/pkg/peer/service"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	_peerSvc peer.Servicer
	mgr      *peer.Manager
)

func initialize(cfg *config.Configuration) (err error) {
	db, err := db.InitDatabase(cfg, "heroku_cleardb")
	if err != nil {
		return err
	}
	log.Println("Init Db complete...")
	err = db.DB().Ping()
	if err != nil {
		log.Printf("main: mysql ping error: %v \n", err)
		return err
	}

	mgr = peer.InitManager()
	_peerSvc = peerService.NewService(mgr)

	return nil
}
func newEchoHandler(cfg *config.Configuration) http.Handler {
	e := echo.New()
	// setting

	e.Debug = false
	e.HTTPErrorHandler = errors.HTTPErrorHandlerForEcho
	// 所有 API 皆經過 CORS middeware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodPatch,
		},
		AllowHeaders: []string{
			"*",
			echo.HeaderAuthorization,
			echo.HeaderContentType,
			echo.HeaderOrigin,
			echo.HeaderContentLength,
		},
	}))
	// cover all api error response
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err != nil {
				logFields := Fields{}

				// get request data
				req := c.Request()
				{
					logFields["requestMethod"] = req.Method
					logFields["requestURL"] = req.URL.String()
				}

				str := fmt.Sprintf("%+v, error message : %+v\n", logFields, err)
				msg := format.GetCMDColor(format.Color_red, "[API ERROR] ")
				log.Printf(msg + str)
			}
			return err
		}
	})

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World!!!")
	})
	e.File("/sfu", "view/sfu.html")
	e.File("/p2p", "view/p2p.html")

	_peerHTTPHandler := peerHTTP.NewHandler(_peerSvc, mgr)
	peerHTTP.SetRoutes(e, _peerHTTPHandler)

	return e

}
