package services

import (
	"fmt"
	"github.com/otiai10/gosseract/v2"
	"gocv.io/x/gocv"
	"image"
	_ "image/jpeg" // Import necessary image formats (jpeg, png, etc.)
	"io/ioutil"
	"mime/multipart"
	"vehicle-plate-recognition/store/entities"
	"vehicle-plate-recognition/store/postgres"
)

type VehicleService struct {
	db *postgres.Database
}

func NewVehicleService(db *postgres.Database) *VehicleService {
	return &VehicleService{db: db}
}

func (s *VehicleService) ProcessVehicleImage(file *multipart.FileHeader) (string, error) {
	// Read the uploaded file into memory
	fileData, err := s.readUploadedFile(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	// Preprocess the image to improve OCR accuracy
	preprocessedImg, err := s.preprocessImage(fileData)
	if err != nil {
		return "", fmt.Errorf("failed to preprocess image: %v", err)
	}

	// Use Tesseract OCR to extract the plate text
	client := gosseract.NewClient()
	defer client.Close()
	client.SetImageFromBytes(preprocessedImg)
	client.SetWhitelist("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	client.SetPageSegMode(gosseract.PSM_SINGLE_LINE) // Assuming single line for license plates

	plateText, err := client.Text()
	if err != nil {
		return "", fmt.Errorf("failed to extract text: %v", err)
	}

	// Save the extracted plate text to the database
	vehicle := entities.Vehicle{Plate: plateText}
	if err := s.db.DB.Create(&vehicle).Error; err != nil {
		return "", fmt.Errorf("failed to save vehicle record: %v", err)
	}

	return plateText, nil
}

// Preprocess the image to improve OCR accuracy
func (s *VehicleService) preprocessImage(imageData []byte) ([]byte, error) {
	img, err := gocv.IMDecode(imageData, gocv.IMReadColor)
	if err != nil || img.Empty() {
		return nil, fmt.Errorf("failed to decode image")
	}
	defer img.Close()

	// Convert to grayscale
	grayImg := gocv.NewMat()
	defer grayImg.Close()
	gocv.CvtColor(img, &grayImg, gocv.ColorBGRToGray)

	// Apply Gaussian blur to reduce noise
	blurredImg := gocv.NewMat()
	defer blurredImg.Close()
	gocv.GaussianBlur(grayImg, &blurredImg, image.Pt(5, 5), 0, 0, gocv.BorderDefault)

	// Apply binary thresholding
	threshImg := gocv.NewMat()
	defer threshImg.Close()
	gocv.Threshold(blurredImg, &threshImg, 0, 255, gocv.ThresholdBinary|gocv.ThresholdOtsu)

	// Encode the preprocessed image back to bytes
	buffer, err := gocv.IMEncode(".jpg", threshImg)
	if err != nil {
		return nil, fmt.Errorf("failed to encode preprocessed image: %v", err)
	}

	// Convert NativeByteBuffer to byte slice
	return buffer.GetBytes(), nil
}

// Helper function to read the uploaded file into memory
func (s *VehicleService) readUploadedFile(file *multipart.FileHeader) ([]byte, error) {
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open source file: %v", err)
	}
	defer src.Close()

	return ioutil.ReadAll(src)
}
