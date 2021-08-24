//+build wireinject

package di

import (
	server1 "bingu/test/server"
	"github.com/go-kirito/pkg/application"
	"github.com/google/wire"
)

func RegisterService(app application.Application) error {
	wire.Build(
		server1.NewStrings,
		server1.NewName,
		server1.NewUser)

	return nil
}
