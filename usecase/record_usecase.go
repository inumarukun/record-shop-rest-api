package usecase

import (
	"record-shop-rest-api/common"
	"record-shop-rest-api/model"
	"record-shop-rest-api/repository"
)

type IRecordUsecase interface {
	CreateRecord(record model.Record) (model.RecordResponse, error)
	GetRecordList() ([]model.RecordResponse, error)
	GetDetail(title string) (model.DetailResponse, error)
	GetRecordByTitle(title string) ([]model.RecordResponse, error)
	GetRecordByArtist(artist string) ([]model.RecordResponse, error)
	UpdateRecord(task model.Record) (model.RecordResponse, error)
	DeleteRecord(id uint) error
}

type recordUsecase struct {
	rr repository.IRecordRepository
}

// constructor injection
func NewRecordUsecase(rr repository.IRecordRepository) IRecordUsecase {
	// &recordUsecase構造体がIrecordUsecaseを満たすため、interfaceの定義を全て実装する必要がある
	return &recordUsecase{rr}
}

func (ru *recordUsecase) CreateRecord(record model.Record) (model.RecordResponse, error) {
	// これで新しいstructが出来る、idはGormが自動で入れる？
	newRecord := model.Record{
		Artist:      record.Artist,
		Title:       record.Title,
		Genre:       record.Genre,
		Style:       record.Style,
		ReleaseYear: record.ReleaseYear,
	}
	if err := ru.rr.CreateRecord(&newRecord); err != nil {
		return model.RecordResponse{}, err
	}
	// CreateUserが成功すれば、newUser、つまり引数が新しいユーザになっている、それを詰めて返す
	resRecord := model.RecordResponse{
		ID:          newRecord.ID,
		Artist:      newRecord.Artist,
		Title:       newRecord.Title,
		Genre:       newRecord.Genre,
		Style:       newRecord.Style,
		ReleaseYear: newRecord.ReleaseYear,
	}
	return resRecord, nil
}

func (ru *recordUsecase) GetRecordList() ([]model.RecordResponse, error) {
	// else句を嫌ったパターン
	// var recordList []model.RecordResponse
	// recordList, err := ru.rr.GetRecordList()
	// if err != nil {
	// 	return nil, err
	// }
	// return recordList, nil
	if recordList, err := ru.rr.GetRecordList(); err != nil {
		return nil, err
	} else {
		// 複数の構造体変換は、mapstructure や automapper 等のライブラリを使う
		// forループで呼んでみたパターン
		var recordResponseList []model.RecordResponse
		for _, record := range recordList {
			recordResponseList = append(recordResponseList, model.RecordResponse{
				ID:          record.ID,
				Title:       record.Title,
				Artist:      record.Artist,
				Genre:       record.Genre,
				Style:       record.Style,
				ReleaseYear: record.ReleaseYear,
			})
		}
		return recordResponseList, nil
	}
}

func (ru *recordUsecase) GetDetail(title string) (model.DetailResponse, error) {
	if record, err := ru.rr.GetDetail(title); err != nil {
		return model.DetailResponse{}, err
	} else {
		return record, nil
	}
}

func (ru *recordUsecase) GetRecordByTitle(title string) ([]model.RecordResponse, error) {
	if recordList, err := ru.rr.GetRecordByTitle(title); err != nil {
		return nil, err
	} else {
		// mapSlice()を作成し、そこから呼んでみたパターン
		return ru.mapSlice(recordList)
	}
}

func (ru *recordUsecase) GetRecordByArtist(artist string) ([]model.RecordResponse, error) {
	if recordList, err := ru.rr.GetRecordByArtist(artist); err != nil {
		return nil, err
	} else {
		// mapSlice()を作成し、そこから呼んでみたパターン
		return ru.mapSlice(recordList)
	}
}

// model.Recordをmodel.RecordResponseに変換
func (*recordUsecase) mapSlice(recordList []model.Record) ([]model.RecordResponse, error) {
	recordResponseList := common.MapSlice(recordList, func(record model.Record) model.RecordResponse {
		return model.RecordResponse{
			Title:       record.Title,
			Artist:      record.Artist,
			Genre:       record.Genre,
			Style:       record.Style,
			ReleaseYear: record.ReleaseYear,
		}
	})
	return recordResponseList, nil
}

func (ru *recordUsecase) UpdateRecord(record model.Record) (model.RecordResponse, error) {
	if err := ru.rr.UpdateRecord(&record); err != nil {
		return model.RecordResponse{}, err
	}
	// Recordリポジトリが成功の場合、引数で渡したアドレスが指し示す先の値が
	// 更新したRecordで書き変わっているので、そこから新しいRecordResponse構造体を作成して返却
	resTask := model.RecordResponse{
		ID:          record.ID,
		Title:       record.Title,
		Artist:      record.Artist,
		Genre:       record.Genre,
		Style:       record.Style,
		ReleaseYear: record.ReleaseYear,
	}
	return resTask, nil
}

func (ru *recordUsecase) DeleteRecord(id uint) error {
	if err := ru.rr.DeleteRecord(id); err != nil {
		return err
	}
	return nil
}
