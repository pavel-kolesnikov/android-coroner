package main

import (
	"log"
	"strings"

	"github.com/zserge/lorca"
)

type UI struct {
	lorca.UI
}

func (ui UI) log(message ...string) {
	m := strings.Join(message, " ")
	ui.Eval(`document.getElementById('log').insertAdjacentHTML('afterbegin', '<div class="log-info">` + m + `</div>')`)
}

func (ui UI) error(message ...string) {
	m := strings.Join(message, " ")
	ui.Eval(`document.getElementById('log').insertAdjacentHTML('afterbegin', '<div class="log-error">` + m + `</div>')`)
}

func (ui UI) errorfn(message string) {
	ui.error(message)
}

func (ui UI) fatal(message ...string) {
	m := strings.Join(message, " ")
	ui.Eval(`document.body.insertAdjacentHTML('afterbegin', '<div class="log-fatal" onclick="quit()">` + m + `</div>')`)
	<-ui.Done()
	log.Fatal(message)
}

func (ui UI) fatalfn(message string) {
	ui.fatal(message)
}
