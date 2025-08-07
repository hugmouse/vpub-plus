package syntax

import (
	"reflect"
	"testing"
)

type testRenderer struct{}

func (b *testRenderer) Convert(text string, wrap bool) string {
	return text
}

func (b *testRenderer) Name() string {
	return "test"
}

func TestGet(t *testing.T) {
	testRegistry := NewRegistry()

	err := testRegistry.Register("test-get-renderer", &testRenderer{})
	if err != nil {
		t.Error(err)
		return
	}

	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    Renderer
		wantErr bool
	}{
		{
			name: "get existing renderer",
			args: args{
				name: "test-get-renderer",
			},
			want:    &testRenderer{},
			wantErr: false,
		},
		{
			name: "get non-existing renderer should error",
			args: args{
				name: "boom",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testRegistry.Get(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestList(t *testing.T) {
	listRegistry := NewRegistry()
	err := listRegistry.Register("list-renderer-1", &testRenderer{})
	if err != nil {
		t.Error(err)
		return
	}

	err = listRegistry.Register("list-renderer-2", &testRenderer{})
	if err != nil {
		t.Error(err)
		return
	}

	tests := []struct {
		name string
		want []Renderer
	}{
		{
			name: "list available render engines",
			want: []Renderer{
				&testRenderer{},
				&testRenderer{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := listRegistry.List(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegister(t *testing.T) {
	testRegistry := NewRegistry()
	type args struct {
		name     string
		renderer Renderer
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "registering a new renderer",
			args: args{
				name:     "test-renderer",
				renderer: &testRenderer{},
			},
			wantErr: false,
		},
		{
			name: "registering two renderers with same name should error",
			args: args{
				name:     "test-renderer",
				renderer: &testRenderer{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testRegistry.Register(tt.args.name, tt.args.renderer); (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
