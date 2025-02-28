package model

import "time"

// Gormはデフォルトでモデル名を複数形でテーブル名として使う
type Record struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Title       string    `json:"title" gorm:"not null; default: ''"`
	Artist      string    `json:"artist" gorm:"not null; default: ''"`
	Genre       string    `json:"genre" gorm:"not null; default: ''"`
	Style       string    `json:"style" gorm:"not null; default: ''"`
	ReleaseYear int       `json:"release_year" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"not null"`
	// time.Time 型の場合、default:nullは扱えない
	// null を許容したい場合は、*time.Time 型を使う
	UpdatedAt *time.Time `json:"updated_at" gorm:"default:null"`
}

type RecordResponse struct {
	// IDはupdateで使うので返す
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Artist      string `json:"artist"`
	Genre       string `json:"genre"`
	Style       string `json:"style"`
	ReleaseYear int    `json:"release_year"`
	// omitempty: フィールドがゼロ値の場合、JSONエンコード時にそのフィールドは省略
	// jsonタグのオプションは,区切りの間に空白入れると警告(警告だが入れないほうが無難)
	Error *ErrorResponse `json:"error,omitempty"` // エラーが無い場合はnil
}

// Recordフィールドを通じて、外部キーがどのモデルのどのフィールドに関連するかを明示
// この構造体フィールドはGORMが自動的に利用するため、JSONタグで返さない設定にしている（json:"-"）
// foreignKey:RecordId: DetailのRecordIdフィールドをRecordテーブルの外部キーとして使用します。
// references:ID: RecordのIDフィールドが外部キーの参照先
// constraint:
//   OnUpdate:CASCADE、RecordのIDが変更された場合、それに応じてDetailも更新
//   OnDelete:SET NULL、Recordが削除された場合、紐づくDetailのRecordIdフィールドをNULLに設定
type Detail struct {
	ID             uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	RecordId       uint   `json:"recordId" gorm:"not null"`
	AlbumImageUrl  string `json:"albumImageUrl" gorm:"default: ''"`
	YoutubeTitle   string `json:"youtubeTitle" gorm:"default: ''"`
	YoutubeVideoId string `json:"youtubeVideoId" gorm:"default: ''"`
	Record         Record `json:"-" gorm:"foreignKey:RecordId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

// TrackはTrackInfoをリストで持たないと→そんなことない ↓は出力形式なだけ
type Track struct {
	ID          uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	DetailId    uint   `json:"detailId" gorm:"not null"`
	TrackNumber uint   `json:"trackNumber"  gorm:"not null"`
	TrackTitle  string `json:"trackTitle" gorm:"not null; default: ''"`
	Detail      Record `json:"-" gorm:"foreignKey:DetailId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type TrackInfo struct {
	TrackNumber uint   `json:"trackNumber"`
	TrackTitle  string `json:"trackTitle"`
}
type DetailResponse struct {
	RecordTitle    string      `json:"recordTitle"`
	AlbumImageUrl  string      `json:"albumImageUrl"`
	YoutubeTitle   string      `json:"youtubeTitle"`
	YoutubeVideoId string      `json:"youtubeVideoId"`
	Tracks         []TrackInfo `json:"tracks"`
}
