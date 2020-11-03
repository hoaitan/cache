package local

import (
	"fmt"
	"reflect"
	"runtime"
	"testing"

	"github.com/hoaitan/cache/test"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		cf   Config
	}{
		{
			"Cache with correct size: 1MB",
			Config{
				Size: 1000000,
			},
		},
		{
			"Cache with invalid size",
			Config{
				Size: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(tt.cf)
			assert.NotNil(t, c)
		})
	}
}

func TestCacheImplement_Enable(t *testing.T) {
	c := New(Config{
		Enable: true,
		Size:   1000000,
	})

	for _, fn := range test.GetTestSuite(true) {
		t.Run(fmt.Sprintf("fn=%s", runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()), func(t *testing.T) {
			fn(t, c)
		})

	}

}
func TestCacheImplement_Disable(t *testing.T) {
	c := New(Config{
		Enable: false,
		Size:   1000000,
	})

	for _, fn := range test.GetTestSuite(false) {
		t.Run(fmt.Sprintf("fn=%s", runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()), func(t *testing.T) {
			fn(t, c)
		})

	}

}

func TestEncodeAndDecode(t *testing.T) {
	var (
		// Bool
		boolData, boolData2 bool = true, false

		// String
		stringData, stringData2 string = "test", ""

		// Int
		intData, intData2     int   = 1, 0
		int8Data, int8Data2   int8  = 1, 0
		int16Data, int16Data2 int32 = 1, 0
		int32Data, int32Data2 int32 = 1, 0
		int64Data, int64Data2 int64 = 1, 0

		uintData, uintData2     uint   = 1, 0
		uint8Data, uint8Data2   uint8  = 1, 0
		uint16Data, uint16Data2 uint32 = 1, 0
		uint32Data, uint32Data2 uint32 = 1, 0
		uint64Data, uint64Data2 uint64 = 1, 0

		byteData, byteData2 byte = 1, 0
		runeData, runeData2 rune = 1, 0

		// Float
		float32Data, float32Data2 float32 = 1, 0
		float64Data, float64Data2 float64 = 1, 0

		// Complex
		complex64Data, complex64Data2   complex64  = 1, 0
		complex128Data, complex128Data2 complex128 = 1, 0

		// Slice
		sliceData, sliceData2 []int = []int{1}, []int{0}

		// Map
		mapData, mapData2 map[string]int = map[string]int{"test": 1}, map[string]int{"test": 0}

		// Struct & Pointer
	)
	tests := []struct {
		data              interface{}
		ptrToSameDataType interface{}
	}{
		{
			boolData,
			&boolData2,
		},
		{
			stringData,
			&stringData2,
		},
		// Int
		{
			intData,
			&intData2,
		},
		{
			int8Data,
			&int8Data2,
		},
		{
			int16Data,
			&int16Data2,
		},
		{
			int32Data,
			&int32Data2,
		},
		{
			int64Data,
			&int64Data2,
		},
		{
			uintData,
			&uintData2,
		},
		{
			uint8Data,
			&uint8Data2,
		},
		{
			uint16Data,
			&uint16Data2,
		},
		{
			uint32Data,
			&uint32Data2,
		},
		{
			uint64Data,
			&uint64Data2,
		},
		{
			byteData,
			&byteData2,
		},
		{
			runeData,
			&runeData2,
		},

		// Float
		{
			float32Data,
			&float32Data2,
		},
		{
			float64Data,
			&float64Data2,
		},

		// Complex
		{
			complex64Data,
			&complex64Data2,
		},
		{
			complex128Data,
			&complex128Data2,
		},

		// Slice
		{
			sliceData,
			&sliceData2,
		},

		// Map
		{
			mapData,
			&mapData2,
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("With type=%T", tt.data), func(t *testing.T) {
			// Set
			b, err := encode(tt.data)
			assert.Nil(t, err)

			// Get
			err = decode(b, tt.ptrToSameDataType)
			assert.Nil(t, err)
			assert.Equal(t, tt.data, reflect.Indirect(reflect.ValueOf(tt.ptrToSameDataType)).Interface())
		})
	}
}
