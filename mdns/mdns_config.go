package mdns

var DeviceIsHost bool = false

type TextRecord struct {
	IPAddress string
	Host      bool
}

const ServiceName string = "_clipboard-sync._tcp"
const ServicePort int = 8000

func IdentifyIsHost() error {
	entries, err := SearchDevices()
	if err != nil {
		return err
	}

	if len(entries) == 0 {
		DeviceIsHost = true
	}

	return nil
}
