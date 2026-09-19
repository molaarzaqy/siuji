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

	ResourceTypeImage = "image"
	ResourceTypeVideo = "video"
	ResourceTypeRaw = "raw"
)

// UploadResult carries both the public URL and the Cloudinary public ID.
// The public ID is required later to delete the asset.
type UploadResult struct {
	URL          string
	PublicID     string
	ResourceType string
}

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

func (s *Service) upload(ctx context.Context, file io.Reader, folder, resourceType string) (*UploadResult, error) {
	result, err := s.cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder:       folder,
		ResourceType: resourceType,
	})
	if err != nil {
		return nil, err
	}
	return &UploadResult{
		URL:          result.SecureURL,
		PublicID:     result.PublicID,
		ResourceType: resourceType,
	}, nil
}

func (s *Service) UploadCertificateTemplate(ctx context.Context, file io.Reader) (*UploadResult, error) {
	return s.upload(ctx, file, FolderCertificateTemplates, ResourceTypeImage)
}

func (s *Service) UploadGeneratedCertificate(ctx context.Context, file io.Reader) (*UploadResult, error) {
	return s.upload(ctx, file, FolderCertificatesGenerated, ResourceTypeRaw)
}

func (s *Service) UploadQuestionAudio(ctx context.Context, file io.Reader) (*UploadResult, error) {
	return s.upload(ctx, file, FolderQuestionAudio, ResourceTypeVideo)
}

func (s *Service) UploadQuestionImage(ctx context.Context, file io.Reader) (*UploadResult, error) {
	return s.upload(ctx, file, FolderQuestionImages, ResourceTypeImage)
}

// Destroy removes a single asset. Safe to call with an empty publicID (no-op),
// so callers don't need to null-check every optional file field.
func (s *Service) Destroy(ctx context.Context, publicID, resourceType string) error {
	if publicID == "" {
		return nil
	}
	_, err := s.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID:     publicID,
		ResourceType: resourceType,
	})
	return err
}