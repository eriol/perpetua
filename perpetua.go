package perpetua

var options Options
var store Store

func Start() {

	options.read()

	store.open(DATABASE_FILE)
	defer store.close()

	startIRC()

}
