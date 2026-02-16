package parser

import (
	"strconv"
	"strings"

	model "github.com/Kowari1/File-Handler/internal/domain"
	"github.com/google/uuid"
)

const (
	// Количество полей в записи.
	FieldsCount = 15

	// Сканирование файла начинаем с первой строки.
	DefaultLineNumber = 0

	// Номер строки с гуидом.
	GuidLineNumber = 3
)

type Parser interface {
	Parse(line string, lineNumber int) (*model.Device, error)
}

type TSVParser struct{}

func NewTSVParser() *TSVParser {
	return &TSVParser{}
}

func (p *TSVParser) Parse(
	line string,
	lineNumber int,
) (*model.Device, error) {

	fields := strings.Split(line, "\t")

	if len(fields) < FieldsCount {
		return nil, &ParseError{
			Line:    lineNumber,
			Message: "invalid field count",
		}
	}

	nValue, err := parseInt(fields[0], lineNumber, "N")
	if err != nil {
		return nil, err
	}

	classValue, err := parseInt(fields[7], lineNumber, "class")
	if err != nil {
		return nil, err
	}

	unitGUID, err := parseUUID(fields[3], lineNumber)
	if err != nil {
		return nil, err
	}

	return &model.Device{
		N:         nValue,
		MQTT:      fields[1],
		Invid:     fields[2],
		UnitGUID:  unitGUID,
		MsgID:     fields[4],
		Text:      fields[5],
		Context:   fields[6],
		Class:     classValue,
		Level:     fields[8],
		Area:      fields[9],
		Addr:      fields[10],
		Block:     fields[11],
		Type:      fields[12],
		Bit:       fields[13],
		InvertBit: fields[14],
	}, nil
}

func parseInt(value string, line int, field string) (int, error) {
	v, err := strconv.Atoi(value)
	if err != nil {
		return 0, &ParseError{
			Line:    line,
			Message: "invalid " + field + " value",
		}
	}
	return v, nil
}

func parseUUID(value string, line int) (uuid.UUID, error) {
	id, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, &ParseError{
			Line:    line,
			Message: "invalid unit_guid",
		}
	}
	return id, nil
}
