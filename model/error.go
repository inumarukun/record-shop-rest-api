package model

type ErrorResponse struct {
	Code    string `json:"code"`    // エラーコード (例: "ValidationError", "InternalError")
	Message string `json:"message"` // ユーザ向けのエラーメッセージ
	Details string `json:"details"` // より詳細な内部情報やデバッグ用メッセージ
}
