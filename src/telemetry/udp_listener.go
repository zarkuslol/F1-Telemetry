// --- telemetry/udp_listener.go ---
package telemetry

import (
	"joselucas/f1-telemetry/src/utils"
	"net"
	"strconv"
)

func ListenUDP(port int) (*net.UDPConn, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(port))
	if err != nil {
		utils.Logger.Printf("Falha ao resolver endereço UDP: %v", err)
		return nil, err
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		utils.Logger.Printf("Falha ao ouvir na porta UDP: %v", err)
		return nil, err
	}
	return conn, nil
}
