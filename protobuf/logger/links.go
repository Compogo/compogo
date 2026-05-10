package logger

import (
	"github.com/Compogo/compogo/logger"
	"github.com/Compogo/types/linker"
)

var (
	LevelToProtobuf = linker.NewLinker[logger.Level, Level](
		linker.NewLink(logger.Panic, Level_Panic),
		linker.NewLink(logger.Error, Level_Error),
		linker.NewLink(logger.Warn, Level_Warn),
		linker.NewLink(logger.Info, Level_Info),
		linker.NewLink(logger.Debug, Level_Debug),
	)

	ProtobufToLevel = linker.NewLinker[Level, logger.Level](
		linker.NewLink(Level_Panic, logger.Panic),
		linker.NewLink(Level_Error, logger.Error),
		linker.NewLink(Level_Warn, logger.Warn),
		linker.NewLink(Level_Info, logger.Info),
		linker.NewLink(Level_Debug, logger.Debug),
	)
)
