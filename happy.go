/*
Package happy implements high performance, minimalist Go web framework.
Example:
package main
import (
    "net/http"
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
)
  // Handler
func hello(c echo.Context) error {
    return c.String(http.StatusOK, "Hello, World!")
}
func main() {
	// Echo instance
	h := happy.Make()
	// Middleware
	h.Use(middleware.Logger())
	h.Use(middleware.Recover())
	// Routes
	h.GET("/", hello)
	// Start server
	h.Logger.Fatal(e.Start(":1323"))
}
Learn more at https://echo.labstack.com
*/

package happy

import "net/http"

func Run() {
	http.ListenAndServe(":8044", &RouterMux {})
}

