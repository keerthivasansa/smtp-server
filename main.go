package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	auth := createAuthMiddleware()
	smtpServer := NewSmtpServer()

	r := gin.Default()

	r.Use(auth.middleware)

	r.GET("/health", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "OK")
	})

	r.POST("/send", func(ctx *gin.Context) {
		var body SMTPRequestBody
		if err := ctx.ShouldBind(&body); err != nil {
			ctx.String(http.StatusBadRequest, "The body is not a valid SMTP request.")
		}
		err := smtpServer.SendMail(body)
		if err != nil {
			ctx.String(http.StatusBadRequest, "Could not send mail")
			return
		}
		ctx.String(http.StatusOK, body.Subject)
	})

	port := os.Getenv("API_PORT")
	addr := fmt.Sprintf(":%s", port)
	err := r.Run(addr)
	if err != nil {
		panic(err)
	}
}
