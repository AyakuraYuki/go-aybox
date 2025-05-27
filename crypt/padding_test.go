package crypt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPKCS7UnPadding(t *testing.T) {
	type args struct {
		src []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr error
	}{
		{
			name:    "len(src)==0",
			args:    args{src: make([]byte, 0)},
			want:    make([]byte, 0),
			wantErr: ErrUnPadding,
		},
		{
			name:    `src=="121"`,
			args:    args{src: []byte{1, 2, 1}},
			want:    []byte{1, 2},
			wantErr: nil,
		},
		{
			name:    `src=="12121"`,
			args:    args{src: []byte{1, 2, 1, 1, 9}},
			want:    []byte{1, 2, 1, 1, 9},
			wantErr: ErrUnPadding,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PKCS7UnPadding(tt.args.src)
			// t.Log(string(got))
			assert.Equal(t, tt.wantErr, err)
			assert.EqualValues(t, tt.want, got)
		})
	}
}
