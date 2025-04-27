package img_tools

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
)

const (
	maxThumbnailWidth  = 300
	maxThumbnailHeight = 300
)

func CreateThumbnail(originalImage []byte) ([]byte, error) {
	img, format, err := image.Decode(bytes.NewReader(originalImage))
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	ratio := float64(bounds.Dx()) / float64(bounds.Dy())

	var newWidth, newHeight uint
	if ratio > 1 {
		newWidth = maxThumbnailWidth
		newHeight = uint(float64(maxThumbnailWidth) / ratio)
	} else {
		newHeight = maxThumbnailHeight
		newWidth = uint(float64(maxThumbnailHeight) * ratio)
	}

	thumbnail := resize.Resize(newWidth, newHeight, img, resize.Lanczos3)

	// Encode thumbnail to bytes
	var buf bytes.Buffer
	switch format {
	case "jpeg":
		err = jpeg.Encode(&buf, thumbnail, &jpeg.Options{Quality: 75})
	case "png":
		err = png.Encode(&buf, thumbnail)
	}
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func GenerateUUID() string {
	uuid := make([]byte, 16)
	_, _ = rand.Read(uuid)

	// Set version (4) and variant bits (RFC 4122 compliance)
	uuid[6] = (uuid[6] & 0x0f) | (4 << 4)
	uuid[8] = uuid[8]&(0xff>>2) | (0x02 << 6)

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		uuid[0:4],
		uuid[4:6],
		uuid[6:8],
		uuid[8:10],
		uuid[10:16],
	)
}

func CreateThumbnailByPath(path string) ([]byte, error) {
	defer func() {
		if err := os.Remove(path); err != nil {
			log.Printf("[image_tools] ⚠️ Failed to remove file %s: %v", path, err)
		}
		log.Printf("[image_tools] 🗑 Removed file %s", path)
	}()

	photo, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read image file %s: %w", path, err)
	}

	thumbnail, err := CreateThumbnail(photo)
	if err != nil {
		return nil, fmt.Errorf("failed to create thumbnail: %w", err)
	}

	return thumbnail, nil
}
