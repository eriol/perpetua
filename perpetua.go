package perpetua

var options Options
var store Store

func Start() {

	options.read()
	startIRC()
}
