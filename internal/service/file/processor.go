package file

import (
	"context"

	model "github.com/Kowari1/File-Handler/internal/domain"
	"github.com/Kowari1/File-Handler/internal/parser"
	"github.com/Kowari1/File-Handler/internal/repository"
	"github.com/google/uuid"
)

type Processor struct {
	deviceRepo   repository.DeviceRepository
	parseErrRepo repository.ParseErrorRepository
	parser       parser.Parser
}

func NewProcessor(
	deviceRepo repository.DeviceRepository,
	parseErrRepo repository.ParseErrorRepository,
	parser parser.Parser,
) *Processor {
	return &Processor{
		deviceRepo:   deviceRepo,
		parseErrRepo: parseErrRepo,
		parser:       parser,
	}
}

func (p *Processor) Process(
	ctx context.Context,
	fileName string,
	lines <-chan string,
) (map[uuid.UUID][]*model.Device, error) {

	const batchSize = 1000

	guidMap := make(map[uuid.UUID][]*model.Device)
	devices := make([]*model.Device, 0, batchSize)
	lineNumber := 0

	for line := range lines {
		lineNumber++

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		device, err := p.parser.Parse(line, lineNumber)
		if err != nil {

			if perr, ok := err.(*parser.ParseError); ok {
				if err := p.parseErrRepo.Save(
					ctx,
					fileName,
					perr.Line,
					perr.Message,
				); err != nil {
					return nil, err
				}
				continue
			}

			return nil, err
		}

		if device == nil {
			continue
		}

		devices = append(devices, device)
		guidMap[device.UnitGUID] = append(guidMap[device.UnitGUID], device)

		if len(devices) >= batchSize {
			if err := p.deviceRepo.Save(ctx, devices); err != nil {
				return nil, err
			}
			devices = devices[:0]
		}
	}

	if len(devices) > 0 {
		if err := p.deviceRepo.Save(ctx, devices); err != nil {
			return nil, err
		}
	}

	return guidMap, nil
}
