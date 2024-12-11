package handlers

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

const (
	HTML               = "html"
	FileProtocol       = "file://"
	EchartsInstanceDom = "div[_echarts_instance_]"
	CanvasJs           = "echarts.getInstanceByDom(document.querySelector('div[_echarts_instance_]'))" +
		".getDataURL({type: '%s', pixelRatio: %d, excludeComponents: ['toolbox']})"
)

type SnapshotConfig struct {
	// Renderer canvas or svg, not used for now
	Renderer string
	// RenderContent the content bytes of charts after rendered
	RenderContent []byte
	// Path the path to save image
	Path string
	// FileName image name
	FileName string
	// Suffix image format, png, jpeg
	Suffix string
	// Quality the generated image quality, aka pixelRatio
	Quality int
	// KeepHtml whether keep the generated html also, default false
	KeepHtml bool
	// HtmlPath where to keep the generated html, default same to image path
	HtmlPath string
	// Timeout  the timeout config
	Timeout time.Duration
}

type SnapshotConfigOption func(config *SnapshotConfig)

func NewSnapshotConfig(content []byte, image string, opts ...SnapshotConfigOption) *SnapshotConfig {
	path, file := filepath.Split(image)
	suffix := filepath.Ext(file)[1:]
	fileName := file[0 : len(file)-len(suffix)-1]

	config := &SnapshotConfig{
		RenderContent: content,
		Path:          path,
		FileName:      fileName,
		Suffix:        suffix,
		Quality:       1,
		KeepHtml:      false,
		Timeout:       0,
	}

	for _, o := range opts {
		o(config)
	}
	return config
}

func MakeChartSnapshot(content []byte, image string) error {
	return MakeSnapshot(NewSnapshotConfig(content, image))
}

func MakeSnapshot(config *SnapshotConfig) error {
	path := config.Path
	fileName := config.FileName
	content := config.RenderContent
	quality := config.Quality
	suffix := config.Suffix
	keepHtml := config.KeepHtml
	htmlPath := config.HtmlPath
	timeout := config.Timeout

	if htmlPath == "" {
		htmlPath = path
	}

	if !filepath.IsAbs(path) {
		path, _ = filepath.Abs(path)
	}

	if !filepath.IsAbs(htmlPath) {
		htmlPath, _ = filepath.Abs(htmlPath)
	}

	// Append the --no-sandbox option to the allocator options
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.NoSandbox,
	)

	// Create a new execution allocator context with the --no-sandbox flag
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// Create a new context with the allocator context
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	if timeout != 0 {
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	htmlFullPath := filepath.Join(htmlPath, fileName+"."+HTML)

	if !keepHtml {
		defer func() {
			err := os.Remove(htmlFullPath)
			if err != nil {
				log.Printf("Failed to delete the file(%s), err: %s\n", htmlFullPath, err)
			}
		}()
	}

	// Write the HTML content to the file
	err := os.WriteFile(htmlFullPath, content, 0o644)
	if err != nil {
		return err
	}

	if quality < 1 {
		quality = 1
	}

	var base64Data string
	executeJS := fmt.Sprintf(CanvasJs, suffix, quality)
	err = chromedp.Run(ctx,
		chromedp.Navigate(fmt.Sprintf("%s%s", FileProtocol, htmlFullPath)),
		chromedp.WaitVisible(EchartsInstanceDom, chromedp.ByQuery),
		chromedp.Evaluate(executeJS, &base64Data),
	)
	if err != nil {
		return err
	}

	// Decode the base64 image data
	imgContent, err := base64.StdEncoding.DecodeString(strings.Split(base64Data, ",")[1])
	if err != nil {
		return err
	}

	// Save the image to the specified path
	imageFullPath := filepath.Join(path, fmt.Sprintf("%s.%s", fileName, suffix))
	if err := os.WriteFile(imageFullPath, imgContent, 0o644); err != nil {
		return err
	}

	log.Printf("Wrote %s.%s success", fileName, suffix)
	return nil
}
