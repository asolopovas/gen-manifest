package genWebmanifest

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/davidbyttow/govips/v2/vips"
)

type JsonConfig struct {
	Metadata  Metadata `json:"metadata"`
	Prefix    string   `json:"prefix"`
	IconSizes []int    `json:"sizes"`
}
type Metadata struct {
	Name            string `json:"name"`
	ShortName       string `json:"short_name"`
	Description     string `json:"description"`
	StartUrl        string `json:"start_url"`
	BackgroundColor string `json:"background_color"`
	Display         string `json:"display"`
	Scope           string `json:"scope"`
	ThemeColor      string `json:"theme_color"`
	Icons           []Icon `json:"icons"`
}

type Icon struct {
	Src   string `json:"src"`
	Sizes string `json:"sizes"`
	Type  string `json:"type"`
}

func GetJsonConfig(configPath string) (JsonConfig, error) {
	result := JsonConfig{}

	jsonConfig, err := os.ReadFile(configPath)

	json.Unmarshal([]byte(jsonConfig), &result)
	return result, err
}

func GenConfig(configPath string) {
	defaultConf := []byte(`{
		"sizes": [48, 72, 96, 144, 168, 192, 384, 512],
		"prefix": "",
		"metadata": {
    		"name": "Application Name",
			"short_name": "AN",
    		"description": "",
			"start_url": "/",
			"theme_color": "#000000",
			"background_color": "#ffffff",
			"scope": "/",
			"display": "standalone"
		}
	}
`)

	writeErr := os.WriteFile(configPath, defaultConf, 0644)
	ErrChk(writeErr)

}

func ResizeImage(src string, dest string, size int, iconsDir string) {
	// Save the original logging settings

	// Define a no-op logging handler
	noOpHandler := func(messageDomain string, messageLevel vips.LogLevel, message string) {}

	// Disable logging
	vips.LoggingSettings(noOpHandler, vips.LogLevelError)

	// Read the source image
	img, err := vips.NewImageFromFile(src)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer img.Close()

	// Resize the image
	err = img.Resize(float64(size)/float64(img.Width()), vips.KernelLanczos3)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	if !PathExist(iconsDir) {
		err := os.MkdirAll(iconsDir, os.ModePerm)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
	}

	ep := vips.NewDefaultPNGExportParams()
	imageBytes, _, err := img.Export(ep)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	err = os.WriteFile(dest, imageBytes, 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
}

func GenWebmanifest(conf JsonConfig, icon string, iconsDir string, manifestName string) {
	for _, size := range conf.IconSizes {
		dest := filepath.Join(iconsDir, fmt.Sprintf("icon-%d.png", size))
		ResizeImage(icon, dest, size, iconsDir)
		src := dest
		if conf.Prefix != "" {
			src = fmt.Sprintf("%s/%s", conf.Prefix, dest)
		}
		conf.Metadata.Icons = append(conf.Metadata.Icons, Icon{
			Src:   src,
			Sizes: fmt.Sprintf("%dx%d", size, size),
			Type:  "image/png",
		})
	}
	manifestData, err := json.MarshalIndent(conf.Metadata, "", "  ")
	if err != nil {
		fmt.Println("Error generating manifest:", err)
		return
	}

	err = os.WriteFile(manifestName, manifestData, 0644)
	if err != nil {
		fmt.Println("Error writing manifest file:", err)
	}
}
