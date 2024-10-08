package utils

import (
	"image"
	"os"

	"github.com/nfnt/resize"
)

func resizeImage(filePath string, width, height uint) (image.Image, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    img, _, err := image.Decode(file)
    if err != nil {
        return nil, err
    }

    newImg := resize.Resize(width, height, img, resize.Lanczos3)
    return newImg, nil
}
