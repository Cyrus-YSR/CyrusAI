package image

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	ort "github.com/yalue/onnxruntime_go"
	"golang.org/x/image/draw"
)

type ImageRecognizer struct {
	session      *ort.Session[float32]
	inputName    string
	outputName   string
	inputH       int
	inputW       int
	labels       []string
	inputTensor  *ort.Tensor[float32]
	outputTensor *ort.Tensor[float32]
}

const (
	defaultInputName  = "data"
	defaultOutputName = "resnetv24_dense0_fwd"
)

var (
	initOnce sync.Once
	initErr  error
)

// NewImageRecognizer 创建识别器（自动使用默认 input/output 名称）
func NewImageRecognizer(modelPath, labelPath string, inputH, inputW int) (*ImageRecognizer, error) {
	if inputH <= 0 || inputW <= 0 {
		inputH, inputW = 224, 224
	}

	// 初始化 ONNX 环境（全局一次）
	initOnce.Do(func() {
		if os.Getenv("ORT_DYLIB_PATH") == "" {
			_ = os.Setenv("ORT_DYLIB_PATH", "/usr/lib/libonnxruntime.so")
		}
		initErr = ort.InitializeEnvironment()
	})
	if initErr != nil {
		return nil, fmt.Errorf("onnxruntime initialize error: %w", initErr)
	}

	inputName := os.Getenv("IMAGE_MODEL_INPUT_NAME")
	if inputName == "" {
		inputName = defaultInputName
	}
	outputName := os.Getenv("IMAGE_MODEL_OUTPUT_NAME")
	if outputName == "" {
		outputName = defaultOutputName
	}

	inputShape := ort.NewShape(1, 3, int64(inputH), int64(inputW))
	inData := make([]float32, inputShape.FlattenedSize())
	inTensor, err := ort.NewTensor(inputShape, inData)
	if err != nil {
		return nil, fmt.Errorf("create input tensor failed: %w", err)
	}

	outShape := ort.NewShape(1, 1000)
	outTensor, err := ort.NewEmptyTensor[float32](outShape)
	if err != nil {
		inTensor.Destroy()
		return nil, fmt.Errorf("create output tensor failed: %w", err)
	}

	session, err := ort.NewSession[float32](
		modelPath,
		[]string{inputName},
		[]string{outputName},
		[]*ort.Tensor[float32]{inTensor},
		[]*ort.Tensor[float32]{outTensor},
	)
	if err != nil {
		inTensor.Destroy()
		outTensor.Destroy()
		return nil, fmt.Errorf("create onnx session failed: %w", err)
	}

	// 读取 label 文件
	labels, err := loadLabels(labelPath)
	if err != nil {
		session.Destroy()
		inTensor.Destroy()
		outTensor.Destroy()
		return nil, err
	}

	return &ImageRecognizer{
		session:      session,
		inputName:    inputName,
		outputName:   outputName,
		inputH:       inputH,
		inputW:       inputW,
		labels:       labels,
		inputTensor:  inTensor,
		outputTensor: outTensor,
	}, nil
}

func (r *ImageRecognizer) Close() {
	if r.session != nil {
		_ = r.session.Destroy()
		r.session = nil
	}
	if r.inputTensor != nil {
		_ = r.inputTensor.Destroy()
		r.inputTensor = nil
	}
	if r.outputTensor != nil {
		_ = r.outputTensor.Destroy()
		r.outputTensor = nil
	}
}

func (r *ImageRecognizer) PredictFromFile(imagePath string) (string, error) {
	file, err := os.Open(filepath.Clean(imagePath))
	if err != nil {
		return "", fmt.Errorf("image not found: %w", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}

	return r.PredictFromImage(img)
}

func (r *ImageRecognizer) PredictFromBuffer(buf []byte) (string, error) {
	img, _, err := image.Decode(bytes.NewReader(buf))
	if err != nil {
		return "", fmt.Errorf("failed to decode image from buffer: %w", err)
	}
	return r.PredictFromImage(img)
}

func (r *ImageRecognizer) PredictFromImage(img image.Image) (string, error) {

	resizedImg := image.NewRGBA(image.Rect(0, 0, r.inputW, r.inputH))

	draw.CatmullRom.Scale(resizedImg, resizedImg.Bounds(), img, img.Bounds(), draw.Over, nil)

	h, w := r.inputH, r.inputW
	ch := 3 // R, G, B
	data := make([]float32, h*w*ch)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			c := resizedImg.At(x, y)

			r, g, b, _ := c.RGBA()

			rf := float32(r>>8) / 255.0
			gf := float32(g>>8) / 255.0
			bf := float32(b>>8) / 255.0

			// NCHW format
			data[y*w+x] = rf
			data[h*w+y*w+x] = gf
			data[2*h*w+y*w+x] = bf
		}
	}

	inData := r.inputTensor.GetData()
	copy(inData, data)

	if err := r.session.Run(); err != nil {
		return "", fmt.Errorf("onnx run error: %w", err)
	}

	outData := r.outputTensor.GetData()
	if len(outData) == 0 {
		return "", errors.New("empty output from model")
	}

	maxIdx := 0
	maxVal := outData[0]
	for i := 1; i < len(outData); i++ {
		if outData[i] > maxVal {
			maxVal = outData[i]
			maxIdx = i
		}
	}

	if maxIdx >= 0 && maxIdx < len(r.labels) {
		return r.labels[maxIdx], nil
	}
	return "Unknown", nil
}

func loadLabels(path string) ([]string, error) {
	resolvedPath, err := ensureLabels(path)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(filepath.Clean(resolvedPath))
	if err != nil {
		return nil, fmt.Errorf("open label file failed: %w", err)
	}
	defer f.Close()

	var labels []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		if line != "" {
			labels = append(labels, line)
		}
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("read labels failed: %w", err)
	}
	if len(labels) == 0 {
		return nil, fmt.Errorf("no labels found in %s", resolvedPath)
	}
	return labels, nil
}

func ensureLabels(path string) (string, error) {
	info, err := os.Stat(path)
	if err == nil {
		if info.IsDir() {
			return ensureLabelsInTemp()
		}
		return path, nil
	}
	if !os.IsNotExist(err) {
		return "", fmt.Errorf("stat label file failed: %w", err)
	}
	if err := downloadLabels(path); err == nil {
		return path, nil
	}
	return ensureLabelsInTemp()
}

func ensureLabelsInTemp() (string, error) {
	tempPath := filepath.Join(os.TempDir(), "imagenet_classes.txt")
	info, err := os.Stat(tempPath)
	if err == nil && !info.IsDir() {
		return tempPath, nil
	}
	if err != nil && !os.IsNotExist(err) {
		return "", fmt.Errorf("stat temp label file failed: %w", err)
	}
	if err := downloadLabels(tempPath); err != nil {
		return "", err
	}
	return tempPath, nil
}

func downloadLabels(path string) error {
	resp, err := http.Get("https://raw.githubusercontent.com/pytorch/hub/master/imagenet_classes.txt")
	if err != nil {
		return fmt.Errorf("download labels failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("download labels failed: status %d", resp.StatusCode)
	}
	f, err := os.Create(filepath.Clean(path))
	if err != nil {
		return fmt.Errorf("create label file failed: %w", err)
	}
	defer f.Close()
	if _, err := io.Copy(f, resp.Body); err != nil {
		return fmt.Errorf("write labels failed: %w", err)
	}
	return nil
}
