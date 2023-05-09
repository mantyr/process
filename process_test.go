package process

import (
	"bytes"
	"context"
	"io"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey" //nolint:golint,stylecheck
)

func TestProcess(t *testing.T) {
	Convey("Check process.Process", t, func() {
		config := NewContext()
		config.
			SetDir("./testdata").
			SetEnv("key1", "value1").
			SetEnv("key2", "value2")
		So(config.EnableStdout(), ShouldBeNil)
		So(config.EnableStderr(), ShouldBeNil)

		stdout := bytes.NewBuffer([]byte{})
		stderr := bytes.NewBuffer([]byte{})
		go io.Copy(stdout, config.Stdout())
		go io.Copy(stderr, config.Stderr())

		Convey("Test1 - Run", func() {
			config.SetCommand("./main test1")
			pr := NewProcess()
			pr.SetContext(*config)
			So(pr.Status(), ShouldEqual, NotRunning)
			So(pr.Run(context.Background()), ShouldBeNil)

			ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancelFunc()

			running := false
			for {
				select {
				case <-ctx.Done():
					break
				default:
				}
				status := pr.Status()
				if status == Running {
					running = true
					break
				}
				time.Sleep(1 * time.Second)
			}
			So(running, ShouldEqual, true)
			Convey("Stop", func() {
				So(pr.Status(), ShouldEqual, Running)
				So(pr.Stop(context.Background()), ShouldBeNil)

				ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancelFunc()

				notRunning := false
				for {
					select {
					case <-ctx.Done():
						break
					default:
					}
					status := pr.Status()
					if status == NotRunning {
						notRunning = true
						break
					}
					time.Sleep(1 * time.Second)
				}
				So(notRunning, ShouldEqual, true)
			})
		})
	})
}
