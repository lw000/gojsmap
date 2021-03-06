package models

import (
	"fmt"
	"gojsmap/global"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type ErrorItem struct {
	Source string `json:"source"`
	Name   string `json:"name"`
	Info   string `json:"info"`
	Line   int    `json:"line"`
	Column int    `json:"column"`
}

func (p ErrorItem) String() string {
	return fmt.Sprintf("Source: %s, Name: %s, Line:%d, Column:%d, ErrMsg:%s", p.Source, p.Name, p.Line, p.Column, p.Info)
}

type LogRecordConfig struct {
	module string
	line   string
	column string
	info   string
}

func (lrc *LogRecordConfig) SetColumn(column string) {
	lrc.column = column
}

func (lrc *LogRecordConfig) SetInfo(info string) {
	lrc.info = info
}

func (lrc *LogRecordConfig) SetModule(module string) {
	lrc.module = module
}

func (lrc *LogRecordConfig) SetLine(line string) {
	lrc.line = line
}

func NewLogConfig() *LogRecordConfig {
	return &LogRecordConfig{}
}

func (lrc *LogRecordConfig) Save(clientIP string) error {
	var (
		err error
		row int
		col int
	)

	row, err = strconv.Atoi(lrc.line)
	if err != nil {
		log.Error(err)
		return err
	}

	col, err = strconv.Atoi(lrc.column)
	if err != nil {
		log.Error(err)
		return err
	}

	go lrc.writeLog(row, col)

	return nil
}

func (lrc *LogRecordConfig) writeLog(row, col int) {
	defer func() {
		if x := recover(); x != nil {
		}
	}()

	source, name, line, column, ok := global.SourceMapManager.Get(row, col)
	if ok {
		p := &ErrorItem{Source: source, Name: name, Line: line, Column: column, Info: lrc.info}
		log.Error(p)
	} else {
		log.WithFields(log.Fields{"row": row, "col": col}).Error("unknown error")
	}
}
