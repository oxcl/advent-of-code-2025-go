package main

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

const FILENAME = "../input.txt"

type Device struct {
	name                string
	connectedDevices    map[string]*Device
	hasCachedRouteCount bool
	routesToOut         int
}

func NewDevice(name string) *Device {
	newDevice := &Device{
		name:                name,
		connectedDevices:    make(map[string]*Device),
		hasCachedRouteCount: false,
		routesToOut:         0,
	}
	return newDevice
}

type ServerRack struct {
	devices map[string]*Device
}

func ParseInput() (*ServerRack, error) {
	file, err := os.Open(FILENAME)
	if err != nil {
		return nil, err
	}
	serverRack := ServerRack{
		devices: make(map[string]*Device),
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		arr := strings.Split(line, ": ")
		if len(arr) < 2 {
			return nil, errors.New("there was an issue parsing this line:" + line)
		}
		deviceName := arr[0]
		connectedDevicesString := arr[1]
		connectedDeviceNames := strings.Split(connectedDevicesString, " ")
		if len(connectedDeviceNames) < 1 {
			return nil, errors.New("devices contianed no connected devices which is crazyyyy")
		}
		device, deviceExists := serverRack.devices[deviceName]
		if !deviceExists {
			newDevice := NewDevice(deviceName)
			serverRack.devices[deviceName] = newDevice
			device = newDevice
		}
		for _, connectedDeviceName := range connectedDeviceNames {
			connectedDevice, connectedDeviceExists := serverRack.devices[connectedDeviceName]
			if !connectedDeviceExists {
				newDevice := NewDevice(connectedDeviceName)
				serverRack.devices[connectedDeviceName] = newDevice
				connectedDevice = newDevice
			}
			device.connectedDevices[connectedDeviceName] = connectedDevice
		}
	}
	return &serverRack, nil
}

func (serverRack *ServerRack) Solve(deviceName string) (int, error) {
	device, deviceExists := serverRack.devices[deviceName]
	if !deviceExists {
		return 0, errors.New("device named " + deviceName + " doesn't exist")
	}
	if device.hasCachedRouteCount {
		return device.routesToOut, nil
	}
	totalRoutes := 0
	for deviceName := range device.connectedDevices {
		if deviceName == "out" {
			totalRoutes++
		} else {
			routes, err := serverRack.Solve(deviceName)
			if err != nil {
				return 0, err
			}
			totalRoutes += routes
		}
	}
	device.hasCachedRouteCount = true
	device.routesToOut = totalRoutes
	return totalRoutes, nil
}

func main() {
	serverRack, err := ParseInput()
	if err != nil {
		println(err)
		return
	}
	result, err := serverRack.Solve("you")
	if err != nil {
		println("Failed to solve the rack: " + err.Error())
		return
	}
	println("Result:", result)
}
