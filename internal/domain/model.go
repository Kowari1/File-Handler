package model

import "github.com/google/uuid"

type Device struct {
	N         int
	MQTT      string
	Invid     string
	UnitGUID  uuid.UUID
	MsgID     string
	Text      string
	Context   string
	Class     int
	Level     string
	Area      string
	Addr      string
	Block     string
	Type      string
	Bit       string
	InvertBit string
}
