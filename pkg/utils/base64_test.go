package utils

import (
	"os"
	"testing"
)

func TestFileToBase64(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	path := pwd + "/" + "cert.crt"
	defer func() {
		err = os.RemoveAll(path)
		if err != nil {
			panic(err)
		}
		err = os.RemoveAll(pwd + "/" + "logs")
		if err != nil {
			panic(err)
		}
	}()
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	_, err = file.Write([]byte(`Hello`))
	if err != nil {
		panic(err)
	}
	file.Sync()
	file.Close()
	tests := []struct {
		name    string
		path    string
		want    []byte
		wantErr bool
	}{
		{
			name:    "TestFileToBase64",
			path:    path,
			want:    []byte("SGVsbG8="),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FileToBase64(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileToBase64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(got) != string(tt.want) {
				t.Errorf("FileToBase64() got = %v, want %v", got, tt.want)
			}
		})
	}
}
