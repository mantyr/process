package process

import (
	"os"
	"os/exec"
	"runtime"
)

type Context struct {
	dir   string
	name  string
	argv  []string
	env   map[string]string
	stdin struct {
		r *os.File
		w *os.File
	}
	stdout struct {
		r *os.File
		w *os.File
	}
	stderr struct {
		r *os.File
		w *os.File
	}
	// readinessProbe указывает на то ждать ли специального сигнала от приложения что бы считать приложение полностью запущенным или нет
	readinessProbe bool
}

func NewContext() *Context {
	return &Context{
		env: make(map[string]string),
	}
}

func (c *Context) SetReadinessProbe(status bool) *Context {
	c.readinessProbe = status
	return c
}

func (c *Context) SetCommand(command string) *Context {
	switch runtime.GOOS {
	case "windows":
		c.argv = []string{"cmd", "/S", "/C", command}
	default:
		c.argv = []string{"/bin/sh", "-c", command}
	}
	c.name = c.argv[0]
	return c
}

func (c *Context) SetCMD(name string, argv []string) *Context {
	c.name = name
	c.argv = argv
	return c
}

func (c *Context) SetDir(dir string) *Context {
	c.dir = dir
	return c
}

func (c *Context) SetEnv(key, value string) *Context {
	c.env[key] = value
	return c
}

func (c *Context) GetEnv(key string) string {
	v, _ := c.env[key]
	return v
}

func (c *Context) GetEnvs() []string {
	result := make([]string, 0, len(c.env))
	for key, value := range c.env {
		result = append(result, key+"="+value)
	}
	return result
}

func (c *Context) EnableStdin() (err error) {
	if c.stdin.r != nil {
		return nil
	}
	c.stdin.r, c.stdin.w, err = os.Pipe()
	return err
}

func (c *Context) EnableStdout() (err error) {
	if c.stdout.r != nil {
		return nil
	}
	c.stdout.r, c.stdout.w, err = os.Pipe()
	return err
}

func (c *Context) EnableStderr() (err error) {
	if c.stderr.r != nil {
		return nil
	}
	c.stderr.r, c.stderr.w, err = os.Pipe()
	return err
}

func (c *Context) Stdin() *os.File {
	return c.stdin.w
}

func (c *Context) Stdout() *os.File {
	return c.stdout.r
}

func (c *Context) Stderr() *os.File {
	return c.stderr.r
}

func (c *Context) StartProcess() (*os.Process, error) {
	path, err := exec.LookPath(c.name)
	if err != nil {
		return nil, err
	}
	attr := &os.ProcAttr{
		Dir:   c.dir,
		Env:   c.GetEnvs(),
		Files: []*os.File{c.stdin.r, c.stdout.w, c.stderr.w},
	}
	return os.StartProcess(path, c.argv, attr)
}
