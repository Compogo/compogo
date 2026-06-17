package helper

import (
	"github.com/Compogo/compogo"
	protobuf "github.com/Compogo/compogo/protobuf/logger"
	"github.com/Compogo/types/linker"
)

var (
	LevelLinkProtoLevel = linker.NewLinker[compogo.Level, protobuf.Level](
		linker.Link(compogo.Panic, protobuf.Level_Panic),
		linker.Link(compogo.Error, protobuf.Level_Error),
		linker.Link(compogo.Warn, protobuf.Level_Warn),
		linker.Link(compogo.Info, protobuf.Level_Info),
		linker.Link(compogo.Debug, protobuf.Level_Debug),
	)

	ProtoLevelLinkLevel = linker.NewLinker[protobuf.Level, compogo.Level](
		linker.Link(protobuf.Level_Panic, compogo.Panic),
		linker.Link(protobuf.Level_Error, compogo.Error),
		linker.Link(protobuf.Level_Warn, compogo.Warn),
		linker.Link(protobuf.Level_Info, compogo.Info),
		linker.Link(protobuf.Level_Debug, compogo.Debug),
	)
)
