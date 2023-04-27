package util

import (
	extdraw "golang.org/x/image/draw"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"os"
)

func Jpg2PngAndResize(jpegImg image.Image, width, height int) (string, error) {
	var pngImg image.Image
	bounds := jpegImg.Bounds()
	if bounds.Dx() == width && bounds.Dy() == height {
		pngImg = jpegImg
	} else {
		pngImg = resize(jpegImg, width, height)
	}
	// 创建PNG文件
	pngFile, err := ioutil.TempFile(os.TempDir(), "image_variation*.png")
	if err != nil {
		return "", err
	}
	// 保存PNG文件
	err = png.Encode(pngFile, pngImg)
	if err != nil {
		return "", err
	}
	pngFile.Close()
	return pngFile.Name(), nil
}

// 调整图像大小的函数
func resize(img image.Image, width, height int) image.Image {
	// 使用双线性插值算法进行图像调整
	newImg := image.NewRGBA(image.Rect(0, 0, width, height))
	extdraw.CatmullRom.Scale(newImg, newImg.Bounds(), img, img.Bounds(), draw.Over, nil)
	return newImg
}
