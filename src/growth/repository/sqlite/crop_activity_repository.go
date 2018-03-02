package sqlite

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/Tanibox/tania-server/src/growth/repository"
	"github.com/Tanibox/tania-server/src/growth/storage"
)

type CropActivityRepositorySqlite struct {
	DB *sql.DB
}

func NewCropActivityRepositorySqlite(db *sql.DB) repository.CropActivityRepository {
	return &CropActivityRepositorySqlite{DB: db}
}

func (f *CropActivityRepositorySqlite) Save(cropActivity *storage.CropActivity, isUpdate bool) <-chan error {
	result := make(chan error)

	go func() {
		at, err := json.Marshal(repository.ActivityTypeWrapper{
			ActivityName: cropActivity.ActivityType.Code(),
			ActivityData: cropActivity.ActivityType,
		})

		if isUpdate {
			_, err = f.DB.Exec(`UPDATE CROP_ACTIVITY
				SET BATCH_ID = ?, CONTAINER_TYPE = ?, ACTIVITY_TYPE = ?, ACTIVITY_TYPE_CODE = ?,
				CREATED_DATE = ?, DESCRIPTION = ?
				WHERE CROP_UID = ?`,
				cropActivity.BatchID,
				cropActivity.ContainerType,
				at,
				cropActivity.ActivityType.Code(),
				cropActivity.CreatedDate.Format(time.RFC3339),
				cropActivity.Description,
				cropActivity.UID)

			if err != nil {
				result <- err
			}
		} else {
			_, err = f.DB.Exec(`INSERT INTO CROP_ACTIVITY
				(CROP_UID, BATCH_ID, CONTAINER_TYPE, ACTIVITY_TYPE, ACTIVITY_TYPE_CODE, CREATED_DATE, DESCRIPTION)
				VALUES (?, ?, ?, ?, ?, ?, ?)`,
				cropActivity.UID,
				cropActivity.BatchID,
				cropActivity.ContainerType,
				at,
				cropActivity.ActivityType.Code(),
				cropActivity.CreatedDate.Format(time.RFC3339),
				cropActivity.Description)

			if err != nil {
				result <- err
			}
		}

		result <- nil
		close(result)
	}()

	return result
}