package manager


func Debug(err error) {
	if err != nil {
		Log.Debug(err.Error())
	}
}


func Info(msg string, args ...interface {}) {
	if msg != nil {
		Log.Info(msg, args...)
	}
}




