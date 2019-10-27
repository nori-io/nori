package log

import (
	"github.com/nori-io/nori-common/logger"
	"github.com/sirupsen/logrus"
)

type Logger struct {
	logger *logrus.Logger
	entry  *logrus.Entry
}

func New() logger.Logger {
	return &Logger{
		logger: logrus.New(),
		entry:  nil,
	}
}

func (l *Logger) Debug(args ...interface{}) {
	if l.logger != nil {
		l.logger.Debug(args...)
	} else if l.entry != nil {
		l.entry.Debug(args...)
	}
}

func (l *Logger) Info(args ...interface{}) {
	if l.logger != nil {
		l.logger.Info(args...)
	} else if l.entry != nil {
		l.entry.Info(args...)
	}
}

func (l *Logger) Print(args ...interface{}) {
	if l.logger != nil {
		l.logger.Print(args...)
	} else if l.entry != nil {
		l.entry.Print(args...)
	}
}

func (l *Logger) Warn(args ...interface{}) {
	if l.logger != nil {
		l.logger.Warn(args...)
	} else if l.entry != nil {
		l.entry.Warn(args...)
	}
}

func (l *Logger) Warning(args ...interface{}) {
	if l.logger != nil {
		l.logger.Warning(args...)
	} else if l.entry != nil {
		l.entry.Warning(args...)
	}
}

func (l *Logger) Error(args ...interface{}) {
	if l.logger != nil {
		l.logger.Error(args...)
	} else if l.entry != nil {
		l.entry.Error(args...)
	}
}

func (l *Logger) Fatal(args ...interface{}) {
	if l.logger != nil {
		l.logger.Fatal(args...)
	} else if l.entry != nil {
		l.entry.Fatal(args...)
	}
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	if l.logger != nil {
		l.logger.Debugf(format, args...)
	} else if l.entry != nil {
		l.entry.Debugf(format, args...)
	}
}

func (l *Logger) Infof(format string, args ...interface{}) {
	if l.logger != nil {
		l.logger.Infof(format, args...)
	} else if l.entry != nil {
		l.entry.Infof(format, args...)
	}
}

func (l *Logger) Printf(format string, args ...interface{}) {
	if l.logger != nil {
		l.logger.Printf(format, args...)
	} else if l.entry != nil {
		l.entry.Printf(format, args...)
	}
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	if l.logger != nil {
		l.logger.Warnf(format, args...)
	} else if l.entry != nil {
		l.entry.Warnf(format, args...)
	}
}

func (l *Logger) Warningf(format string, args ...interface{}) {
	if l.logger != nil {
		l.logger.Warningf(format, args...)
	} else if l.entry != nil {
		l.entry.Warningf(format, args...)
	}
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	if l.logger != nil {
		l.logger.Errorf(format, args...)
	} else if l.entry != nil {
		l.entry.Errorf(format, args...)
	}
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	if l.logger != nil {
		l.logger.Fatalf(format, args...)
	} else if l.entry != nil {
		l.entry.Fatalf(format, args...)
	}
}

func (l *Logger) Debugln(args ...interface{}) {
	if l.logger != nil {
		l.logger.Debugln(args...)
	} else if l.entry != nil {
		l.entry.Debugln(args...)
	}
}

func (l *Logger) Infoln(args ...interface{}) {
	if l.logger != nil {
		l.logger.Infoln(args...)
	} else if l.entry != nil {
		l.entry.Infoln(args...)
	}
}

func (l *Logger) Println(args ...interface{}) {
	if l.logger != nil {
		l.logger.Println(args...)
	} else if l.entry != nil {
		l.entry.Println(args...)
	}
}

func (l *Logger) Warnln(args ...interface{}) {
	if l.logger != nil {
		l.logger.Warnln(args...)
	} else if l.entry != nil {
		l.entry.Warnln(args...)
	}
}

func (l *Logger) Warningln(args ...interface{}) {
	if l.logger != nil {
		l.logger.Warningln(args...)
	} else if l.entry != nil {
		l.entry.Warningln(args...)
	}
}

func (l *Logger) Errorln(args ...interface{}) {
	if l.logger != nil {
		l.logger.Errorln(args...)
	} else if l.entry != nil {
		l.entry.Errorln(args...)
	}
}

func (l *Logger) Fatalln(args ...interface{}) {
	if l.logger != nil {
		l.logger.Fatalln(args...)
	} else if l.entry != nil {
		l.entry.Fatalln(args...)
	}
}

func (l *Logger) WithField(key string, value interface{}) logger.Logger {
	var entry *logrus.Entry
	if l.logger != nil {
		entry = l.logger.WithField(key, value)
	} else if l.entry != nil {
		entry = l.entry.WithField(key, value)
	}
	return &Logger{
		entry: entry,
	}
}

func (l *Logger) WithFields(fields logger.Fields) logger.Logger {
	flds := logrus.Fields{}
	for k, v := range fields {
		flds[k] = v
	}
	var entry *logrus.Entry
	if l.logger != nil {
		entry = l.logger.WithFields(flds)
	} else if l.entry != nil {
		entry = l.entry.WithFields(flds)
	}
	return &Logger{
		entry: entry,
	}
}
