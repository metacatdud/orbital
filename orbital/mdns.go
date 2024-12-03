package orbital

import (
	"github.com/hashicorp/mdns"
)

type MDNSService struct {
	services map[string]*mdns.Server
}

func (s *MDNSService) Register(name, host string, port int) error {
	//info := []string{fmt.Sprintf("Registering %s:%d", host, port)}
	//entry := mdns.MDNSService{
	//	Instance: name,
	//	Service:  "_orbital._tcp",
	//	Domain:   "local",
	//	HostName: ,
	//	Port:     0,
	//	IPs:      nil,
	//	TXT:      info,
	//}
	return nil
}

func NewMDNSService() *MDNSService {
	return &MDNSService{
		services: make(map[string]*mdns.Server),
	}
}
