package parsex_test

import (
	"errors"
	"testing"

	"github.com/bbfh-dev/parsex/v2"
	"gotest.tools/assert"
)

func TestRuntimeErrorCases(t *testing.T) {
	cases := []struct {
		name        string
		program     parsex.Program
		programArgs []string
		args        []string
		wantErrKind parsex.ErrKind
		wantErrType error
	}{
		{
			name:        "DataNotPointer",
			program:     parsex.Program{Data: 123, Name: "", Desc: "", Exec: nil},
			programArgs: []string{},
			args:        []string{},
			wantErrType: parsex.ErrProgramData{},
			wantErrKind: parsex.ErrKindMustbePointer,
		},
		{
			name:        "DataPointerNotStruct",
			program:     func() parsex.Program { d := 123; return parsex.Program{Data: &d, Name: "", Desc: "", Exec: nil} }(),
			programArgs: []string{},
			args:        []string{},
			wantErrType: parsex.ErrProgramData{},
			wantErrKind: parsex.ErrKindPointToStruct,
		},
		{
			name:        "ExecIsNil",
			program:     parsex.Program{Data: nil, Name: "", Desc: "", Exec: nil},
			programArgs: []string{},
			args:        []string{},
			wantErrType: parsex.ErrExecution{},
			wantErrKind: parsex.ErrKindExecIsNil,
		},
		{
			name: "ExecReturnsError",
			program: parsex.Program{
				Data: nil,
				Name: "",
				Desc: "",
				Exec: func(_ []string) error { return errors.New("fail") },
			},
			programArgs: []string{},
			args:        []string{},
			wantErrType: parsex.ErrExecution{},
			wantErrKind: parsex.ErrKindExecution,
		},
		{
			name: "NotEnoughArgs",
			program: parsex.Program{
				Data: nil,
				Name: "",
				Desc: "",
				Exec: func(args []string) error { return nil },
			},
			programArgs: []string{"arg1"},
			args:        []string{},
			wantErrType: parsex.ErrInput{},
			wantErrKind: parsex.ErrKindNotEnoughArgs,
		},
		{
			name: "UnknownOption",
			program: parsex.Program{
				Data: nil,
				Name: "",
				Desc: "",
				Exec: func(args []string) error { return nil },
			},
			programArgs: []string{},
			args:        []string{"--unknown"},
			wantErrType: parsex.ErrOption{},
			wantErrKind: parsex.ErrKindUnknownOption,
		},
		{
			name: "UnknownOption",
			program: func() parsex.Program {
				var data struct {
					Unknown string
				}
				return parsex.Program{
					Data: &data,
					Name: "",
					Desc: "",
					Exec: func(args []string) error { return nil },
				}
			}(),
			programArgs: []string{},
			args:        []string{"--unknown"},
			wantErrType: parsex.ErrOption{},
			wantErrKind: parsex.ErrKindOptionNeedsValue,
		},
		{
			name: "SettingOption",
			program: func() parsex.Program {
				var data struct {
					Unknown int
				}
				return parsex.Program{
					Data: &data,
					Name: "",
					Desc: "",
					Exec: func(args []string) error { return nil },
				}
			}(),
			programArgs: []string{},
			args:        []string{"--unknown", "abc"},
			wantErrType: parsex.ErrOption{},
			wantErrKind: parsex.ErrKindSettingOption,
		},
		{
			name: "UnknownCluster",
			program: func() parsex.Program {
				var data struct {
					A bool `alt:"a"`
					B bool `alt:"b"`
				}
				return parsex.Program{
					Data: &data,
					Name: "",
					Desc: "",
					Exec: func(args []string) error { return nil },
				}
			}(),
			programArgs: []string{},
			args:        []string{"-abc"},
			wantErrType: parsex.ErrOption{},
			wantErrKind: parsex.ErrKindUnknownCluster,
		},
		{
			name: "MistypedCluster",
			program: func() parsex.Program {
				var data struct {
					A bool   `alt:"a"`
					B string `alt:"b"`
				}
				return parsex.Program{
					Data: &data,
					Name: "",
					Desc: "",
					Exec: func(args []string) error { return nil },
				}
			}(),
			programArgs: []string{},
			args:        []string{"-ab"},
			wantErrType: parsex.ErrOption{},
			wantErrKind: parsex.ErrKindMistypedCluster,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			runtime := testCase.program.Runtime().SetPosArgs(testCase.programArgs...)
			err := runtime.Run(testCase.args)

			if testCase.wantErrType == nil {
				assert.NilError(t, err)
				return
			}

			switch err := err.(type) {
			case parsex.ErrProgramData:
				assert.Equal(t, err.ErrKind, testCase.wantErrKind)
			case parsex.ErrExecution:
				assert.Equal(t, err.ErrKind, testCase.wantErrKind)
			case parsex.ErrInput:
				assert.Equal(t, err.ErrKind, testCase.wantErrKind)
			case parsex.ErrOption:
				assert.Equal(t, err.ErrKind, testCase.wantErrKind)

			default:
				t.Fatalf("Test %q: unexpected error type: %T", testCase.name, err)
			}
		})
	}
}
