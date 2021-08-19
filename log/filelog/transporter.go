package filelog

import (
	"fmt"
	"github.com/eolinker/eosc"
	eosc_log "github.com/eolinker/eosc/log"
	"github.com/eolinker/goku/log"
	"time"
)

type Transporter struct {
	*eosc_log.Transporter
	writer *FileWriterByPeriod
}

func (t *Transporter) Close() error {
	t.writer.Close()
	return nil
}

func (t *Transporter) Reset(c interface{}, f eosc_log.Formatter) error {
	conf, ok := c.(*Config)
	if !ok {
		return fmt.Errorf("need %s,now %s", eosc.TypeNameOf((*Config)(nil)), eosc.TypeNameOf(c))
	}

	t.Transporter.SetFormatter(f)
	return t.reset(conf)
}

func (t *Transporter) reset(c *Config) error {
	t.SetOutput(t.writer)
	t.SetLevel(c.Level)

	t.writer.Set(
		c.Dir,
		c.File,
		c.Period,
		time.Duration(c.Expire)*time.Hour*24,
	)
	t.writer.Open()
	return nil
}

func createTransporter(conf *Config, formatter eosc_log.Formatter) (log.TransporterReset, error) {

	fileWriterByPeriod := NewFileWriteByPeriod()

	transport := &Transporter{
		Transporter: eosc_log.NewTransport(fileWriterByPeriod, conf.Level, formatter),
		writer:      fileWriterByPeriod,
	}

	e := transport.Reset(conf, formatter)
	if e != nil {
		return nil, e
	}
	return transport, nil
}
