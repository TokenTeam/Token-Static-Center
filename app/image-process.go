// Token-Static-Center
// 图片处理模块
// 负责图片的格式转换，长宽修改，水印添加
// LiuFuXin @ Token Team 2018 <loli@lurenjia.in>

package app

import (
	"bytes"
	"errors"
	"gopkg.in/gographics/imagick.v2/imagick"
	"strconv"
)

// 图片缩放
func ImageResize(inputFileStream []byte, width int, keepRatio bool) (outputFileStream []byte, err error) {
	imageHandler := imagick.NewMagickWand()
	// 回收内存
	defer imageHandler.Destroy()

	err = imageHandler.ReadImageBlob(inputFileStream)

	if err != nil {
		return nil, errors.New("图片缩放过程中读取文件时遭遇致命错误：" + err.Error())
	}

	imageWidth := imageHandler.GetImageWidth()
	imageHeight := imageHandler.GetImageHeight()

	imageRatio := float32(imageHeight) / float32(imageWidth)

	switch keepRatio {
	case true:
		// 按比例缩放
		newImageWidth := uint(width)
		newImageHeight := uint(float32(width) * imageRatio)

		// 采取多相位图像插值算法（速度最快，最适合网页图片渲染）
		err = imageHandler.ResizeImage(newImageWidth, newImageHeight, imagick.FILTER_LANCZOS, 1)

		break
	case false:
		// 不按比例缩放
		newImageWidth := uint(width)
		newImageHeight := imageHeight

		// 采取多相位图像插值算法（速度最快，最适合网页图片渲染）
		err = imageHandler.ResizeImage(newImageWidth, newImageHeight, imagick.FILTER_LANCZOS, 1)

		break
	}

	if err != nil {
		return nil, errors.New("图片缩放过程中缩放图片时遭遇致命错误：" + err.Error())
	}

	// 输出文件
	outputFileStream = imageHandler.GetImageBlob()

	return outputFileStream, nil
}

// 图片水印
// 1. watermarkPosition参数
// - 0: 左上角
// - 1: 右上角
// - 2: 左下角
// - 3: 右下角
// - 4: 正中央
// 2. watermarkOpacity参数
// - 0~100，数值越高，越透明
// 3. watermarkSize参数
// - 0~100，数值代表宽度占比(百分比)
func ImageWatermark(inputFileStream []byte, inputWatermarkStream []byte,
	watermarkPosition uint, watermarkOpacity uint,
	watermarkSize uint) (outputFileStream []byte, err error) {
	imageHandler := imagick.NewMagickWand()
	// 回收内存
	defer imageHandler.Destroy()

	watermarkHandler := imagick.NewMagickWand()
	defer watermarkHandler.Destroy()

	// 读取图片资源
	err = imageHandler.ReadImageBlob(inputFileStream)

	if err != nil {
		return nil, errors.New("图片水印处理过程中读取文件时遭遇致命错误：" + err.Error())
	}

	err = watermarkHandler.ReadImageBlob(inputWatermarkStream)

	if err != nil {
		return nil, errors.New("图片水印处理过程中读取文件时遭遇致命错误：" + err.Error())
	}

	// 调整水印透明度
	// 对于包含Alpha通道的图片资源，不能使用SetImageOpacity，否则会失去Alpha通道
	// 需要使用CompositeImage方法进行透明度设置
	// EvaluateImage方法中1.0为不透明，0.0为完全透明
	// 与本方法0~100，100为完全透明的设定不一样，此处顺便进行转换
	imageOpacity := 1.0 - (float64(watermarkOpacity) / float64(100))
	watermarkHandler.EvaluateImage(imagick.EVAL_OP_MULTIPLY, imageOpacity)

	// 获取图片尺寸
	imageWidth := imageHandler.GetImageWidth()
	imageHeight := imageHandler.GetImageHeight()

	// 获取水印尺寸（原图）
	watermarkWidth := watermarkHandler.GetImageWidth()
	watermarkHeight := watermarkHandler.GetImageHeight()

	// 调整水印尺寸
	// 按比例进行缩放
	watermarkRatio := float32(watermarkHeight) / float32(watermarkWidth)
	newWatermarkWidth := uint(float32(imageWidth) * float32(watermarkSize) / 100)
	newWatermarkHeight := uint(float32(newWatermarkWidth) * watermarkRatio)

	err = watermarkHandler.ResizeImage(newWatermarkWidth, newWatermarkHeight, imagick.FILTER_LANCZOS, 1)

	if err != nil {
		return nil, errors.New("图片水印处理过程中缩放水印时遭遇致命错误：" + err.Error())
	}

	// 计算水印位置
	watermarkCoordinateX, watermarkCoordinateY, err := getWatermarkPosition(imageWidth, imageHeight, watermarkPosition, newWatermarkWidth, newWatermarkHeight)

	if err != nil {
		return nil, errors.New("图片水印处理过程中计算水印位置时遭遇致命错误：" + err.Error())
	}

	err = imageHandler.CompositeImage(watermarkHandler, imagick.COMPOSITE_OP_OVER, int(watermarkCoordinateX), int(watermarkCoordinateY))

	if err != nil {
		return nil, errors.New("图片水印处理过程中合成图像时遭遇致命错误：" + err.Error())
	}

	outputFileStream = imageHandler.GetImageBlob()

	return outputFileStream, nil
}

// 文字水印，参数含义同上
func TextWatermark(inputFileStream []byte, watermarkPosition uint,
	watermarkOpacity uint, watermarkSize uint, watermarkColor string,
	watermarkText string, watermarkStyle string) (outputFileStream []byte, err error) {
	imageHandler := imagick.NewMagickWand()
	// 回收内存
	defer imageHandler.Destroy()

	textWatermarkHandler := imagick.NewDrawingWand()
	defer textWatermarkHandler.Destroy()

	imageWatermarkHandler := imagick.NewMagickWand()
	defer imageWatermarkHandler.Destroy()

	colorPalette := imagick.NewPixelWand()
	defer colorPalette.Destroy()

	// 获取静态资源存储目录
	fontPath, err := getStorageRoot()

	if err != nil {
		return nil, errors.New("文字水印处理过程中获取静态资源目录时遭遇致命错误：" + err.Error())
	}

	// 读取图片
	imageHandler.ReadImageBlob(inputFileStream)

	// 获取图片长宽
	imageWidth := imageHandler.GetImageWidth()
	imageHeight := imageHandler.GetImageHeight()

	// 根据文字样式定义字体
	fontName := ""
	switch watermarkStyle {
	// 常规
	case "regular":
		fontName = "msyh.ttf"
		break
		// 细字体
	case "light":
		fontName = "msyh-light.ttf"
		break
		// 粗体
	case "bold":
		fontName = "msyh-bold.ttf"
		break
		// 错误捕获
	default:
		return nil, errors.New("字体类型错误，期望值为regular/light/bold，传入值为" + watermarkStyle)
		break
	}

	// 定义字体路径
	// 此处使用微软雅黑，以便支持多种语言（微博也用微软雅黑，效果非常好）
	// 此前尝试过思源黑体，效果不佳，因为这个字体是为了高PPI显示屏设计的，在图片水印处理上的表现不尽人意
	fontPath = fontPath + "font/" + fontName

	err = textWatermarkHandler.SetFont(fontPath)
	if err != nil {
		return nil, errors.New("文字水印处理过程中调取字体时候遭遇致命错误：" + err.Error())
	}

	// 调色板配色
	colorPalette.SetColor("#" + watermarkColor)
	// 文字填充颜色
	textWatermarkHandler.SetFillColor(colorPalette)
	textWatermarkHandler.SetStrokeWidth(0)
	textWatermarkHandler.SetStrokeAntialias(true)
	textWatermarkHandler.SetStrokeOpacity(0.3)
	textWatermarkHandler.SetStrokeColor(colorPalette)
	// 文字透明度
	imageOpacity := 1.0 - (float64(watermarkOpacity) / float64(100))
	textWatermarkHandler.SetOpacity(imageOpacity)
	// 文字间距
	textWatermarkHandler.SetTextKerning(1)
	// 文字编码
	textWatermarkHandler.SetTextEncoding("UTF-8")

	// 文字对齐方式
	// 根据水印位置确认
	switch watermarkPosition {
	// 左上角
	case 1:
		textWatermarkHandler.SetGravity(imagick.GRAVITY_WEST)
		break
	// 右上角
	case 2:
		textWatermarkHandler.SetGravity(imagick.GRAVITY_EAST)
		break
	// 左下角
	case 3:
		textWatermarkHandler.SetGravity(imagick.GRAVITY_WEST)
		break
	// 右下角
	case 4:
		textWatermarkHandler.SetGravity(imagick.GRAVITY_EAST)
		break
	// 正中央
	case 5:
		textWatermarkHandler.SetGravity(imagick.GRAVITY_CENTER)
		break
	// 错误捕获
	default:
		return nil, errors.New("水印位置错误，期望值为1/2/3/4/5，传入值为" + strconv.Itoa(int(watermarkPosition)))
		break
	}

	// 文字大小
	textWatermarkHandler.SetFontSize(float64(watermarkSize))
	textWatermarkHandler.SetTextAntialias(true)
	// 文字文本
	textWatermarkHandler.Annotation(0, 0, watermarkText)

	// 文字背景为透明
	colorPalette.SetColor("none")

	// 新建水印图片
	// 宽度与长度均由字体大小进行动态设置
	// 为了保持字体清晰锐利，显现出浮雕效果，采取放大4倍后进行超采样的手段
	imageWatermarkHandler.NewImage(uint(float64(bytes.Count([]byte(watermarkText), nil))*float64(watermarkSize)*4), uint(float64(watermarkSize) * 6), colorPalette)
	imageWatermarkZoomWidth := imageWatermarkHandler.GetImageWidth()
	imageWatermarkZoomHeight := imageWatermarkHandler.GetImageHeight()
	imageWatermarkHandler.ResizeImage(imageWatermarkZoomWidth/4, imageWatermarkZoomHeight/4, imagick.FILTER_LANCZOS, 1)

	// 绘制水印文字
	err = imageWatermarkHandler.DrawImage(textWatermarkHandler)

	if err != nil {
		return nil, errors.New("文字水印处理过程中执行文字渲染时遭遇致命错误：" + err.Error())
	}

	// 绘制文字阴影
	shadowLayer := imageWatermarkHandler.Clone()
	shadowLayer.ShadowImage(30, 1.2, 0, 0)
	// 默认的Shadow是白色的，进行反色处理，变成黑色
	shadowLayer.NegateImage(true)
	err = imageWatermarkHandler.CompositeImage(shadowLayer, imagick.COMPOSITE_OP_DST_OVER, -2, 0)

	if err != nil {
		return nil, errors.New("文字水印处理过程中绘制阴影时遭遇致命错误：" + err.Error())
	}

	// 获取文字宽度
	watermarkWidth := imageWatermarkHandler.GetImageWidth()
	watermarkHeight := imageWatermarkHandler.GetImageHeight()

	// 计算位置
	watermarkCoordinateX, watermarkCoordinateY, err := getWatermarkPosition(imageWidth, imageHeight, watermarkPosition, watermarkWidth, watermarkHeight)

	if err != nil {
		return nil, errors.New("文字水印处理过程中计算水印位置时遭遇致命错误：" + err.Error())
	}

	// 图像合成
	err = imageHandler.CompositeImage(imageWatermarkHandler, imagick.COMPOSITE_OP_OVER, int(watermarkCoordinateX), int(watermarkCoordinateY))

	if err != nil {
		return nil, errors.New("文字水印处理过程中合成图像时遭遇致命错误：" + err.Error())
	}

	outputFileStream = imageHandler.GetImageBlob()

	return outputFileStream, nil
}

// 图片格式修改
func ImageReformat(fileStream []byte, targetFormat string) (outputFileStream []byte, err error) {
	imageHandler := imagick.NewMagickWand()
	// 回收内存
	defer imageHandler.Destroy()

	// 读取图片文件
	err = imageHandler.ReadImageBlob(fileStream)
	if err != nil {
		return nil, errors.New("图片格式修改过程中读取图像时遭遇致命错误：" + err.Error())
	}

	// 设置目标格式
	err = imageHandler.SetImageFormat(targetFormat)
	if err != nil {
		return nil, errors.New("图片格式修改过程中转换格式时遭遇致命错误：" + err.Error())
	}

	outputFileStream = imageHandler.GetImageBlob()

	return outputFileStream, nil
}

// 图片压缩
// imageQuality 图像质量
// 从1~100，100为最佳，体积最大，反之亦然
func ImageCompress(fileStream []byte, imageQuality uint) (outputFileStream []byte, err error) {
	imageHandler := imagick.NewMagickWand()
	// 回收内存
	defer imageHandler.Destroy()

	// 判断图片质量有效性
	if imageQuality > 100 || imageQuality < 0 {
		return nil, errors.New("图片压缩过程中判断数值有效性时遭遇致命错误：压缩比率无效，期望值为0~100整数，传入值为" + strconv.Itoa(int(imageQuality)))
	}

	// 防止出现过低质量的图片（糊到没法看）
	if imageQuality < 15 {
		imageQuality = 15
	}

	imageHandler.ReadImageBlob(fileStream)

	// 转换任意图片格式到JPEG
	imageHandler.SetImageFormat("JPEG")

	// 配置压缩算法
	imageHandler.SetCompression(imagick.COMPRESSION_JPEG)

	// 执行压缩
	err = imageHandler.SetImageCompressionQuality(imageQuality)

	if err != nil {
		return nil, errors.New("图片压缩过程中判断数值有效性时遭遇致命错误：压缩出现错误：" + err.Error())
	}

	imageHandler.StripImage()
	outputFileStream = imageHandler.GetImageBlob()

	return outputFileStream, nil
}

// 水印位置计算
func getWatermarkPosition(imageWidth uint, imageHeight uint, watermarkPosition uint,
	watermarkWidth uint, watermarkHeight uint) (watermarkCoordinateX uint,
	watermarkCoordinateY uint, err error) {
	// 边界尺寸
	bordermarginWidth := uint(float32(imageWidth) * 0.025)
	bordermarginHeight := uint(float32(imageHeight) * 0.015)

	// 水印定位
	watermarkCoordinateX = uint(0)
	watermarkCoordinateY = uint(0)

	switch watermarkPosition {
	// 左上角
	case 1:
		// 直接赋值
		watermarkCoordinateX = bordermarginWidth
		watermarkCoordinateY = bordermarginHeight
		break
		// 右上角
	case 2:
		// 需要对X轴进行重计算
		watermarkCoordinateX = imageWidth - bordermarginWidth - watermarkWidth
		watermarkCoordinateY = bordermarginHeight
		break
		// 左下角
	case 3:
		// 需要对Y轴进行重计算
		watermarkCoordinateX = bordermarginWidth
		watermarkCoordinateY = imageHeight - bordermarginHeight - watermarkHeight
		break
		// 右下角
	case 4:
		// XY轴都需要重计算
		watermarkCoordinateX = imageWidth - bordermarginWidth - watermarkWidth
		watermarkCoordinateY = imageHeight - bordermarginHeight - watermarkHeight
		break
		// 正中央
	case 5:
		// 计算方法与上面不太一致
		// 不需要borderMargin参数
		watermarkCoordinateX = (imageWidth - watermarkWidth) / 2
		watermarkCoordinateY = (imageHeight - watermarkHeight) / 2
		break
	default:
		return 0, 0, errors.New("水印位置错误，期望值为1/2/3/4/5，传入值为" + strconv.Itoa(int(watermarkPosition)))
	}
	return watermarkCoordinateX, watermarkCoordinateY, nil
}
