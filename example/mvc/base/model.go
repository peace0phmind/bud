package base

import (
	"errors"
	"github.com/google/uuid"
	"github.com/peace0phmind/bud/util"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

type IdInf interface {
	SetId(id any) error
}

type UUIDBase struct {
	Id        *uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`
	CreatedAt time.Time  `json:"createdAt" gorm:"type:datetime(3);"`
	UpdatedAt time.Time  `json:"updatedAt" gorm:"type:datetime(3);"`
}

const EmptyUUID = "00000000-0000-0000-0000-000000000000"

func (ub *UUIDBase) BeforeCreate(*gorm.DB) error {
	if ub.Id == nil || ub.Id.String() == EmptyUUID {
		uId := uuid.New()
		ub.Id = &uId
		//tx.Statement.SetColumn("id", ub.Id)
		//} else {
		//	tx.Statement.SetColumn("id", ub.Id)
	}

	ub.CreatedAt = time.Time{}
	ub.UpdatedAt = time.Time{}

	return nil
}

func (ub *UUIDBase) SetId(id any) (err error) {
	if util.IsNil(id) {
		ub.Id = nil
		return
	}

	if sId, ok := id.(string); ok {
		if uId, err1 := uuid.Parse(sId); err1 != nil {
			logrus.Errorf("base uuid parse: %+v", err1)
			return err1
		} else {
			ub.Id = &uId
			return nil
		}
	} else if uId, ok := id.(uuid.UUID); ok {
		ub.Id = &uId
		return nil
	} else if puid, ok := id.(*uuid.UUID); ok {
		ub.Id = puid
		return nil
	} else {
		return errors.New("id must be string or uuid.UUID")
	}
}
