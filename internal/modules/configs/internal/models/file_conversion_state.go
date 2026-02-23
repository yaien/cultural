package models

import (
	"errors"
	"slices"
)

type ConversionPreset struct {
	Convert     ConvertFunc
	Variants    []int
	ContentType string
}

type ConversionState struct {
	BiggestFormat           Format
	BiggestFormatIndex      int
	BiggestFormatIsDropable bool
	MissingVariants         []int
	DropableVariantIndexes  []int
	Preset                  *ConversionPreset
}

var ConversionPresets = map[string]*ConversionPreset{
	"image": {
		Convert:     ConvertImage,
		Variants:    []int{320, 640, 1280},
		ContentType: "image/webp",
	},
	"video": {
		Convert:     ConvertVideo,
		Variants:    []int{420, 720},
		ContentType: "video/mp4",
	},
}

var ErrFileHasNoFormats = errors.New("file has no formats")
var ErrFileHasNotConversionPreset = errors.New("file has no conversion preset")

func (f *File) ConversionState() (state ConversionState, err error) {
	if len(f.Formats) == 0 {
		return state, ErrFileHasNoFormats
	}

	preset, ok := ConversionPresets[f.Preset]
	if !ok {
		return state, ErrFileHasNotConversionPreset
	}

	state.Preset = preset

	already := make(map[int]bool)

	for index, format := range f.Formats {
		if state.BiggestFormat.Variant < format.Variant {
			state.BiggestFormat = format
			state.BiggestFormatIndex = index
			state.BiggestFormatIsDropable = format.Variant > slices.Max(preset.Variants) || format.ContentType != preset.ContentType
		}

		if format.ContentType == preset.ContentType {
			already[format.Variant] = true
		}
	}

	for _, variant := range preset.Variants {
		if already[variant] || variant > state.BiggestFormat.Variant {
			continue
		}

		state.MissingVariants = append(state.MissingVariants, variant)

	}
	return state, nil
}
