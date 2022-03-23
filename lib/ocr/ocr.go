package ocr

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"os"
	"path"
	"strconv"

	"github.com/otiai10/gosseract"
	"gocv.io/x/gocv"
)

func dist(origin image.Point, contour []image.Point) int {
	sX, sY := 0, 0
	for _, p := range contour {
		sX += p.X - origin.X
		sY += p.Y - origin.Y
	}
	sX /= len(contour)
	sY /= len(contour)
	return sX*sX + sY*sY
}

func max(nums ...int) int {
	ret := nums[0]
	for _, v := range nums {
		if v > ret {
			ret = v
		}
	}
	return ret
}

type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

// RGB画像を受け取ってOCR
func ImageToSudoku(img image.Image, table [][]int, savepath string) error {
	// 画像処理の途中経過を保存するディレクトリ
	if _, err := os.Stat(savepath); os.IsNotExist(err) {
		err = os.MkdirAll(savepath, os.FileMode(0777))
		if err != nil {
			return err
		}
	}

	rgbMat, err := gocv.ImageToMatRGB(img)
	if err != nil {
		return err
	}

	grayMat := RGBToGray(rgbMat)
	blurredMat := blur(grayMat)
	go gocv.IMWrite(path.Join(savepath, "gray.png"), blurredMat)
	binaryMat := binarization(blurredMat)
	go gocv.IMWrite(path.Join(savepath, "binary.png"), binaryMat)

	vertex, err := findVertex(binaryMat)
	if err != nil {
		return err
	}
	go func() {
		contours := gocv.NewPointsVectorFromPoints([][]image.Point{vertex.ToPoints()})
		gocv.DrawContours(&rgbMat, contours, -1, color.RGBA{255, 0, 0, 0}, 2)
		gocv.IMWrite(path.Join(savepath, "contour.png"), rgbMat)
	}()
	// 射影変換
	tableMat := transform(binaryMat, *vertex)
	go gocv.IMWrite(path.Join(savepath, "table.png"), tableMat)

	// SubImageを使う前準備
	tableImg, err := tableMat.ToImage()
	if err != nil {
		return err
	}
	// OCR
	return tableOCR(tableImg, table)
}

// グレースケール
func RGBToGray(rgbMat gocv.Mat) gocv.Mat {
	grayMat := gocv.NewMatWithSize(rgbMat.Rows(), rgbMat.Cols(), gocv.MatTypeCV8UC1)
	gocv.CvtColor(rgbMat, &grayMat, gocv.ColorRGBToGray)
	return grayMat
}

// ガウシアンフィルタ
func blur(grayMat gocv.Mat) gocv.Mat {
	blurredMat := gocv.NewMatWithSize(grayMat.Rows(), grayMat.Cols(), gocv.MatTypeCV8UC1)
	gocv.GaussianBlur(grayMat, &blurredMat, image.Point{3, 3}, 0, 0, gocv.BorderReflect101)
	return blurredMat
}

// 二値化
func binarization(grayMat gocv.Mat) gocv.Mat {
	binaryMat := gocv.NewMatWithSize(grayMat.Rows(), grayMat.Cols(), gocv.MatTypeCV8UC1)
	// AdaptiveThreshold(src Mat, dst *Mat, maxValue float32, adaptiveTyp AdaptiveThresholdType, typ ThresholdType, blockSize int, c float32)
	gocv.AdaptiveThreshold(grayMat, &binaryMat, 255, gocv.AdaptiveThresholdGaussian, gocv.ThresholdBinaryInv, 9, 5)
	return binaryMat
}

// 4つの頂点を見つける
func findVertex(binaryMat gocv.Mat) (*gocv.PointVector, error) {
	center := image.Point{binaryMat.Cols() / 2, binaryMat.Rows() / 2}
	// 輪郭検出
	contours := gocv.FindContours(binaryMat, gocv.RetrievalExternal, gocv.ChainApproxSimple)
	// "頂点の平均位置"が画像の中心に最も近い輪郭を求める
	contour := contours.At(0)
	minDist := dist(center, contour.ToPoints())
	for i := 1; i < contours.Size(); i++ {
		d := dist(center, contours.At(i).ToPoints())
		if d < minDist {
			contour = contours.At(i)
			minDist = d
		}
	}
	// 4つの頂点に縮約
	arcLen := gocv.ArcLength(contour, true)
	vertex := gocv.ApproxPolyDP(contour, arcLen*0.01, true)
	if vertex.Size() != 4 {
		return nil, fmt.Errorf("Cannot detect 4 vertices")
	}
	return &vertex, nil
}

// 色反転
func colorInversion(mat gocv.Mat) gocv.Mat {
	invMat := gocv.NewMatWithSize(mat.Rows(), mat.Cols(), mat.Type())
	gocv.BitwiseNot(mat, &invMat)
	return invMat
}

// 射影変換
func transform(binaryMat gocv.Mat, vertex gocv.PointVector) gocv.Mat {
	arcLen := gocv.ArcLength(vertex, true)
	sqSize := max(int(arcLen/4), 300) //400
	fmt.Printf("size of square: %dx%d\n", sqSize, sqSize)
	tfSize := image.Point{sqSize, sqSize}
	corners := gocv.NewPointVectorFromPoints([]image.Point{{0, 0}, {0, tfSize.Y}, {tfSize.X, tfSize.Y}, {tfSize.X, 0}})
	tfMat := gocv.GetPerspectiveTransform(vertex, corners)
	sqMat := gocv.NewMatWithSize(sqSize, sqSize, binaryMat.Type())
	gocv.WarpPerspective(binaryMat, &sqMat, tfMat, tfSize)
	// 白ベースのMatに変換
	return colorInversion(sqMat)
}

func cellOCR(client *gosseract.Client, imgBytes []byte) int {
	client.SetImageFromBytes(imgBytes)
	text, err := client.Text()
	if err != nil {
		return 0
	} else if len(text) != 1 {
		return 0
	}
	num, err := strconv.Atoi(text)
	if err != nil {
		return 0
	}
	return num
}

func tableOCR(tableImg image.Image, table [][]int) error {
	// gosseract
	client := gosseract.NewClient()
	defer client.Close()
	client.SetPageSegMode(gosseract.PSM_SINGLE_CHAR)
	client.SetWhitelist("123456789")

	size := tableImg.Bounds().Dx()
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			rect := image.Rect(j*size/9+2*(j/3+1), i*size/9+2*(i/3+1), (j+1)*size/9-2*(j/3+1), (i+1)*size/9-2*(i/3+1))
			subImg := tableImg.(SubImager).SubImage(rect)
			// subImg to []bytes
			buf := new(bytes.Buffer)
			if err := png.Encode(buf, subImg); err != nil {
				return err
			}
			// OCR
			num := cellOCR(client, buf.Bytes())
			table[i][j] = num
		}
	}
	return nil
}
