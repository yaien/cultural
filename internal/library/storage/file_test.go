package storage

import (
	"testing"
)

func TestFileGetFormat(t *testing.T) {
	tests := []struct {
		name     string
		variants []int
		query    int
		want     int
		err      string
	}{
		{
			name: "If No Formats: return err",
			err:  "file has no formats",
		},
		{
			name:     "If only one format, return it",
			variants: []int{100},
			query:    300,
			want:     100,
		},
		{
			name:     "If query is 0, return the best format",
			variants: []int{200, 300, 100},
			query:    0,
			want:     300,
		},
		{
			name:     "If query is less than zero formats, return the best format",
			variants: []int{200, 300, 100},
			query:    -1,
			want:     300,
		},
		{
			name:     "If query is greater than all formats, return the best format",
			variants: []int{200, 300, 100},
			query:    400,
			want:     300,
		},
		{
			name:     "If query is between formats, return the best format",
			variants: []int{200, 300, 100},
			query:    250,
			want:     300,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			formats := make([]Format, len(test.variants))
			for i, f := range test.variants {
				formats[i] = Format{Variant: f}
			}

			file := &File{Formats: formats}
			format, err := file.GetFormat(test.query)

			if test.err != "" {
				if err == nil {
					t.Errorf("Expected error '%s', got nil", test.err)
					return
				}

				if err.Error() != test.err {
					t.Errorf("Expected error '%s', got '%s'", test.err, err.Error())
					return
				}

				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %s", err.Error())
			}

			if format.Variant != test.want {
				t.Errorf("Expected format quality %d, got %d", test.want, format.Variant)
			}

		})
	}

}
