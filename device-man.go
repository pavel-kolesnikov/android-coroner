package main

import (
	"fmt"
	"io/ioutil"
	"time"

	adb "github.com/zach-klippenstein/goadb"
)

type DeviceManager struct {
	*adb.Adb
}

func (dm *DeviceManager) watchWithUI(ui UI) {
	makeDeviceWatcher := func() {
		watcher := dm.NewDeviceWatcher()
		defer watcher.Shutdown()
		ui.log("device watcher started")

		for event := range watcher.C() {
			if event.CameOnline() {
				dm.handleDeviceCameOnline(event.Serial, ui)
			} else if event.WentOffline() {
				ui.log(event.Serial + " offline")
			}
		}
		if e := watcher.Err(); e != nil {
			ui.error("device watcher dies: ", e.Error())
		}
	}

	for {
		select {
		case <-ui.Done():
			return
		default:
			makeDeviceWatcher()
			time.Sleep(1 * time.Second)
		}
	}
}

func pre(s string) string {
	return "<pre>" + s + "</pre>"
}

func (dm *DeviceManager) handleDeviceCameOnline(serial string, ui UI) {
	ui.log(serial + " online")

	device := dm.Device(adb.DeviceWithSerial(serial))

	cmd := "dumpsys -l"
	ui.log(cmd)
	cmdOutput, err := device.RunCommand(cmd)
	if err != nil {
		ui.error(fmt.Sprintf("error running command `%s`: %s", cmd, err))
	}
	ui.log(pre(cmdOutput))

	path := "/data/tombstones"
	ui.log(fmt.Sprintf(`files in "%s":`, path))
	entries, err := device.ListDirEntries(path)
	if err != nil {
		ui.error("error listing files:", err.Error())
	} else {
		for entries.Next() {
			ui.log(fmt.Sprintf("<div>%+v</div>", *entries.Entry()))
		}
		if entries.Err() != nil {
			ui.error("error listing files:", err.Error())
		}
	}

	loadavgReader, err := device.OpenRead("/proc/loadavg")
	if err != nil {
		ui.error("error opening file:", err.Error())
	} else {
		loadAvg, err := ioutil.ReadAll(loadavgReader)
		if err != nil {
			ui.error("error reading file:", err.Error())
		} else {
			ui.log("loadavg", string(loadAvg))
		}
	}
}
