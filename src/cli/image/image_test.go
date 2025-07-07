package image

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetOutputPath(t *testing.T) {
	cases := []struct {
		Case     string
		Config   string
		Path     string
		Expected string
	}{
		{Case: "default config", Expected: "prompt.png"},
		{Case: "hidden file", Config: ".posh.omp.json", Expected: "posh.png"},
		{Case: "hidden file toml", Config: ".posh.omp.toml", Expected: "posh.png"},
		{Case: "hidden file yaml", Config: ".posh.omp.yaml", Expected: "posh.png"},
		{Case: "hidden file yml", Config: ".posh.omp.yml", Expected: "posh.png"},
		{Case: "path provided", Path: "mytheme.png", Expected: "mytheme.png"},
		{Case: "relative, no omp", Config: "~/jandedobbeleer.json", Expected: "jandedobbeleer.png"},
		{Case: "relative path", Config: "~/jandedobbeleer.omp.json", Expected: "jandedobbeleer.png"},
		{Case: "invalid config name", Config: "~/jandedobbeleer.omp.foo", Expected: "prompt.png"},
	}

	for _, tc := range cases {
		image := &Renderer{
			Path: tc.Path,
		}

		image.setOutputPath(tc.Config)

		assert.Equal(t, tc.Expected, image.Path, tc.Case)
	}
}

func TestHexToRGB(t *testing.T) {
	cases := []struct {
		name     string
		hex      HexColor
		expected *RGB
		hasError bool
	}{
		{
			name:     "Valid hex with hash",
			hex:      "#FF0000",
			expected: &RGB{255, 0, 0},
			hasError: false,
		},
		{
			name:     "Valid hex without hash",
			hex:      "00FF00",
			expected: &RGB{0, 255, 0},
			hasError: false,
		},
		{
			name:     "Valid hex blue",
			hex:      "#0000FF",
			expected: &RGB{0, 0, 255},
			hasError: false,
		},
		{
			name:     "Invalid hex too short",
			hex:      "#FFF",
			expected: nil,
			hasError: true,
		},
		{
			name:     "Invalid hex too long",
			hex:      "#FFFFFFF",
			expected: nil,
			hasError: true,
		},
		{
			name:     "Invalid hex characters",
			hex:      "#GGGGGG",
			expected: nil,
			hasError: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := tc.hex.RGB()

			if tc.hasError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestGetColorNameFromCode(t *testing.T) {
	cases := []struct {
		expected  string
		colorCode int
	}{
		{"black", 30},
		{"red", 31},
		{"green", 32},
		{"yellow", 33},
		{"blue", 34},
		{"magenta", 35},
		{"cyan", 36},
		{"white", 37},
		{"black", 40}, // background
		{"red", 41},   // background
		{"darkGray", 90},
		{"lightRed", 91},
		{"lightGreen", 92},
		{"lightYellow", 93},
		{"lightBlue", 94},
		{"lightMagenta", 95},
		{"lightCyan", 96},
		{"lightWhite", 97},
		{"", 999}, // invalid code
	}

	for _, tc := range cases {
		t.Run(tc.expected, func(t *testing.T) {
			result := colorNameFromCode(tc.colorCode)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestSetBase16ColorWithOverrides(t *testing.T) {
	renderer := &Renderer{
		defaultForegroundColor: &RGB{255, 255, 255},
		Settings: Settings{
			Colors: map[string]HexColor{
				"red":  "#FF6B6B",
				"blue": "#4ECDC4",
			},
		},
	}

	// Test with color override
	renderer.setBase16Color("31") // Red foreground
	assert.Equal(t, &RGB{255, 107, 107}, renderer.foregroundColor)

	// Test with color override for background
	renderer.setBase16Color("44") // Blue background
	assert.Equal(t, &RGB{78, 205, 196}, renderer.backgroundColor)

	// Test without override (should use default)
	renderer.setBase16Color("32") // Green foreground (no override)
	assert.Equal(t, &RGB{57, 181, 74}, renderer.foregroundColor)
}

func TestSetBase16ColorWithoutOverrides(t *testing.T) {
	renderer := &Renderer{
		defaultForegroundColor: &RGB{255, 255, 255},
		Settings:               Settings{}, // No color overrides
	}

	// Test default red
	renderer.setBase16Color("31")
	assert.Equal(t, &RGB{222, 56, 43}, renderer.foregroundColor)

	// Test default blue background
	renderer.setBase16Color("44")
	assert.Equal(t, &RGB{0, 111, 184}, renderer.backgroundColor)
}

func TestSetBase16ColorInvalidInput(t *testing.T) {
	renderer := &Renderer{
		defaultForegroundColor: &RGB{255, 255, 255},
		Settings:               Settings{},
	}

	// Test invalid color code
	renderer.setBase16Color("invalid")
	assert.Equal(t, &RGB{255, 255, 255}, renderer.foregroundColor)
}
