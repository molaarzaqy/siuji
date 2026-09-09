package http

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v3"
)

// detectAndValidateContentType sniffs the actual file bytes (not the
// client-supplied header, which is unreliable across HTTP clients like
// Postman/curl/browsers) and returns a fresh io.Reader with the content intact.
func detectAndValidateContentType(file io.Reader, allowTypes []string, fileName string) (io.Reader, error) {
	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to read "+fileName+" file")
	}

	contentType := http.DetectContentType(buf[:n])

	valid := false
	for _, t := range allowTypes {
		if contentType == t {
			valid = true
			break
		}
	}
	if !valid {
		return nil, fiber.NewError(fiber.StatusBadRequest, fileName+" has unsupported file type: "+contentType)
	}
	// Gabungkan kembali bytes yang sudah "dicicipi" dengan sisa file,
	// supaya tidak ada data yang hilang saat nanti diupload.
	return io.MultiReader(bytes.NewReader(buf[:n]), file), nil
}

func extractOptionalFile(c fiber.Ctx, field string, allowedTypes []string) (io.Reader, func(), error) {
	noop := func() {}

	fh, err := c.FormFile(field)
	if err != nil {
		return nil, noop, nil // file tidak dikirim, bukan error
	}

	file, err := fh.Open()
	if err != nil {
		return nil, noop, fiber.NewError(fiber.StatusInternalServerError, "failed to read "+field+" file")
	}

	validated, err := detectAndValidateContentType(file, allowedTypes, field)
	if err != nil {
		file.Close()
		return nil, noop, err
	}

	return validated, func() { file.Close() }, nil
}