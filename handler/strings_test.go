package handler

import "testing"

func TestFilExt(t *testing.T) {
	tests := []struct {
		file, ext string
	}{
		{"my.file.tar.gz", ".tar.gz"},
		{"my.file.tar.bz2", ".tar.bz2"},
		{"my.file.tar.bz", ".tar.bz"},
		{"my.file.bz2", ".bz2"},
		{"my.file.gz", ".gz"},
		{"my.file.tar.zip", ".tar.zip"}, // :(
	}
	for _, tc := range tests {
		ext := getFileExt(tc.file)
		if ext != tc.ext {
			t.Fatalf("getFileExt(%s) = %s, want %s", tc.file, ext, tc.ext)
		}
	}
}

func TestArch(t *testing.T) {
	tests := []struct {
		file, arch string
	}{
		{"test-armv8-2.11.5.gz", "arm64"},
	}
	for _, tc := range tests {
		ext := getArch(tc.file)
		if ext != tc.arch {
			t.Fatalf("getFileExt(%s) = %s, want %s", tc.file, ext, tc.arch)
		}
	}
}
