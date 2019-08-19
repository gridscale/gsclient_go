package main

import (
	"bufio"
	"os"
	"time"

	"github.com/gridscale/gsclient-go"
	log "github.com/sirupsen/logrus"
)

const LocationUuid = "45ed677b-3702-4b36-be2a-a2eab9827950"

func main() {
	uuid := os.Getenv("GRIDSCALE_UUID")
	token := os.Getenv("GRIDSCALE_TOKEN")
	config := gsclient.NewConfiguration(
		"https://api.gridscale.io",
		uuid,
		token,
		true,
	)
	client := gsclient.NewClient(config)
	log.Info("gridscale client configured")

	log.Info("Create PaaS and Security zone: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	//Get template for creating paas
	paasTemplates, err := client.GetPaaSTemplateList()
	if err != nil {
		log.Fatal("Get PaaS templates has failed with error", err)
	}

	//Create security zone
	secZoneRequest := gsclient.PaaSSecurityZoneCreateRequest{
		Name:         "go-client-security-zone",
		LocationUuid: LocationUuid,
	}
	cSCZ, err := client.CreatePaaSSecurityZone(secZoneRequest)
	if err != nil {
		log.Fatal("Create security zone has failed with error", err)
	}
	log.WithFields(log.Fields{
		"securityzone_uuid": cSCZ.ObjectUuid,
	}).Info("Security zone successfully created")

	//Create PaaS service
	paasRequest := gsclient.PaaSServiceCreateRequest{
		Name:                    "go-client-paas",
		PaaSServiceTemplateUuid: paasTemplates[0].Properties.ObjectUuid,
		PaaSSecurityZoneUuid:    cSCZ.ObjectUuid,
	}
	cPaaS, err := client.CreatePaaSService(paasRequest)
	if err != nil {
		log.Fatal("Create PaaS service has failed with error", err)
	}
	log.WithFields(log.Fields{
		"paas_uuid": cPaaS.ObjectUuid,
	}).Info("PaaS service create successfully")

	log.Info("Update PaaS and Security zone: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	//Get a security zone to update
	secZone, err := client.GetPaaSSecurityZone(cSCZ.ObjectUuid)
	if err != nil {
		log.Fatal("Get security zone has failed with error", err)
	}
	secZoneUpdateRequest := gsclient.PaaSSecurityZoneUpdateRequest{
		Name:                 "updated security zone",
		LocationUuid:         secZone.Properties.LocationUuid,
		PaaSSecurityZoneUuid: secZone.Properties.ObjectUuid,
	}
	//Update security zone
	err = client.UpdatePaaSSecurityZone(secZone.Properties.ObjectUuid, secZoneUpdateRequest)
	if err != nil {
		log.Fatal("Update security zone has failed with error", err)
	}
	log.Info("Security Zone successfully updated")

	//Get a PaaS service to update
	paas, err := client.GetPaaSService(cPaaS.ObjectUuid)
	if err != nil {
		log.Fatal("Get PaaS service has failed with error", err)
	}

	//Update PaaS service
	paasUpdateRequest := gsclient.PaaSServiceUpdateRequest{
		Name:           "updated paas",
		Labels:         paas.Properties.Labels,
		Parameters:     paas.Properties.Parameters,
		ResourceLimits: paas.Properties.ResourceLimits,
	}
	err = client.UpdatePaaSService(paas.Properties.ObjectUuid, paasUpdateRequest)
	if err != nil {
		log.Fatal("Update PaaS service has failed with error", err)
	}
	log.Info("PaaS service successfully updated")

	//Clean up
	log.Info("Delete PaaS and Security zone: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	//PaaS has to be deleted before deleting Security zone
	//because we cannot delete security zone when it is in use
	err = client.DeletePaaSService(paas.Properties.ObjectUuid)
	if err != nil {
		log.Fatal("Delete PaaS service has failed with error", err)
	}
	log.Info("PaaS service successfully deleted")

	//Wait until paas deleted successfully
	//it takes around a minute
	time.Sleep(60 * time.Second)
	err = client.DeletePaaSSecurityZone(secZone.Properties.ObjectUuid)
	if err != nil {
		log.Fatal("Delete security zone has failed with error", err)
	}
	log.Info("Security zone successfully deleted")
}
