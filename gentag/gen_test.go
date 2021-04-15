package gentag

import (
	"reflect"
	"testing"
)

func TestGen(t *testing.T) {
	type args struct {
		data      []byte
		tagtype   string
		omitempty bool
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		want1      string
		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation after test
	}{
		{
			name: "empty",
			args: func(t *testing.T) args {
				return args{
					data:      []byte{},
					tagtype:   "json",
					omitempty: false,
				}
			},
			want1:   "",
			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				if err.Error() != "unexpected end of JSON input" {
					t.FailNow()
				}
			},
		},
		{
			name: "empty struct",
			args: func(t *testing.T) args {
				return args{
					data:      []byte(`{}`),
					tagtype:   "json",
					omitempty: false,
				}
			},
			want1:   "type Struct struct {\n}\n\n",
			wantErr: false,
		},
		{
			name: "empty array",
			args: func(t *testing.T) args {
				return args{
					data:      []byte(`[]`),
					tagtype:   "json",
					omitempty: false,
				}
			},
			want1:   "",
			wantErr: false,
		},
		{
			name: "num array",
			args: func(t *testing.T) args {
				return args{
					data:      []byte(`[1, 2]`),
					tagtype:   "json",
					omitempty: false,
				}
			},
			want1:   "",
			wantErr: false,
		},
		{
			name: "struct1",
			args: func(t *testing.T) args {
				return args{
					data:      []byte(`{"hello_world": 1}`),
					tagtype:   "json",
					omitempty: true,
				}
			},
			want1:   "type Struct struct {\n\tHelloWorld float64 `json:\"hello_world,omitempty\"`\n}\n\n",
			wantErr: false,
		},
		{
			name: "struct2",
			args: func(t *testing.T) args {
				return args{
					data: []byte(`
{
    "hello": {
        "aaa1": 111,
        "bbb2": 222,
        "ccc3": 333
    },
    "world": {
        "aaa1": 0,
        "bbb2": 1,
        "ccc3": 2
    },
	"a_bool": true,
	"hello_world": "hello_world"
}
`),
					tagtype:   "json",
					omitempty: false,
				}
			},
			want1:   "type Hello struct {\n\tAaa1 float64 `json:\"aaa1\"`\n\tBbb2 float64 `json:\"bbb2\"`\n\tCcc3 float64 `json:\"ccc3\"`\n}\n\ntype Struct struct {\n\tABool bool `json:\"a_bool\"`\n\tHello *Hello `json:\"hello\"`\n\tHelloWorld string `json:\"hello_world\"`\n\tWorld *Hello `json:\"world\"`\n}\n\n",
			wantErr: false,
		},
		{
			name: "array",
			args: func(t *testing.T) args {
				return args{
					data: []byte(`
[{
	"hello": "world"
}
]
`),
					tagtype: "json",
				}
			},
			want1: "type Struct struct {\n\tHello string `json:\"hello\"`\n}\n\n",
		},
		{
			name: "nil",
			args: func(t *testing.T) args {
				return args{
					data: []byte(`
{
	"hello": null,
	"nil_array": [],
	"sub_array": [[{
	"a_field": 123
	}]]
}
`),
					tagtype: "json",
				}
			},
			want1: "type Struct struct {\n\tHello interface{} `json:\"hello\"`\n\tNilArray []interface{} `json:\"nil_array\"`\n\tSubArray [][]*Struct81 `json:\"sub_array\"`\n}\n\ntype Struct81 struct {\n\tAField float64 `json:\"a_field\"`\n}\n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			got1, err := Gen(tArgs.data, tArgs.tagtype, tArgs.omitempty)

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Gen got1 = `%v`, want1: `%v`", got1, tt.want1)
			}

			if (err != nil) != tt.wantErr {
				t.Fatalf("Gen error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}
