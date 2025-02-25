package repository

import (
	"fmt"
	"record-shop-rest-api/model"

	"gorm.io/gorm"
)

type IRecordRepository interface {
	CreateRecord(record *model.Record) error
	GetRecordList() ([]model.Record, error)
	GetDetail(title string) (model.DetailResponse, error)
	GetRecordByTitle(title string) ([]model.Record, error)
	GetRecordByArtist(artist string) ([]model.Record, error)
	UpdateRecord(task *model.Record) error
	DeleteRecord(id uint) error
}

type recordRepository struct {
	db *gorm.DB
}

func NewRecordRepository(db *gorm.DB) IRecordRepository {
	return &recordRepository{db}
}

func (rr *recordRepository) CreateRecord(record *model.Record) error {
	// Gormはこのように引数を直接変更する、.Errorにチェインさせるため
	// 引数変えていいのか、違和感があるが変えたくなければ、recordRrepositoryの例見て
	// 引数受け取らず、戻り値で([]model.Record, error) を返してる
	if err := rr.db.Create(record).Error; err != nil {
		return err
	}
	return nil
}

func (rr *recordRepository) GetRecordList() ([]model.Record, error) {
	var records []model.Record
	if err := rr.db.
		Order("release_year ASC, artist ASC, title ASC").
		Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

func (rr *recordRepository) GetDetail(title string) (model.DetailResponse, error) {
	var records []struct {
		RecordTitle   string
		AlbumImageUrl string
		TrackNumber   uint
		TrackTitle    string
	}
	result := rr.db.
		Table("records").
		Select(
			"records.title AS record_title, details.album_image_url, tracks.track_number, tracks.track_title").
		Joins("JOIN details ON details.record_id = records.id").
		Joins("JOIN tracks ON tracks.detail_id = details.id").
		Where("records.title = ?", title).
		Scan(&records)

	if result.Error != nil {
		fmt.Printf("Error occurred while querying: %v\n", result.Error)
		return model.DetailResponse{}, result.Error
	}

	// クエリが成功したが結果が空の場合
	if result.RowsAffected == 0 {
		fmt.Println("No records found for the given title.")
		return model.DetailResponse{}, result.Error
	}

	response := model.DetailResponse{
		RecordTitle:   records[0].RecordTitle,
		AlbumImageUrl: records[0].AlbumImageUrl,
	}
	for _, r := range records {
		track := model.TrackInfo{
			TrackNumber: r.TrackNumber,
			TrackTitle:  r.TrackTitle,
		}
		response.Tracks = append(response.Tracks, track)
	}
	return response, nil
}

func (rr *recordRepository) GetRecordByTitle(title string) ([]model.Record, error) {
	var records []model.Record
	if err := rr.db.
		Order("release_year ASC, artist ASC, title ASC").
		Where("title=?", title).Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

func (rr *recordRepository) GetRecordByArtist(artist string) ([]model.Record, error) {
	var records []model.Record
	if err := rr.db.
		Order("release_year ASC, artist ASC, title ASC").
		Where("artist=?", artist).Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

func (rr *recordRepository) UpdateRecord(record *model.Record) error {
	// Save: レコードが存在すればその全てのフィールドを更新、存在しなければ、新規作成
	// つまりupsert
	// idを条件にSave、time.Time型のupdate_atを自動で更新してくれる
	// .Clauses(clause.Returning{})は使えない、Updatesは使える
	// result := tr.db.Save(record)
	// SaveはCreatedAtまで更新してしまう（今回何も渡してないので0001-01-01 00:00:00+00）
	// にしてしまう、なのでOmitで更新対象外とする
	result := rr.db.Model(record).Omit("CreatedAt").Save(record)
	// Updates: 複数のフィールドを更新する時に使う、1項目の場合Updateを使う
	// time.Time型のupdate_atを自動で更新してくれる
	// resultU := tr.db.Model(&model.Record{}).Where("id = ?", id).Updates(map[string]interface{}{
	// 	"Title":       record.Title,
	// 	"Artist":      record.Artist,
	// 	"Genre":       record.Genre,
	// 	"ReleaseYear": record.ReleaseYear,
	// })
	if result.Error != nil {
		return result.Error
	}
	// 更新対象が存在しない場合、エラーにならないのでこのチェックがいる
	if result.RowsAffected < 1 {
		return fmt.Errorf("object does not exist")
	}
	return nil
}

func (rr *recordRepository) DeleteRecord(id uint) error {
	result := rr.db.Where("id=?", id).Delete(&model.Record{})
	if result.Error != nil {
		return result.Error
	}
	// 削除対象が存在しない場合、エラーにならないのでこのチェックがいる
	if result.RowsAffected < 1 {
		return fmt.Errorf("object does not exist")
	}
	return nil
}
