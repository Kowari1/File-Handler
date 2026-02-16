package file

import (
	"bufio"
	"context"
	"os"
	"path/filepath"

	model "github.com/Kowari1/File-Handler/internal/domain"
	"github.com/Kowari1/File-Handler/internal/pdf"
	"github.com/Kowari1/File-Handler/internal/repository"
)

type FileService struct {
	processor         *Processor
	processedFileRepo repository.ProcessedFileRepository
	pdfGenerator      pdf.PDFGenerator
	outputDir         string
}

func NewFileService(
	processor *Processor,
	processedFileRepo repository.ProcessedFileRepository,
	pdfGenerator pdf.PDFGenerator,
	outputDir string,
) *FileService {
	return &FileService{
		processor:         processor,
		processedFileRepo: processedFileRepo,
		pdfGenerator:      pdfGenerator,
		outputDir:         outputDir,
	}
}

func (s *FileService) Handle(
	ctx context.Context,
	filePath string,
) error {

	fileName := filepath.Base(filePath)

	exists, err := s.processedFileRepo.Exists(ctx, fileName)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	if err := s.processedFileRepo.Create(ctx, fileName); err != nil {
		return err
	}

	file, err := os.Open(filePath)
	if err != nil {
		s.markFailed(ctx, fileName, err)
		return err
	}
	defer file.Close()

	lines := make(chan string)

	go func() {
		defer close(lines)

		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return
			case lines <- scanner.Text():
			}
		}
	}()

	guidMap, err := s.processor.Process(ctx, fileName, lines)
	if err != nil {
		s.markFailed(ctx, fileName, err)
		return err
	}

	for guid, devices := range guidMap {
		err := s.pdfGenerator.Generate(
			s.outputDir,
			guid,
			devices,
		)
		if err != nil {
			s.markFailed(ctx, fileName, err)
			return err
		}
	}

	if err := s.processedFileRepo.UpdateStatus(
		ctx,
		fileName,
		model.FileStatusCompleted,
		nil,
	); err != nil {
		return err
	}

	return nil
}

func (s *FileService) markFailed(
	ctx context.Context,
	fileName string,
	err error,
) {
	errMsg := err.Error()

	_ = s.processedFileRepo.UpdateStatus(
		ctx,
		fileName,
		model.FileStatusFailed,
		&errMsg,
	)
}
