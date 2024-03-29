package homework

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	type args struct {
		v any
	}
	type TestStruct struct {
		LenS      []string `validate:"len:2"`
		MaxI      []int    `validate:"max:3"`
		MaxS      []string `validate:"max:5"`
		MinI      []int    `validate:"min:2"`
		MinS      []string `validate:"min:2"`
		InI       []int    `validate:"in:10,25,30"`
		InS       []string `validate:"in:foo,bar"`
		Len       string   `validate:"len:20"`
		LenZ      string   `validate:"len:0"`
		InInt     int      `validate:"in:20,25,30"`
		InNeg     int      `validate:"in:-20,-25,-30"`
		InStr     string   `validate:"in:foo,bar"`
		MinInt    int      `validate:"min:10"`
		MinIntNeg int      `validate:"min:-10"`
		MinStr    string   `validate:"min:10"`
		MaxInt    int      `validate:"max:20"`
		MaxIntNeg int      `validate:"max:-2"`
		MaxStr    string   `validate:"max:20"`
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		checkErr func(err error) bool
	}{
		{
			name: "valid nested struct",
			args: args{
				v: struct {
					TestStruct TestStruct
					Lower      string `validate:"len:2"`
					Higher     string `validate:"len:5"`
					Zero       string `validate:"len:0"`
				}{
					TestStruct: TestStruct{
						LenS:      []string{"11", "22"},
						MaxI:      []int{1, 2, 3},
						MaxS:      []string{"1", "2", "5"},
						MinI:      []int{11, 22, 11},
						MinS:      []string{"11", "22", "77"},
						InI:       []int{10, 25, 30},
						InS:       []string{"foo", "foo", "bar"},
						Len:       "abcdefghjklmopqrstvu",
						LenZ:      "",
						InInt:     25,
						InNeg:     -25,
						InStr:     "bar",
						MinInt:    15,
						MinIntNeg: -9,
						MinStr:    "abcdefghjkl",
						MaxInt:    16,
						MaxIntNeg: -3,
						MaxStr:    "abcdefghjklmopqrst",
					},
					Lower:  "ab",
					Higher: "abcde",
					Zero:   "",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid nested struct",
			args: args{
				v: struct {
					TestStruct TestStruct
				}{
					TestStruct: TestStruct{
						LenS:      []string{"1", "2"},
						MaxI:      []int{1, 2, 7},
						MaxS:      []string{"111111", "2", "5"},
						MinI:      []int{1, 2, 1},
						MinS:      []string{"1", "2", "7"},
						InI:       []int{10, 178, -7},
						InS:       []string{"ooo", "bar", "aaa"},
						Len:       "abcdefghjklmopqrstvuwxyz",
						LenZ:      "1",
						InInt:     1000,
						InNeg:     -2500,
						InStr:     "str",
						MinInt:    9,
						MinIntNeg: -90,
						MinStr:    "ab",
						MaxInt:    160,
						MaxIntNeg: -1,
						MaxStr:    "abcdefghjklmopqrstvuwxyz",
					},
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 18)
				return true
			},
		},
		{
			name: "valid slice",
			args: args{
				v: struct {
					LenS []string `validate:"len:2"`
					MaxI []int    `validate:"max:3"`
					MaxS []string `validate:"max:5"`
					MinI []int    `validate:"min:1"`
					MinS []string `validate:"min:1"`
					InI  []int    `validate:"in:10,25,30"`
					InS  []string `validate:"in:foo,bar"`
				}{
					LenS: []string{"11", "22"},
					MaxI: []int{1, 2, 3},
					MaxS: []string{"1", "2", "5"},
					MinI: []int{1, 2, 1},
					MinS: []string{"1", "2", "7"},
					InI:  []int{10, 25, 30},
					InS:  []string{"foo", "foo", "bar"},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid slice",
			args: args{
				v: struct {
					LenS []string `validate:"len:2"`
					MaxI []int    `validate:"max:1"`
					MaxS []string `validate:"max:1"`
					MinI []int    `validate:"min:2"`
					MinS []string `validate:"min:7"`
					InI  []int    `validate:"in:10,25,30"`
					InS  []string `validate:"in:foo,bar"`
				}{
					LenS: []string{"1", "2"},
					MaxI: []int{1, 2, 3},
					MaxS: []string{"11", "2", "5"},
					MinI: []int{1, 2, 1},
					MinS: []string{"1", "2", "7"},
					InI:  []int{10, 178, -7},
					InS:  []string{"ooo", "bar", "aaa"},
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 7)
				return true
			},
		},
		{
			name: "invalid struct: interface",
			args: args{
				v: new(any),
			},
			wantErr: true,
			checkErr: func(err error) bool {
				return errors.Is(err, ErrNotStruct)
			},
		},
		{
			name: "invalid struct: map",
			args: args{
				v: map[string]string{},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				return errors.Is(err, ErrNotStruct)
			},
		},
		{
			name: "invalid struct: string",
			args: args{
				v: "some string",
			},
			wantErr: true,
			checkErr: func(err error) bool {
				return errors.Is(err, ErrNotStruct)
			},
		},
		{
			name: "valid struct with no fields",
			args: args{
				v: struct{}{},
			},
			wantErr: false,
		},
		{
			name: "valid struct with untagged fields",
			args: args{
				v: struct {
					f1 string
					f2 string
				}{},
			},
			wantErr: false,
		},
		{
			name: "valid struct with unexported fields",
			args: args{
				v: struct {
					foo string `validate:"len:10"`
				}{},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				e := &ValidationErrors{}
				return errors.As(err, e) && e.Error() == ErrValidateForUnexportedFields.Error()
			},
		},
		{
			name: "invalid validator syntax",
			args: args{
				v: struct {
					Foo string `validate:"len:abcdef"`
				}{},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				e := &ValidationErrors{}
				return errors.As(err, e) && e.Error() == ErrInvalidValidatorSyntax.Error()
			},
		},
		{
			name: "valid struct with tagged fields",
			args: args{
				v: struct {
					Len       string `validate:"len:20"`
					LenZ      string `validate:"len:0"`
					InInt     int    `validate:"in:20,25,30"`
					InNeg     int    `validate:"in:-20,-25,-30"`
					InStr     string `validate:"in:foo,bar"`
					MinInt    int    `validate:"min:10"`
					MinIntNeg int    `validate:"min:-10"`
					MinStr    string `validate:"min:10"`
					MinStrNeg string `validate:"min:-1"`
					MaxInt    int    `validate:"max:20"`
					MaxIntNeg int    `validate:"max:-2"`
					MaxStr    string `validate:"max:20"`
				}{
					Len:       "abcdefghjklmopqrstvu",
					LenZ:      "",
					InInt:     25,
					InNeg:     -25,
					InStr:     "bar",
					MinInt:    15,
					MinIntNeg: -9,
					MinStr:    "abcdefghjkl",
					MinStrNeg: "abc",
					MaxInt:    16,
					MaxIntNeg: -3,
					MaxStr:    "abcdefghjklmopqrst",
				},
			},
			wantErr: false,
		},
		{
			name: "wrong length",
			args: args{
				v: struct {
					Lower    string `validate:"len:24"`
					Higher   string `validate:"len:5"`
					Zero     string `validate:"len:3"`
					BadSpec  string `validate:"len:%12"`
					Negative string `validate:"len:-6"`
				}{
					Lower:    "abcdef",
					Higher:   "abcdef",
					Zero:     "",
					BadSpec:  "abc",
					Negative: "abcd",
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 5)
				return true
			},
		},
		{
			name: "wrong in",
			args: args{
				v: struct {
					InA     string `validate:"in:ab,cd"`
					InB     string `validate:"in:aa,bb,cd,ee"`
					InC     int    `validate:"in:-1,-3,5,7"`
					InD     int    `validate:"in:5-"`
					InEmpty string `validate:"in:"`
				}{
					InA:     "ef",
					InB:     "ab",
					InC:     2,
					InD:     12,
					InEmpty: "",
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 5)
				return true
			},
		},
		{
			name: "wrong min",
			args: args{
				v: struct {
					MinA string `validate:"min:12"`
					MinB int    `validate:"min:-12"`
					MinC int    `validate:"min:5-"`
					MinD int    `validate:"min:"`
					MinE string `validate:"min:"`
				}{
					MinA: "ef",
					MinB: -22,
					MinC: 12,
					MinD: 11,
					MinE: "abc",
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 5)
				return true
			},
		},
		{
			name: "wrong max",
			args: args{
				v: struct {
					MaxA string `validate:"max:2"`
					MaxB string `validate:"max:-7"`
					MaxC int    `validate:"max:-12"`
					MaxD int    `validate:"max:5-"`
					MaxE int    `validate:"max:"`
					MaxF string `validate:"max:"`
				}{
					MaxA: "efgh",
					MaxB: "ab",
					MaxC: 22,
					MaxD: 12,
					MaxE: 11,
					MaxF: "abc",
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 6)
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.args.v)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, tt.checkErr(err), "test expect an error, but got wrong error type")
			} else {
				assert.NoError(t, err)
			}
		})
	}

}
