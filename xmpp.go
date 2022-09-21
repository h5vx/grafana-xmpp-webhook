package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"

	"gosrc.io/xmpp"
	"gosrc.io/xmpp/stanza"
)

var client *xmpp.Client

const presenceMucJoin string = `
<presence to='%s/%s'>
	<priority>1</priority>
	<c xmlns='http://jabber.org/protocol/caps' node='webhook'/>
	<x xmlns='http://jabber.org/protocol/muc'/>
</presence>
`

const presenceMucLeave string = `
<presence to='%s/%s' type='unavailable'/>
`

const iqVersionResponse string = `
<iq type='result' to='%s' id='%s'>
	<query xmlns='jabber:iq:version'>
		<name>%s</name>
		<version>%s (%s)</version>
	</query>
</iq>
`

func ConnectToXmpp() error {
	if !GlobalConfig.Xmpp.Enabled {
		fmt.Println("XMPP: Refuse to connect, XMPP is disabled")
		return nil
	}

	config := xmpp.Config{
		TransportConfiguration: xmpp.TransportConfiguration{
			Address:   GlobalConfig.Xmpp.Server,
			TLSConfig: &tls.Config{InsecureSkipVerify: GlobalConfig.Xmpp.TLS_Skip_Verify},
		},
		Jid:          GlobalConfig.Xmpp.Jid,
		Credential:   xmpp.Password(GlobalConfig.Xmpp.Password),
		StreamLogger: os.Stdout,
		Insecure:     !GlobalConfig.Xmpp.TLS,
	}

	router := xmpp.NewRouter()
	router.HandleFunc("iq", handleIq)
	_client, err := xmpp.NewClient(&config, router, errorHandler)

	client = _client

	if err != nil {
		return err
	}

	client.PostConnectHook = (func() error {
		RoomJoin(GlobalConfig.Xmpp.Room, GlobalConfig.Xmpp.Nickname)
		return nil
	})

	if err != nil {
		return err
	}

	go connectionManager(client)
	return nil
}

func RoomJoin(room string, nickname string) {
	presence := fmt.Sprintf(presenceMucJoin, room, nickname)
	client.SendRaw(presence)
}

func RoomSendMessage(room string, body string) {
	if !GlobalConfig.Xmpp.Enabled {
		fmt.Println("Refuse to send message: XMPP is disabled")
		return
	}

	message := stanza.NewMessage(stanza.Attrs{To: room, Type: "groupchat"})
	message.Body = body

	client.Send(message)
}

func RoomLeave(room string, nickname string) {
	presence := fmt.Sprintf(presenceMucLeave, room, nickname)
	client.SendRaw(presence)
}

func handleIq(s xmpp.Sender, p stanza.Packet) {
	iq := p.(*stanza.IQ)

	if iq.Payload != nil && iq.Payload.Namespace() == "jabber:iq:version" {
		if iq.From == iq.To {
			return
		}

		response := fmt.Sprintf(iqVersionResponse, iq.From, iq.Id, PackageName, PackageVersion, PackageCommitHash)
		client.SendRaw(response)
	}
}

func connectionManager(client *xmpp.Client) {
	cm := xmpp.NewStreamManager(client, nil)
	log.Fatal(cm.Run())
}

func errorHandler(err error) {
	fmt.Println(err.Error())
}
