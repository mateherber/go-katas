package simulator

import (
	"errors"
	"log"
	"os/exec"
	"strings"

	"github.com/coreos/go-semver/semver"
)

const deviceName = "iPhone Xs Max"

func StartIosSimulator() {
	cmd := exec.Command("xcrun", "simctl", "list", "--json")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Couldn't execute simctl list %s\n", err)
	}
	simulators, err := UnmarshalSimulators(out)
	if err != nil {
		log.Fatalf("Couldn't parse simctl list response %s\n", err)
	}
	deviceType, err := simulators.getDeviceType(deviceName)
	if err != nil {
		log.Fatalf("Couldn't find device type with name %s\n", deviceName)
	}
	device, err := simulators.getDevice(deviceType.Identifier)
	var udid string
	if err != nil {
		runtime, err := simulators.getLastIosRuntime()
		if err != nil {
			log.Fatalf("Couldn't find iOS runtime")
		}
		cmd = exec.Command("xcrun", "simctl", "create", deviceName, deviceType.Identifier, runtime.Identifier)
		out, err = cmd.CombinedOutput()
		if err != nil {
			log.Fatalf("Couldn't create new device %s\n", err)
		}
		udidWords := strings.Fields(string(out))
		udid = udidWords[len(udidWords)-1]
	} else {
		udid = device.Udid
	}
	cmd = exec.Command("xcrun", "simctl", "bootstatus", udid, "-b")
	out, err = cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Couldn't start device %s\n", err)
	}
	log.Println("Simulator booted successfully")
}

func (s *Simulators) getLastIosRuntime() (*Runtime, error) {
	var maxIosRuntime *Runtime
	for _, r := range s.Runtimes {
		runtime := r
		if runtime.IsAvailable && strings.Contains(runtime.Identifier, "iOS") {
			if maxIosRuntime != nil {
				if semver.New(maxIosRuntime.Version).LessThan(*semver.New(runtime.Version)) {
					maxIosRuntime = &runtime
				}
			} else {
				maxIosRuntime = &runtime
			}
		}
	}
	if maxIosRuntime != nil {
		return maxIosRuntime, nil
	}
	return nil, errors.New("runtime was not found")
}

func (s *Simulators) getDeviceType(name string) (*DeviceType, error) {
	for _, deviceType := range s.DeviceTypes {
		if deviceType.Name == name {
			return &deviceType, nil
		}
	}
	return nil, errors.New("device type was not found")
}

func (s *Simulators) getDevice(deviceTypeIdentifier string) (*Device, error) {
	for _, devices := range s.Devices {
		for _, device := range devices {
			if device.IsAvailable && device.DeviceTypeIdentifier == deviceTypeIdentifier {
				return &device, nil
			}
		}
	}
	return nil, errors.New("device was not found")
}
