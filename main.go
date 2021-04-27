package main

import (
	"net/url"
	"strconv"

	"github.com/zserge/lorca"

	adb "github.com/zach-klippenstein/goadb"
)

const MainUI string = `
<!doctype html>
<html>
	<head>
		<title>Android Coroner</title>
		<style>
			html,body {box-sizing: border-box;background-color:#111; color:#ccc; font: 12pt "Consolas",monospace;}
			body {width:100%;height:100%; overflow:hidden}
			#log {position:absolute; height: 50vh;bottom:0;}
			#log > div {padding: 1pt 12pt}
			.log-error {background-color: #500;}
			.log-fatal {
				padding: 1cm;
				border: 3px solid #a00;
				background-color: #500;
				color: #eee;
				font-size: 300%;
				text-align: center;

				max-width: 75%;
				position: absolute;
				top: 50%;
				left: 50%;
				margin-right: -50%;
				transform: translate(-50%, -50%)
			}
		</style>
	</head>
	<body>
		<div id="log"></div>
		<script>
			document.body.addEventListener("keyup", function (e) {if (e.key == "Escape") {quit()}});
			document.body.addEventListener("keyup", function (e) {if (e.key == "u") {uplot()}});
		</script>
	</body>
</html>
`

func main() {
	lorcaUI, err := lorca.New("", "", 640, 480)
	mustnot(err)
	defer lorcaUI.Close()

	ui := UI{lorcaUI}

	must(ui.Bind("quit", ui.Close))

	must(ui.Load("data:text/html," + url.PathEscape(MainUI)))
	mustfn(ui.SetBounds(lorca.Bounds{0, 0, 640, 480, lorca.WindowStateFullscreen}), ui.fatalfn)

	adbPort := 15307
	adbClient, e := adb.NewWithConfig(adb.ServerConfig{Port: adbPort})
	mustnotfn(e, ui.fatalfn)

	ui.log("starting adb server on port", strconv.Itoa(adbPort))
	mustfn(adbClient.StartServer(), ui.fatalfn)
	defer must(adbClient.KillServer())

	serverVersion, err := adbClient.ServerVersion()
	mustnotfn(err, ui.fatalfn)
	ui.log("adb server version:", strconv.Itoa(serverVersion))

	dm := DeviceManager{adbClient}
	go dm.watchWithUI(ui)

	<-ui.Done()
}
