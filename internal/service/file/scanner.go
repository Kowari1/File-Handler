package file

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Kowari1/File-Handler/internal/repository"
	"go.uber.org/zap"
)

type Scanner struct {
	inputDir string
	interval time.Duration
	repo     repository.ProcessedFileRepository
	queue    chan string
	log      *zap.Logger
}

func NewScanner(
	inputDir string,
	interval time.Duration,
	repo repository.ProcessedFileRepository,
	queue chan string,
	log *zap.Logger,
) *Scanner {
	return &Scanner{
		inputDir: inputDir,
		interval: interval,
		repo:     repo,
		queue:    queue,
		log:      log,
	}
}

func (s *Scanner) Start(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.scan(ctx)
		}
	}
}

func (s *Scanner) scan(ctx context.Context) {
	files, err := os.ReadDir(s.inputDir)
	if err != nil {
		s.log.Error("failed to read input directory", zap.Error(err))
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if !strings.HasSuffix(file.Name(), ".tsv") {
			continue
		}

		exists, err := s.repo.Exists(ctx, file.Name())
		if err != nil {
			s.log.Error("failed to check file existence", zap.Error(err))
			continue
		}

		if exists {
			continue
		}

		fullPath := filepath.Join(s.inputDir, file.Name())

		select {
		case s.queue <- fullPath:
			s.log.Info("file added to queue", zap.String("file", file.Name()))
		default:
			s.log.Warn("queue is full, skipping file", zap.String("file", file.Name()))
		}
	}
}
