package handler

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

// AuthMiddleware : ヘッダーからユーザー名を取得し、Contextにセットするミドルウェア
func AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var userName string

			// 1. 環境変数 LOCAL が "true" の場合はローカル環境とみなし固定値をセット
			if os.Getenv("LOCAL") == "true" {
				userName = "ikura-hamu"
			} else {
				// 2. 本番環境 (NeoShowcase) : 付与されたヘッダーからユーザー名を取得
				userName = c.Request().Header.Get("X-Forwarded-User")
				if userName == "" {
					userName = c.Request().Header.Get("X-Showcase-User")
				}
			}

			// 3. ユーザー名が取得できなかった場合は 401 エラーを返す
			if userName == "" {
				return echo.NewHTTPError(http.StatusUnauthorized)
			}

			// 4. EchoのContextに "userName" というキーでセット
			c.Set("userName", userName)

			return next(c)
		}
	}
}
