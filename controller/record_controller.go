package controller

import (
	"net/http"
	"record-shop-rest-api/model"
	"record-shop-rest-api/usecase"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type IRecordController interface {
	CreateRecord(c echo.Context) error
	ViewList(c echo.Context) error
	GetDetail(c echo.Context) error
	GetRecordByTitle(c echo.Context) error
	GetRecordByArtist(c echo.Context) error
	UpdateRecord(c echo.Context) error
	DeleteRecord(c echo.Context) error
}

type recordController struct {
	ru usecase.IRecordUsecase
}

func NewRecordController(ru usecase.IRecordUsecase) IRecordController {
	return &recordController{ru}
}

func (rc *recordController) CreateRecord(c echo.Context) error {
	record := model.Record{}
	// clientから送られてくるリクエストBodyをRecordオブジェクトのポインタが指し示す先の値に格納する
	// つまり構造体にBind
	if err := c.Bind(&record); err != nil {
		// 変換に失敗した場合
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	recordRes, err := rc.ru.CreateRecord(record)
	if err != nil {
		if recordRes.Error != nil {
			// ここでErrorを返しているから、フロント側でerr.response.data.messageで受けれる
			return c.JSON(http.StatusBadRequest, recordRes.Error)
		}
		// map[string]string: のstringは、key, value
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusCreated, recordRes)
}

func (rc *recordController) ViewList(c echo.Context) error {
	// このリソースは公開リソースとするため以下の処理は不要
	// echojwtのmiddleware内部で"user"キーを自動付与
	// user := c.Get("user").(*jwt.Token)
	// if user == nil {
	// 	return c.JSON(http.StatusUnauthorized, "Unauthorized")
	// }

	recordResponse, err := rc.ru.GetRecordList()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, recordResponse)
}

func (rc *recordController) GetDetail(c echo.Context) error {
	title := c.Param("title")
	recordReponse, err := rc.ru.GetDetail(title)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, recordReponse)
}

func (rc *recordController) GetRecordByTitle(c echo.Context) error {
	// クライアントから受取るリクエストボディを構造体に変換
	record := model.Record{}
	if err := c.Bind(&record); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	userResponse, err := rc.ru.GetRecordByTitle(record.Title)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, userResponse)
}

func (rc *recordController) GetRecordByArtist(c echo.Context) error {
	// クライアントから受取るリクエストボディを構造体に変換
	record := model.Record{}
	if err := c.Bind(&record); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	userResponse, err := rc.ru.GetRecordByArtist(record.Artist)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, userResponse)
}

func (rc *recordController) UpdateRecord(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	if user == nil {
		return c.JSON(http.StatusUnauthorized, "Unauthorized")
	}

	record := model.Record{}
	if err := c.Bind(&record); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	recordRes, err := rc.ru.UpdateRecord(record)
	if err != nil {
		if recordRes.Error != nil {
			// ここでErrorを返しているから、フロント側でerr.response.data.messageで受けれる
			return c.JSON(http.StatusBadRequest, recordRes.Error)
		}
		// map[string]string: のstringは、key, value
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	return c.JSON(http.StatusOK, recordRes)
}

func (rc *recordController) DeleteRecord(c echo.Context) error {
	// echojwtのmiddleware内部で"user"キーを自動付与
	user := c.Get("user").(*jwt.Token)
	if user == nil {
		return c.JSON(http.StatusUnauthorized, "Unauthorized")
	}
	id, _ := strconv.Atoi(c.Param("id"))

	err := rc.ru.DeleteRecord(uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
