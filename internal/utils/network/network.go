package network

import (
	"fmt"
	"net"
)

func IsAddressInNetwork(ipAddress string, networkCIDR string) (bool, error) {
	// Преобразуем строковые представления IP-адреса и CIDR в объекты типа net.IP и net.IPNet
	ip := net.ParseIP(ipAddress)
	_, network, err := net.ParseCIDR(networkCIDR)
	if err != nil {
		return false, err
	}
	// Проверяем, принадлежит ли IP-адрес сети
	return network.Contains(ip), nil
}

func ParseIP() (string, error) {
	address, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	var ip string
	for _, addr := range address {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				if inNetwork, err := IsAddressInNetwork(ipnet.IP.String(), "169.254.0.0/16"); err != nil {
					return "", fmt.Errorf("невозможно определить подсеть IP-адреса %w", err)
				} else if !inNetwork {
					ip = ipnet.IP.String()
					break
				}
			}
		}
	}
	if len(ip) == 0 {
		return "", fmt.Errorf("невозможно определить IP-адрес")
	}
	return ip, nil
}
