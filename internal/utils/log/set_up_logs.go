package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/seqra/seqra/internal/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	logFile *os.File
)

// OpenLogFile creates and returns a file for logging at the specified path.
// It creates the directory structure if it doesn't exist.
// The file handle is stored in a global variable and can be closed with CloseLogFile().
func OpenLogFile() (*os.File, string, error) {
	seqraHomePath, err := utils.GetSeqraHome()
	cobra.CheckErr(err)
	logDir := filepath.Join(seqraHomePath, "logs")
	// Create log file with timestamp
	logPath := filepath.Join(logDir, time.Now().Format("2006-01-02_15-04-05.log"))

	// Ensure the directory exists
	if err := os.MkdirAll(filepath.Dir(logPath), 0755); err != nil {
		return nil, "", err
	}

	// Open or create the log file
	logFile, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		return nil, "", err
	}

	return logFile, logPath, nil
}

// SetUpLogs configures logging with the specified output and level.
// If logFilePath is provided, logs will also be written to the specified file
// in addition to the provided output writer.
func SetUpLogs(out io.Writer, level string) error {
	// Parse log level
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}

	// Set the log level
	logrus.SetLevel(lvl)

	// If no log file is specified, just use the provided output
	logrus.SetOutput(out)
	logrus.AddHook(&plainHook{w: os.Stdout})

	return nil
}

// plainHook prints only Entry.Message for InfoLevel (no timestamp / level)
type plainHook struct{ w io.Writer }

func (h *plainHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
	}
}

func (h *plainHook) Fire(e *logrus.Entry) error {
	_, err := fmt.Fprintln(h.w, e.Message)
	return err
}

// CloseLogFile closes the log file if it's open.
// This should be called when the application is shutting down.
func CloseLogFile() error {
	if logFile != nil {
		return logFile.Close()
	}
	return nil
}
