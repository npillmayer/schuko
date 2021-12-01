package appender

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/npillmayer/schuko"
	"github.com/npillmayer/schuko/tracing"
)

func AppenderFromConfig(conf schuko.Configuration) (io.Writer, error) {
	if dest := conf.GetString("tracing.destination"); dest != "" {
		var err error
		var w io.Writer
		//fmt.Printf("@@@ dest/appender = %q\n", dest)
		tracing.Infof("opening tracing destination %q\n", dest)
		if w, err = Destination(dest); err != nil {
			err = fmt.Errorf("re-directing trace output failed: %w", err)
			tracing.Errorf(err.Error())
			return os.Stderr, err
		}
		return w, err
	}
	return os.Stderr, nil
}

// Destination opens a tracing destination as an io.Writer. dest may be one of
//
// a) literals "Stdout" or "Stderr"
//
// b) a file URI ("file: //my.log")
//
// More to come.
//
func Destination(dest string) (io.WriteCloser, error) {
	switch strings.ToLower(dest) {
	case "stdout":
		return os.Stdout, nil
	case "stderr":
		return os.Stderr, nil
	}
	u, err := url.Parse(dest)
	if err != nil {
		return os.Stderr, err
	}
	if strings.ToLower(u.Scheme) == "file" {
		fname := u.Path
		if fname == "" {
			fname = u.Host
		}
		f, err := os.OpenFile(fname, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err == nil {
			return f, nil
		}
		if errors.Is(err, fs.ErrNotExist) {
			if err := os.MkdirAll(filepath.Dir(fname), 0777); err != nil {
				return nil, err
			}
			return os.Create(fname)
		} else {
			return nil, err
		}
	}
	return os.Stderr, err
}
