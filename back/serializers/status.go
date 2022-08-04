package serializers

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sti26/ovpn_freeipa_mgmt/models"
)

func StatusSerialazer(s string) *gin.H {
	status := models.Status{}

	rows := strings.Split(s, "\n")
	for i := 0; i < len(rows); i++ {

		switch {
		case strings.HasPrefix(rows[i], "TITLE"):
			status.Title = rows[i][6:]
		case strings.HasPrefix(rows[i], "CLIENT_LIST"):
			var client models.Client

			row := strings.Split(rows[i], ",")

			if len(row) == 12 {
				// new version

				client = models.Client{
					CommonName:         row[1],
					RealAddress:        row[2],
					VirtualAddress:     row[3],
					VirtualIPv6Address: row[4],
					BytesReceived:      row[5],
					BytesSent:          row[6],
					ConnectedSince:     row[7],
					Username:           row[9],
					ClientID:           row[10],
					PeerID:             row[11],
				}
			} else if len(row) == 9 {
				// old version

				client = models.Client{
					CommonName:     row[1],
					RealAddress:    row[2],
					VirtualAddress: row[3],
					BytesReceived:  row[4],
					BytesSent:      row[5],
					ConnectedSince: row[6],
					Username:       row[8],
				}
			} else {
				client = models.Client{CommonName: "parce error"}
			}

			status.ClientList = append(status.ClientList, client)
		}
	}

	return &gin.H{
		"error":  "",
		"status": status,
	}
}
