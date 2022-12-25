package transport

import (
	"bytes"
	"crypto/subtle"
	"reflect"
	"testing"
)

func TestSerializeSlice(t *testing.T) {
	tests := []struct {
		name  string
		slice []byte
		want  []byte
	}{
		{
			name:  "normal case",
			slice: []byte{1, 2, 3},
			want:  []byte{3, 1, 2, 3},
		},
		{
			name:  "empty slice",
			slice: []byte{},
			want:  []byte{0},
		},
		{
			name:  "length overflow",
			slice: make([]byte, 256),
			want:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SerializeSlice(tt.slice); subtle.ConstantTimeCompare(got, tt.want) != 1 {
				t.Errorf("SerializeSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeserializeSlice(t *testing.T) {
	tests := []struct {
		name        string
		bytesToRead []byte
		want        []byte
		wantErr     bool
	}{
		{
			name:        "normal case",
			bytesToRead: []byte{3, 1, 2, 3},
			want:        []byte{1, 2, 3},
			wantErr:     false,
		},
		{
			name:        "expected empty slice",
			bytesToRead: []byte{0},
			want:        nil,
			wantErr:     false,
		},
		{
			name:        "empty slice",
			bytesToRead: []byte{},
			wantErr:     true,
		},
		{
			name:        "not enough to read",
			bytesToRead: []byte{3, 1, 2},
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DeserializeSlice(bytes.NewReader(tt.bytesToRead))
			if (err != nil) != tt.wantErr {
				t.Errorf("DeserializeSlice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeserializeSlice() got = %v, want %v", got, tt.want)
			}
		})
	}
}
