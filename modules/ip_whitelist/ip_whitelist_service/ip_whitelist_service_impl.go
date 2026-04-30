package ip_whitelist_service

import (
	"core-ticket/modules/ip_whitelist/ip_whitelist_repository"
	"errors"
	"fmt"
	"net"
	"strings"
)

type IpWhitelistServiceImpl struct {
	IpWhitelistRepository ip_whitelist_repository.IpWhitelistRepository
}

func NewIpWhitelistService(
	ipWhitelistRepo ip_whitelist_repository.IpWhitelistRepository,
) IpWhitelistService {
	return &IpWhitelistServiceImpl{
		IpWhitelistRepository: ipWhitelistRepo,
	}
}

func (service *IpWhitelistServiceImpl) InWhitelist(clientIpAddress string) (bool, error) {
	data, err := service.IpWhitelistRepository.FindAll()
	fmt.Println(data)
	if err != nil {
		return false, err
	}

	var ipWhitelist []string
	for _, ipW := range data {
		ipWhitelist = append(ipWhitelist, ipW.IpAddress)
	}

	ipNets := getIpNets(ipWhitelist)

	clientIP := net.ParseIP(clientIpAddress)
	for _, ipNet := range ipNets {
		if ipNet.Contains(clientIP) {
			return true, nil
		}
	}

	return false, errors.New("Forbidden")
}

func getIpNets(ipWhitelist []string) []*net.IPNet {
	var ipNets []*net.IPNet
	for _, s := range ipWhitelist {
		_, ipNet, err := net.ParseCIDR(strings.TrimSpace(s))
		if err != nil {
			ip := net.ParseIP(strings.TrimSpace(s))
			if ip != nil {
				ipNets = append(ipNets, &net.IPNet{IP: ip, Mask: net.CIDRMask(32, 32)})
			}
		} else {
			ipNets = append(ipNets, ipNet)
		}
	}
	return ipNets
}
