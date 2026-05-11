package logger

import (
	"github.com/Compogo/compogo/logger"
	"github.com/Compogo/types/linker"
)

var (
	LevelToProtobuf = linker.NewLinker[logger.Level, Level](
		linker.Link(logger.Panic, Level_Panic),
		linker.Link(logger.Error, Level_Error),
		linker.Link(logger.Warn, Level_Warn),
		linker.Link(logger.Info, Level_Info),
		linker.Link(logger.Debug, Level_Debug),
	)

	ProtobufToLevel = linker.NewLinker[Level, logger.Level](
		linker.Link(Level_Panic, logger.Panic),
		linker.Link(Level_Error, logger.Error),
		linker.Link(Level_Warn, logger.Warn),
		linker.Link(Level_Info, logger.Info),
		linker.Link(Level_Debug, logger.Debug),
	)
)
