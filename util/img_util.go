package util

import (
	"bytes"
	extdraw "golang.org/x/image/draw"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"os"
)

const _4M = 4 * 1024 * 1024

func Jpg2PngAndResize(img image.Image, sideLen int) (*os.File, error) {
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
		return nil, err
	}
	if png.Encode(pngFile, pngImg) != nil {
		return nil, err
	}
	// 强制刷盘PNG文件
	if err := pngFile.Sync(); err != nil {
		return nil, err
	}
	//将读写偏移量置于文件起始位置
	if _, err := pngFile.Seek(0, 0); err != nil {
		return nil, err
	}
	fileStat, err := pngFile.Stat()
	if err != nil {
		return nil, err
	}
	if fileStat.Size() < _4M {
		return pngFile, nil
	}
	//压缩会生成新的图片，所以先将原高分辨率的图片删除
	os.Remove(pngFile.Name())
	return compressImg(pngImg)

}

func compressImg(pngImg image.Image) (*os.File, error) {
	//压缩高分辨率的图片
	buffer, err := compressPngImage(pngImg, _4M)
	if err != nil {
		return nil, err

	}
	compressPngFile, err := ioutil.TempFile(os.TempDir(), "compress_for_gpt*.png")
	if err != nil {
		return nil, err

	}
	_, err = compressPngFile.Write(buffer.Bytes())
	if err != nil {
		return nil, err
	}
	// 强制刷盘PNG文件
	if err := compressPngFile.Sync(); err != nil {
		return nil, err
	}
	//将读写偏移量置于文件起始位置
	if _, err := compressPngFile.Seek(0, 0); err != nil {
		return nil, err
	}
	return compressPngFile, nil
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

func compressPngImage(img image.Image, maxSize int) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	err := png.Encode(buf, img)
	if err != nil {
		return nil, err
	}
	for {
		// 如果文件大小小于指定大小，则退出循环
		if buf.Len() < maxSize {
			break
		}

		// 对 PNG 数据进行压缩
		compressedBuf := new(bytes.Buffer)
		err := png.Encode(compressedBuf, img)
		if err != nil {
			return nil, err
		}

		// 如果压缩后的数据大小大于原始数据大小，则退出循环
		if compressedBuf.Len() >= buf.Len() {
			break
		}

		// 将压缩后的数据保存到内存缓冲区中
		buf = compressedBuf
	}
	return buf, nil
}
