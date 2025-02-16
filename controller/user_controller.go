package controller

import (
	"net/http"
	"os"
	"record-shop-rest-api/model"
	"record-shop-rest-api/usecase"
	"time"

	"github.com/labstack/echo/v4"
)

type IUserControler interface {
	SignUp(c echo.Context) error
	LogIn(c echo.Context) error
	LogOut(c echo.Context) error
	CsrfToken(c echo.Context) error
}

type userController struct {
	uu usecase.IUserUsecase
}

func NewUserController(uu usecase.IUserUsecase) IUserControler {
	return &userController{uu}
}

func (uc *userController) CsrfToken(c echo.Context) error {
	// csrf Tokenはecho.Contextから取得出来る、string型でアサーション
	token := c.Get("csrf").(string)
	// JSONでクライアントにcsrf tokenをレスポンスで返す
	return c.JSON(http.StatusOK, echo.Map{"csrf_token": token})
}

func (uc *userController) SignUp(c echo.Context) error {
	user := model.User{}
	// clientから送られてくるリクエストBodyをUserオブジェクトのポインタが指し示す先の値に格納する
	// つまり構造体にBind
	if err := c.Bind(&user); err != nil {
		// 変換に失敗した場合
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	userRes, err := uc.uu.SignUp(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, userRes)
}

func (uc *userController) LogIn(c echo.Context) error {
	user := model.User{}
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	// JWTを生成するので
	tokenString, err := uc.uu.Login(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	// JWTをサーバーサイドでCookieに設定
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = tokenString
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.Path = "/"
	cookie.Domain = os.Getenv("API_DOMAIN")
	cookie.Secure = true                    // trueにしておく必要があるが、いったんPOSTMANで確認したいのでコメントアウト
	cookie.HttpOnly = true                  // クライアントのJSからTokenの値が読み取れないように
	cookie.SameSite = http.SameSiteNoneMode // frontendとbackendのdomainが違うクロスドメイン間のCookie送受信になるので、クロスサイト・スクリプティング攻撃やセッションハイジャックなどのリスクを軽減
	c.SetCookie(cookie)                     // 上で設定したCookieをHTTPレスポンスに含める
	return c.NoContent(http.StatusOK)
}

func (uc *userController) LogOut(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = ""           // 値をクリア
	cookie.Expires = time.Now() // 有効期限がすぐ切れるように
	cookie.Path = "/"
	cookie.Domain = os.Getenv("API_DOMAIN")
	cookie.Secure = true // trueにしておく必要があるが、いったんPOSTMANで確認したいのでコメントアウト
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteNoneMode
	c.SetCookie(cookie)
	return c.NoContent(http.StatusOK)
}
