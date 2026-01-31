package output

import (
	"io"
	"os"

	"github.com/ExecutiveOrder6102/postgres-bloat/internal/bloat"
	"github.com/ExecutiveOrder6102/postgres-bloat/internal/config"
)

func Write(cfg config.Config, indexes []bloat.IndexBloat, tables []bloat.TableBloat, concurrentSupported bool) error {
	var writer io.Writer = os.Stdout
	if cfg.OutputFile != "" {
		file, err := os.Create(cfg.OutputFile)
		if err != nil {
			return err
		}
		defer file.Close()
		writer = file
	}

	switch cfg.Output {
	case "console":
		return WriteConsole(writer, indexes, tables, concurrentSupported)
	case "csv":
		return WriteCSV(writer, indexes, tables, concurrentSupported)
	default:
		return nil
	}
}
