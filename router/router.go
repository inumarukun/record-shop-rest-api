package router

import (
	"net/http"
	"os"
	"record-shop-rest-api/controller"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter(uc controller.IUserControler, rc controller.IRecordController) *echo.Echo {
	e := echo.New()
	// CORS middleware
	// クロスオリジンリソース共有 (CORS) は、悪意のあるウェブサイトが明示的な権限を持たずに
	// 他のサイト(クロスドメイン)のデータ (Box APIなど) にアクセスするのを防ぐために、
	// ウェブブラウザで利用されているセキュリティメカニズム
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		// 許可するオリジン一覧、REACT用ローカルホストと、Vercelで取得したドメイン
		// AllowOrigins: []string{"http://localhost:3000", os.Getenv("FE_URL")},
		AllowOrigins: []string{"http://localhost:5173", os.Getenv("FE_URL")},
		// 許可するヘッダ一覧、echo.HeaderXCSRFTokenでヘッダ経由でCERF tokenを受取れるように
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept,
			echo.HeaderAccessControlAllowHeaders, echo.HeaderXCSRFToken},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		// Cookieの送受信を可能に
		AllowCredentials: true,
	}))

	// CSRF middleware
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		// CSRFトークンを含むCookieが '/'パス以下全てのリクエストに対して送信されることを指定
		CookiePath: "/",
		// CSRFトークンを含むCookieを指定したドメインに対してのみ送信
		CookieDomain:   os.Getenv("API_DOMAIN"),
		CookieHTTPOnly: true,
		CookieSameSite: http.SameSiteNoneMode,
		// CookieSameSite: http.SameSiteDefaultMode, // POSTMAN動作確認用(SecudeMode: false)
		// CookieMaxAge: 60,　// csrf tokenの有効期限、デフォルト24H、秒単位
	}))

	e.POST("/signup", uc.SignUp)
	e.POST("/login", uc.LogIn)
	e.POST("/logout", uc.LogOut)
	e.GET("/csrf", uc.CsrfToken)

	r := e.Group("/records")
	// 実質これでGET: /records
	r.GET("", rc.ViewList)

	// /records以下の全てのルートに対して、JWT認証を適用
	// リクエストにcookie: token が含まれている場合、
	// JWTトークンが検証され、認証情報がリクエストに追加される
	// つまりloginしていないと/records以下にはアクセス出来ない
	// これは先頭にlogin画面を配備し、loginしていないと以降の処理を許可しない場合に有効
	r.Use(echojwt.WithConfig(echojwt.Config{
		// Jwtを生成した時と同じ秘密鍵
		SigningKey: []byte(os.Getenv("SECRET")),
		// クライアントから送られてくるJWTがどこに格納されているか
		// 今回はCookieにtokenという形で実装している
		TokenLookup: "cookie:token",
	}))

	// 実質これでPOST: /records
	r.POST("", rc.CreateRecord)
	// /records/idのidいらないだろうと思ったけどupdateのwhere条件にいる)
	// TODO：：でも結局idはc.BindでBodyから取得しているから不要
	r.PUT("/:id", rc.UpdateRecord)
	r.DELETE("/:id", rc.DeleteRecord)
	return e
}
