package cloudinary

import (
	"context"
	"io"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

const (
	FolderCertificateTemplates  = "siuji/certificates/templates"
	FolderCertificatesGenerated = "siuji/certificates/generated"
	FolderQuestionAudio         = "siuji/questions/audio"
	FolderQuestionImages        = "siuji/questions/images"
)

type Service struct {
	cld *cloudinary.Cloudinary
}

func NewService(cloudinaryURL string) (*Service, error) {
	cld, err := cloudinary.NewFromURL(cloudinaryURL)
	if err != nil {
		return nil, err
	}
	return &Service{cld: cld}, nil
}

func (s *Service) upload(ctx context.Context, file io.Reader, folder, resorceType string) (string, error) {
	result, err := s.cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder: folder,
		ResourceType: resorceType,
	})
	if err != nil {
		return "", err
	}
	return result.SecureURL, nil
}

func (s *Service) UploadCertificateTemplate(ctx context.Context, file io.Reader) (string, error) {
	return s.upload(ctx, file, FolderCertificateTemplates, "image")
}

func (s *Service) UploadGeneratedCertificate(ctx context.Context, file io.Reader) (string, error) {
	return s.upload(ctx, file, FolderCertificatesGenerated, "raw")
}

func (s *Service) UploadQuestionAudio(ctx context.Context, file io.Reader) (string, error) {
	return s.upload(ctx, file, FolderQuestionAudio, "video") // Cloudinary menyimpan audio di bawah resource type "video"
}

func (s *Service) UploadQuestionImage(ctx context.Context, file io.Reader) (string, error) {
	return s.upload(ctx, file, FolderQuestionImages, "image")
}