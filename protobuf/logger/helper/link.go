package helper

import (
	"github.com/Compogo/compogo/logger"
	protobuf "github.com/Compogo/compogo/protobuf/logger"
	"github.com/Compogo/types/linker"
)

var (
	LevelLinkProtoLevel = linker.NewLinker[logger.Level, protobuf.Level](
		linker.Link(logger.Panic, protobuf.Level_Panic),
		linker.Link(logger.Error, protobuf.Level_Error),
		linker.Link(logger.Warn, protobuf.Level_Warn),
		linker.Link(logger.Info, protobuf.Level_Info),
		linker.Link(logger.Debug, protobuf.Level_Debug),
	)

	ProtoLevelLinkLevel = linker.NewLinker[protobuf.Level, logger.Level](
		linker.Link(protobuf.Level_Panic, logger.Panic),
		linker.Link(protobuf.Level_Error, logger.Error),
		linker.Link(protobuf.Level_Warn, logger.Warn),
		linker.Link(protobuf.Level_Info, logger.Info),
		linker.Link(protobuf.Level_Debug, logger.Debug),
	)
)
