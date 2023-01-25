package internal_test

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/matansh/enver/internal"
)

// setting this test up as a table test requires nasty reflections.
func TestPopulateStructFromEnv(t *testing.T) {
	t.Parallel()
	type innerConfig struct {
		InnerStr     string        `testTag:"INNER_STR"`
		InnerInt     int           `testTag:"INNER_INT"`
		InnerInt8    int8          `testTag:"INNER_INT8"`
		InnerInt16   int16         `testTag:"INNER_INT16"`
		InnerInt32   int32         `testTag:"INNER_INT32"`
		InnerInt64   int64         `testTag:"INNER_INT64"`
		InnerUint    uint          `testTag:"INNER_UINT"`
		InnerUint8   uint8         `testTag:"INNER_UINT8"`
		InnerUint16  uint16        `testTag:"INNER_UINT16"`
		InnerUint32  uint32        `testTag:"INNER_UINT32"`
		InnerUint64  uint64        `testTag:"INNER_UINT64"`
		InnerFloat32 float32       `testTag:"INNER_FLOAT32"`
		InnerFloat64 float64       `testTag:"INNER_FLOAT64"`
		InnerBool    bool          `testTag:"INNER_BOOL"`
		InnerTimeout time.Duration `testTag:"INNER_TIMEOUT"`
		InnerArrStr  []string      `testTag:"INNER_ARR_STR"`
	}
	type config struct {
		Str     string        `testTag:"STR"`
		Int     int           `testTag:"INT"`
		Int8    int8          `testTag:"INT8"`
		Int16   int16         `testTag:"INT16"`
		Int32   int32         `testTag:"INT32"`
		Int64   int64         `testTag:"INT64"`
		Uint    uint          `testTag:"UINT"`
		Uint8   uint8         `testTag:"UINT8"`
		Uint16  uint16        `testTag:"UINT16"`
		Uint32  uint32        `testTag:"UINT32"`
		Uint64  uint64        `testTag:"UINT64"`
		Float32 float32       `testTag:"FLOAT32"`
		Float64 float64       `testTag:"FLOAT64"`
		Bool    bool          `testTag:"BOOL"`
		Timeout time.Duration `testTag:"TIMEOUT"`
		ArrStr  []string      `testTag:"ARR_STR"`
		Inner   innerConfig
	}
	testData := map[string]string{
		"STR":           "test ",
		"INT":           fmt.Sprint(int(1)),
		"INT8":          fmt.Sprint(int8(2)),
		"INT16":         fmt.Sprint(int16(3)),
		"INT32":         fmt.Sprint(int32(4)),
		"INT64":         fmt.Sprint(int64(5)),
		"UINT":          fmt.Sprint(uint(6)),
		"UINT8":         fmt.Sprint(uint8(7)),
		"UINT16":        fmt.Sprint(uint16(8)),
		"UINT32":        fmt.Sprint(uint32(9)),
		"UINT64":        fmt.Sprint(uint64(10)),
		"FLOAT32":       fmt.Sprint(float32(1.1)),
		"FLOAT64":       fmt.Sprint(float64(1.2)),
		"BOOL":          fmt.Sprint(true),
		"TIMEOUT":       "1m",
		"ARR_STR":       "first,second, third",
		"INNER_STR":     "\"inner test\"",
		"INNER_INT":     fmt.Sprint(int(13)),
		"INNER_INT8":    fmt.Sprint(int8(14)),
		"INNER_INT16":   fmt.Sprint(int16(15)),
		"INNER_INT32":   fmt.Sprint(int32(16)),
		"INNER_INT64":   fmt.Sprint(int64(17)),
		"INNER_UINT":    fmt.Sprint(uint(18)),
		"INNER_UINT8":   fmt.Sprint(uint8(19)),
		"INNER_UINT16":  fmt.Sprint(uint16(20)),
		"INNER_UINT32":  fmt.Sprint(uint32(21)),
		"INNER_UINT64":  fmt.Sprint(uint64(22)),
		"INNER_FLOAT32": fmt.Sprint(float32(23.6)),
		"INNER_FLOAT64": fmt.Sprint(float64(24.7)),
		"INNER_BOOL":    fmt.Sprint(false),
		"INNER_TIMEOUT": "10",
		"INNER_ARR_STR": "fourth,fifth, sixth",
	}
	// setting data into env
	for key, value := range testData {
		if err := os.Setenv(key, value); err != nil {
			t.Fatal(err) // stopping test execution
		}
	}
	// populating test struct
	var cfg config
	if err := internal.ScanStructFromEnv(&cfg, "testTag"); err != nil {
		t.Fatal(err) // stopping test execution
	}
	// assertions - these cant be in a table test cuz we need to read the value assigned to the struct
	isEqual := reflect.DeepEqual(cfg, config{
		Str:     "test",
		Int:     int(1),
		Int8:    int8(2),
		Int16:   int16(3),
		Int32:   int32(4),
		Int64:   int64(5),
		Uint:    uint(6),
		Uint8:   uint8(7),
		Uint16:  uint16(8),
		Uint32:  uint32(9),
		Uint64:  uint64(10),
		Float32: float32(1.1),
		Float64: float64(1.2),
		Bool:    true,
		Timeout: time.Minute,
		ArrStr:  []string{"first", "second", "third"},
		Inner: innerConfig{
			InnerStr:     "inner test",
			InnerInt:     int(13),
			InnerInt8:    int8(14),
			InnerInt16:   int16(15),
			InnerInt32:   int32(16),
			InnerInt64:   int64(17),
			InnerUint:    uint(18),
			InnerUint8:   uint8(19),
			InnerUint16:  uint16(20),
			InnerUint32:  uint32(21),
			InnerUint64:  uint64(22),
			InnerFloat32: float32(23.6),
			InnerFloat64: float64(24.7),
			InnerBool:    false,
			InnerTimeout: time.Second * 10,
			InnerArrStr:  []string{"fourth", "fifth", "sixth"},
		},
	})
	if !isEqual {
		t.Errorf("cfg was not in its wanted state: %+v", cfg)
	}
}
