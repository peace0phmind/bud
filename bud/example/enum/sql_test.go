package enum

import (
	"github.com/peace0phmind/bud/util/opt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSQLExtras(t *testing.T) {
	assert.Equal(t, "ProjectStatus(22).Name", ProjectStatus(22).String(), "String value is not correct")

	_, err := ParseProjectStatus(`NotAStatus`)
	assert.Error(t, err, "Should have had an error parsing a non status")

	var (
		intVal            int    = 3
		strVal            string = "completed"
		strIntVal         string = "2"
		nullInt           *int
		nullInt64         *int64
		nullFloat64       *float64
		nullUint          *uint
		nullUint64        *uint64
		nullString        *string
		nullProjectStatus *ProjectStatus
	)

	tests := map[string]struct {
		input  interface{}
		result *opt.SqlNull[ProjectStatus]
	}{
		"nil": {},
		"val": {
			input:  ProjectStatusRejected,
			result: opt.NewSqlNull[ProjectStatus](ProjectStatusRejected),
		},
		"ptr": {
			input:  ProjectStatusCompleted.Ptr(),
			result: opt.NewSqlNull[ProjectStatus](ProjectStatusCompleted),
		},
		"string": {
			input:  strVal,
			result: opt.NewSqlNull[ProjectStatus](ProjectStatusCompleted),
		},
		"*string": {
			input:  &strVal,
			result: opt.NewSqlNull[ProjectStatus](ProjectStatusCompleted),
		},
		"*string as int": {
			input:  &strIntVal,
			result: opt.NewSqlNull[ProjectStatus](ProjectStatusCompleted),
		},
		"invalid string": {
			input: "random value",
		},
		"[]byte": {
			input:  []byte(ProjectStatusInWork.String()),
			result: opt.NewSqlNull[ProjectStatus](ProjectStatusInWork),
		},
		"int": {
			input:  intVal,
			result: opt.NewSqlNull[ProjectStatus](ProjectStatusRejected),
		},
		"*int": {
			input:  &intVal,
			result: opt.NewSqlNull[ProjectStatus](ProjectStatusRejected),
		},
		"nullInt": {
			input: nullInt,
		},
		"nullInt64": {
			input: nullInt64,
		},
		"nullUint": {
			input: nullUint,
		},
		"nullUint64": {
			input: nullUint64,
		},
		"nullFloat64": {
			input: nullFloat64,
		},
		"nullString": {
			input: nullString,
		},
		"nullProjectStatus": {
			input: nullProjectStatus,
		},
		"int as []byte": {
			input:  []byte("1"),
			result: opt.NewSqlNull[ProjectStatus](ProjectStatusInWork),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			status := &opt.SqlNull[ProjectStatus]{}
			err1 := status.Scan(tc.input)
			if err1 != nil {
				t.Errorf("scan err: %v", err1)
			}
			assert.Equal(t, tc.result, status)
		})
	}
}

//type SQLMarshalType struct {
//	Status    NullProjectStatus    `json:"status"`
//	StatusStr NullProjectStatusStr `json:"status_str"`
//	Status2   *ProjectStatus       `json:"status2,omitempty"`
//}
//
//func TestSQLMarshal(t *testing.T) {
//	var val SQLMarshalType
//	var val2 SQLMarshalType
//
//	result, err := json.Marshal(val)
//	require.NoError(t, err)
//	assert.Equal(t, `{"status":null,"status_str":null}`, string(result))
//
//	require.NoError(t, json.Unmarshal([]byte(`{}`), &val2))
//	assert.Equal(t, val, val2)
//
//	require.NoError(t, json.Unmarshal(result, &val2))
//	val.Status.Set = true
//	val.StatusStr.Set = true
//	assert.Equal(t, val, val2)
//
//	val.Status = NewNullProjectStatus(1)
//	val.StatusStr = NewNullProjectStatusStr(2)
//	result, err = json.Marshal(val)
//	require.NoError(t, err)
//	assert.Equal(t, `{"status":"inWork","status_str":"completed"}`, string(result))
//
//	require.NoError(t, json.Unmarshal(result, &val2))
//	assert.Equal(t, val, val2)
//
//	require.NoError(t, json.Unmarshal([]byte(`{"status":"inWork"}`), &val2))
//	assert.Equal(t, val, val2)
//
//	require.NoError(t, json.Unmarshal([]byte(`{"status":"2"}`), &val2))
//	val.Status = NewNullProjectStatus(2)
//	assert.Equal(t, val, val2)
//
//	require.NoError(t, json.Unmarshal([]byte(`{"status":3}`), &val2))
//	val.Status = NewNullProjectStatus(3)
//	assert.Equal(t, val, val2)
//
//	require.NoError(t, json.Unmarshal([]byte(`{"status":null}`), &val2))
//	val.Status = NullProjectStatus{Set: true}
//	assert.Equal(t, val, val2)
//
//	val2 = SQLMarshalType{} // reset it so that the `set` value is false.
//
//	require.NoError(t, json.Unmarshal([]byte(`{"status2":"rejected"}`), &val2))
//	val.Status = NullProjectStatus{}
//	val.StatusStr = NullProjectStatusStr{}
//	val.Status2 = ProjectStatusRejected.Ptr()
//	assert.Equal(t, val, val2)
//
//	require.Error(t, json.Unmarshal([]byte(`{"status2":"xyz"}`), &val2))
//	val2 = SQLMarshalType{} // reset it so that the `set` value is false.
//
//	require.NoError(t, json.Unmarshal([]byte(`{"status_str":"rejected"}`), &val2))
//	val.Status = NullProjectStatus{}
//	val.StatusStr = NewNullProjectStatusStr(3)
//	val.Status2 = nil
//	assert.Equal(t, val, val2)
//
//	require.Error(t, json.Unmarshal([]byte(`{"status2":"xyz"}`), &val2))
//}
