package storage

import "testing"

func TestGetFileDimension(t *testing.T) {
	var tests = []struct {
		name        string
		file        string
		contentType string
		width       int
		height      int
		variant     int
	}{
		{
			name:        "big_photo",
			file:        "testdata/big_photo.jpg",
			contentType: "image/jpeg",
			width:       3303,
			height:      4954,
			variant:     4954,
		},
		{
			name:        "big_video",
			file:        "testdata/big_video.mp4",
			contentType: "video/mp4",
			width:       1920,
			height:      1080,
			variant:     1080,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := t.Context()
			width, height, quality, err := GetDimensionByContentType(ctx, test.file, test.contentType)
			if err != nil {
				t.Fatalf("dimension returned an error: %v", err)
			}
			if width != test.width {
				t.Errorf("Expected width %d, got %d", test.width, width)
			}
			if height != test.height {
				t.Errorf("Expected height %d, got %d", test.height, height)
			}
			if quality != test.variant {
				t.Errorf("Expected quality %d, got %d", test.variant, quality)
			}
		})
	}
}
