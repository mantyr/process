package process

import (
	"context"
)

// Environment это окружение для контролируемого запуска команд
type Environment interface {
	SetDir(dir string) Environment
	SetEnv(key, value string) Environment
	GetEnvs() []string
	GetEnv(key string) string

	RunCommand(ctx context.Context, cmd string) ([]byte, error)
}
