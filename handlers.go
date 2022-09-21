package main

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Alert struct {
	// Status is current status of the alert, firing or resolved
	Status string `json:"status" xml:"status" form:"status"`

	// Labels that are part of this alert, map of string keys to string values
	Labels map[string]interface{} `json:"labels" xml:"labels" form:"labels"`

	// StartsAt is the start time of the alert
	StartsAt time.Time `json:"startsAt" xml:"startsAt" form:"startsAt"`

	// EndsAt is the end time of the alert, default value when not resolved is 0001-01-01T00:00:00Z
	EndsAt time.Time `json:"endsAt" xml:"endsAt" form:"endsAt"`

	// Annotations that are part of this alert, map of string keys to string values
	Annotations map[string]interface{} `json:"annotations" xml:"annotations" form:"annotations"`

	// ValueString contains values that triggered the current status
	ValueString string `json:"valueString" xml:"valueString" form:"valueString"`
}

type AlertPayload struct {
	// Receiver is name of the webhook
	Receiver string `json:"receiver" xml:"receiver" form:"receiver"`

	// Status is current status of the alert, firing or resolved
	Status string `json:"status" xml:"status" form:"status"`

	// OrgID is ID of the organization related to the payload
	OrgID int `json:"orgId" xml:"orgId" form:"orgId"`

	// ExternalURL to the Grafana instance sending this webhook
	ExternalURL string `json:"externalURL" xml:"externalURL" form:"externalURL"`

	// Version of the payload
	Version string `json:"version" xml:"version" form:"version"`

	Alerts []Alert `json:"alerts" xml:"alerts" form:"alerts"`

	// Message contains the whole message about alert. It will be deprecated soon
	Message string `json:"message" xml:"message" form:"message"`
}

func AlertHandler(c *fiber.Ctx) error {
	c.Accepts("json", "text")
	c.AcceptsCharsets("utf-8")

	alert := new(AlertPayload)

	if err := c.BodyParser(alert); err != nil {
		return err
	}

	tmpl, err := template.New("alert").Parse(GlobalConfig.Alerting.Template)

	if err != nil {
		return err
	}

	var renderBuf bytes.Buffer

	if err := tmpl.Execute(&renderBuf, alert); err != nil {
		fmt.Println(err.Error())
		return err
	}

	renderedAlert := renderBuf.String()
	renderedAlert = strings.Trim(renderedAlert, " \n")

	RoomSendMessage(GlobalConfig.Xmpp.Room, alert.Message)

	fmt.Println("Sending alert:", renderedAlert)
	fmt.Println(strings.Repeat("*", 80))

	return nil
}

func SetupHandlers(app *fiber.App) {
	app.Post("/alert", AlertHandler)
}
