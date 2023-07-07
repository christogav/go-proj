package logging

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Module(
		"logging",

		fx.Provide(newZeroLogger),
		fx.Provide(newConfig),
	),
	fx.WithLogger(newFxLogger),
)
