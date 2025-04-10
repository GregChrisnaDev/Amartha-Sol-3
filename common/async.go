package common

var AsyncFunc = func(fn func()) {
	go fn()
}

// make async proccess to sync (for testing)
func AsyncMakeSync() {
	AsyncFunc = func(fn func()) {
		fn()
	}
}

func AsyncMakeDefault() {
	AsyncFunc = func(fn func()) {
		go fn()
	}
}
