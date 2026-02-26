package english

import "log/slog"

func bad() {
	slog.Info("запуск сервера")             // want `log message must be in English only`
	slog.Error("ошибка подключения к базе") // want `log message must be in English only`
	slog.Warn("предупреждение системы")     // want `log message must be in English only`
	slog.Debug("服务器启动失败")                   // want `log message must be in English only`
}

func good() {
	slog.Info("starting server")
	slog.Error("connection failed")
	slog.Warn("system warning")
	slog.Debug("request timeout after 30s")
	slog.Info("user john@example.com created")
}
