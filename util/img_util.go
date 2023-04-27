package util

import (
	extdraw "golang.org/x/image/draw"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"os"
)

func Jpg2PngAndResize(img image.Image, sideLen int) (string, error) {
	pngImg := img
	if img.Bounds().Dy() != img.Bounds().Dx() {
		pngImg = toSquare(img)
	}
	if pngImg.Bounds().Dy() != sideLen {
		pngImg = resize(pngImg, sideLen, sideLen)
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
	//强制写入
	pngFile.Close()
	return pngFile.Name(), nil
}

func toSquare(img image.Image) *image.RGBA {
	width := img.Bounds().Size().X
	height := img.Bounds().Size().Y

	// 计算新图片的大小
	size := width
	if height > width {
		size = height
	}

	// 创建新的图片对象
	newImg := image.NewRGBA(image.Rect(0, 0, size, size))

	// 使用白色填充新图片的背景
	white := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	draw.Draw(newImg, newImg.Bounds(), &image.Uniform{C: white}, image.Point{}, draw.Src)

	// 计算要绘制的图片的位置
	x := (size - width) / 2
	y := (size - height) / 2
	pt := image.Pt(x, y)

	// 将原图片绘制到新图片上
	draw.Draw(newImg, img.Bounds().Add(pt), img, image.Point{}, draw.Src)
	return newImg
}

// 调整图像大小的函数
func resize(img image.Image, width, height int) image.Image {
	// 使用双线性插值算法进行图像调整
	newImg := image.NewRGBA(image.Rect(0, 0, width, height))
	extdraw.CatmullRom.Scale(newImg, newImg.Bounds(), img, img.Bounds(), draw.Over, nil)
	return newImg
}
