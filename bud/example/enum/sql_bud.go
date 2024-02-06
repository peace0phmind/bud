package enum

import (
	"database/sql/driver"
	"errors"
	"fmt"
)

const (
	ProjectStatusPending ProjectStatus = iota
	ProjectStatusInWork
	ProjectStatusCompleted
	ProjectStatusRejected
)

var ErrInvalidProjectStatus = errors.New("not a valid ProjectStatus")

var _ProjectStatusName = "pendinginWorkcompletedrejected"

var _ProjectStatusMapName = map[ProjectStatus]string{
	ProjectStatusPending:   _ProjectStatusName[0:7],
	ProjectStatusInWork:    _ProjectStatusName[7:13],
	ProjectStatusCompleted: _ProjectStatusName[13:22],
	ProjectStatusRejected:  _ProjectStatusName[22:30],
}

// Name is the attribute of ProjectStatus.
func (x ProjectStatus) Name() string {
	if v, ok := _ProjectStatusMapName[x]; ok {
		return v
	}
	panic(ErrInvalidProjectStatus)
}

// Val is the attribute of ProjectStatus.
func (x ProjectStatus) Val() int {
	if x.IsValid() {
		return int(x)
	}
	panic(ErrInvalidProjectStatus)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x ProjectStatus) IsValid() bool {
	_, ok := _ProjectStatusMapName[x]
	return ok
}

// String implements the Stringer interface.
func (x ProjectStatus) String() string {
	return x.Name()
}

var _ProjectStatusNameMap = map[string]ProjectStatus{
	_ProjectStatusName[0:7]:   ProjectStatusPending,
	_ProjectStatusName[7:13]:  ProjectStatusInWork,
	_ProjectStatusName[13:22]: ProjectStatusCompleted,
	_ProjectStatusName[22:30]: ProjectStatusRejected,
}

// ParseProjectStatus converts a string to a ProjectStatus.
func ParseProjectStatus(value string) (ProjectStatus, error) {
	if x, ok := _ProjectStatusNameMap[value]; ok {
		return x, nil
	}
	return ProjectStatus(0), fmt.Errorf("%s is %w", value, ErrInvalidProjectStatus)
}

func (x ProjectStatus) Ptr() *ProjectStatus {
	return &x
}

// MarshalText implements the text marshaller method.
func (x ProjectStatus) MarshalText() ([]byte, error) {
	return []byte(x.String()), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *ProjectStatus) UnmarshalText(text []byte) error {
	val, err := ParseProjectStatus(string(text))
	if err != nil {
		return err
	}
	*x = val
	return nil
}

var errProjectStatusNilPtr = errors.New("value pointer is nil")

// Scan implements the Scanner interface.
func (x *ProjectStatus) Scan(value any) (err error) {
	if value == nil {
		*x = ProjectStatus(0)
		return
	}

	switch v := value.(type) {
	case int:
		*x = ProjectStatus(v)
	case int64:
		*x = ProjectStatus(v)
	case uint:
		*x = ProjectStatus(v)
	case uint64:
		*x = ProjectStatus(v)
	case float64:
		*x = ProjectStatus(v)
	case *int:
		if v == nil {
			return errProjectStatusNilPtr
		}
		*x = ProjectStatus(*v)
	case *int64:
		if v == nil {
			return errProjectStatusNilPtr
		}
		*x = ProjectStatus(*v)
	case *uint:
		if v == nil {
			return errProjectStatusNilPtr
		}
		*x = ProjectStatus(*v)
	case *uint64:
		if v == nil {
			return errProjectStatusNilPtr
		}
		*x = ProjectStatus(*v)
	case *float64:
		if v == nil {
			return errProjectStatusNilPtr
		}
		*x = ProjectStatus(*v)
	}

	if !x.IsValid() {
		return ErrInvalidProjectStatus
	}
	return
}

// Value implements the driver Valuer interface.
func (x ProjectStatus) Value() (driver.Value, error) {
	return x.Val(), nil
}
