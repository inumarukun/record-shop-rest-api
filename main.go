package main

import (
	"record-shop-rest-api/controller"
	"record-shop-rest-api/db"
	"record-shop-rest-api/repository"
	"record-shop-rest-api/router"
	"record-shop-rest-api/usecase"
	"record-shop-rest-api/validator"
)

// $env:GO_ENV="dev"; go run main.go
// debug: 左サイドバー Run and Debug
// Notice: docker desktop起動、record-shop-rest-api(postgres)を起動させておくこと
func main() {
	db := db.NewDB()
	userValidator := validator.NewUserValidator()
	recordValidator := validator.NewRecordValidator()
	userRepository := repository.NewUserRepository(db)
	recordRepository := repository.NewRecordRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepository, userValidator)
	recordUsecase := usecase.NewRecordUsecase(recordRepository, recordValidator)
	userController := controller.NewUserController(userUsecase)
	recordController := controller.NewRecordController(recordUsecase)

	e := router.NewRouter(userController, recordController)
	// server起動
	// error発生時、log出力して終了
	e.Logger.Fatal(e.Start(":8080"))
}
