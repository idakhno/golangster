package special_chars

import "log/slog"

func bad() {
	slog.Info("server started 🚀")        // want `log message must not contain emoji`
	slog.Error("connection failed!!!")    // want `log message must not contain special character '!'`
	slog.Warn("something went wrong...")  // want `log message must not contain`
	slog.Debug("are you sure?")           // want `log message must not contain special character`
	slog.Info("hello world!")             // want `log message must not contain special character '!'`
}

func good() {
	slog.Info("server started")
	slog.Error("connection failed")
	slog.Warn("something went wrong")
	slog.Debug("request processed")
	slog.Info("version 1.2.3 deployed")
}
