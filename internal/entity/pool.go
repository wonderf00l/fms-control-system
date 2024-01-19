package entity

import (
	"github.com/wonderf00l/fms-control-system/internal/entity/conveyor"
	"github.com/wonderf00l/fms-control-system/internal/entity/lathe"
	"github.com/wonderf00l/fms-control-system/internal/entity/miller"
	"github.com/wonderf00l/fms-control-system/internal/entity/recognition"
	"github.com/wonderf00l/fms-control-system/internal/entity/storage"
)

type Pool struct {
	Storage     storage.Storage
	Conveyor    conveyor.Conveyor
	Recognition recognition.Recognition
	Lathe       lathe.Lathe
	Miller      miller.Miller
}
