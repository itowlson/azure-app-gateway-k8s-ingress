package controller

import (
	"github.com/Azure/azure-sdk-for-go/arm/network"
	"github.com/golang/glog"
)

func (lbc *LoadBalancerController) putPublicIP(publicIP network.PublicIPAddress) (*network.PublicIPAddress, error) {
	client := network.NewPublicIPAddressesClientWithBaseURI(lbc.azureAuth.BaseURI, lbc.azureAuth.SubscriptionID)
	client.Authorizer = lbc.azureAuth

	cancel := make(chan struct{})
	rsrcch, errch := client.CreateOrUpdate(lbc.resourceGroup(), *publicIP.Name, publicIP, cancel)
	err := <-errch
	if err != nil {
		return nil, err
	}

	rsrc := <-rsrcch

	rsrcName := safe(rsrc.Name)
	rsrcIPAddress := safe(rsrc.IPAddress)  // TODO: apparently sometimes this can panic with a NPE - implies name is set but properties is nil (because .IPAddress is shorthand for .<properties>.IPAddress)
	glog.V(1).Infof("created or updated %s with IP addr %s", rsrcName, rsrcIPAddress)

	return &rsrc, nil
}

func (lbc *LoadBalancerController) getPublicIP(publicIP network.PublicIPAddress) (*network.PublicIPAddress, error) {
	client := network.NewPublicIPAddressesClientWithBaseURI(lbc.azureAuth.BaseURI, lbc.azureAuth.SubscriptionID)
	client.Authorizer = lbc.azureAuth

	rsrc, err := client.Get(lbc.resourceGroup(), *publicIP.Name, "")
	if err != nil {
		return nil, err
	}

	glog.V(1).Infof("got pip %s with IP addr %s", safe(rsrc.Name), safe(rsrc.IPAddress))

	return &rsrc, nil
}

func (lbc *LoadBalancerController) deletePublicIP(name string) error {
	client := network.NewPublicIPAddressesClientWithBaseURI(lbc.azureAuth.BaseURI, lbc.azureAuth.SubscriptionID)
	client.Authorizer = lbc.azureAuth

	cancel := make(chan struct{})
	_, errch := client.Delete(lbc.resourceGroup(), name, cancel)
	err := <-errch
	if err != nil {
		return err
	}
	return nil
}
