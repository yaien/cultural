package models

import (
	"os"
	"testing"
)

func TestGetFileDimension(t *testing.T) {
	var tests = []struct {
		name        string
		file        string
		contentType string
		width       int
		height      int
		quality     int
	}{
		{
			name:        "big_photo",
			file:        "big_photo.jpg",
			contentType: "image/jpeg",
			width:       3303,
			height:      4954,
			quality:     4954,
		},
		{
			name:        "big_video",
			file:        "big_video.mp4",
			contentType: "video/mp4",
			width:       1920,
			height:      1080,
			quality:     1080,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input, err := os.OpenInRoot("testdata", test.file)
			if err != nil {
				t.Fatalf("failed getting input file: %v", err)
			}

			defer input.Close()

			width, height, quality, err := GetFileDimension(input, test.contentType)
			if err != nil {
				t.Fatalf("GetFileDimension returned an error: %v", err)
			}
			if width != test.width {
				t.Errorf("Expected width %d, got %d", test.width, width)
			}
			if height != test.height {
				t.Errorf("Expected height %d, got %d", test.height, height)
			}
			if quality != test.quality {
				t.Errorf("Expected quality %d, got %d", test.quality, quality)
			}
		})
	}
}
