package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/pivotal-cf/brokerapi"
	"github.com/pivotal-golang/lager"
)

var services = []brokerapi.Service{
	{
		ID:            "f614fcc2-3cb2-4400-aa93-87714417f2cf",
		Name:          "MySQL",
		Description:   "MySQL Service",
		Bindable:      true,
		Tags:          []string{"SQL", "Database"},
		PlanUpdatable: true,
		Plans: []brokerapi.ServicePlan{
			{
				ID:          "476217c3-110c-4b88-b43a-74d5b91b8643",
				Name:        "standalone",
				Description: "Standalone MySQL Service",
				Free:        func() *bool { b := false; return &b }(),
			},
			{
				ID:          "726aa087-5a18-4a04-b449-5d455cecb29b",
				Name:        "shared",
				Description: "Shared MySQL Service",
				Free:        func() *bool { b := true; return &b }(),
			},
		},
	},
}

type credentials struct {
	Uri      string `json:"uri"`
	Hostname string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name,omitempty"`
}

type serviceInstance struct {
	ID      string
	details brokerapi.ProvisionDetails
	spec    brokerapi.ProvisionedServiceSpec
}

var instances = map[string]serviceInstance{}

type myServiceBroker struct {
}

func (myBroker *myServiceBroker) Services() []brokerapi.Service {
	return services
}

func (myBroker *myServiceBroker) Provision(
	instanceID string,
	details brokerapi.ProvisionDetails,
	asyncAllowed bool,
) (brokerapi.ProvisionedServiceSpec, error) {
	if _, present := instances[instanceID]; present {
		log.Println("Instance " + instanceID + " already exists.")
		return brokerapi.ProvisionedServiceSpec{}, brokerapi.ErrInstanceAlreadyExists
	}

	instance := serviceInstance{
		ID:      instanceID,
		details: details,
		spec: brokerapi.ProvisionedServiceSpec{
			IsAsync:      true,
			DashboardURL: "",
			Credentials: credentials{
				Uri:      "user:pass@tcp(mysql.service.com:3306)/database",
				Hostname: "mysql.service.com",
				Port:     "3306",
				Username: "user",
				Password: "pass",
				Name:     "database",
			},
		},
	}
	instances[instanceID] = instance

	log.Println("Instance " + instanceID + " created successfully.")

	return instance.spec, nil
}

func (myBroker *myServiceBroker) LastOperation(instanceID string) (brokerapi.LastOperation, error) {
	log.Println("Successful query last operation for instance " + instanceID)

	return brokerapi.LastOperation{
		State:       brokerapi.Succeeded,
		Description: "Succeeded",
	}, nil
}

func (myBroker *myServiceBroker) Deprovision(instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.IsAsync, error) {

	if _, present := instances[instanceID]; !present {
		log.Println("Instance " + instanceID + " doesn't exists.")
		return brokerapi.IsAsync(false), brokerapi.ErrInstanceDoesNotExist
	}

	delete(instances, instanceID)

	log.Println("Instance " + instanceID + " is deleted successfully.")

	return false, nil
}

func (myBroker *myServiceBroker) Bind(instanceID, bindingID string, details brokerapi.BindDetails) (brokerapi.Binding, error) {

	if instance, present := instances[instanceID]; !present {
		log.Println("Instance " + instanceID + " doesn't exists.")
		return brokerapi.Binding{}, brokerapi.ErrInstanceDoesNotExist
	} else {

		var myBinding = brokerapi.Binding{
			Credentials: instance.spec.Credentials,
		}

		log.Println("Instance " + instanceID + " bind " + bindingID + " done.")
		return myBinding, nil
	}
}

func (myBroker *myServiceBroker) Unbind(instanceID, bindingID string, details brokerapi.UnbindDetails) error {

	if instance, present := instances[instanceID]; !present {
		log.Println("Instance " + instanceID + " doesn't exists.")
		return brokerapi.ErrInstanceDoesNotExist
	} else {
		_ = instance
		log.Println("Instance " + instanceID + " unind " + bindingID + " done.")
		return nil
	}
}

func (myBroker *myServiceBroker) Update(instanceID string, details brokerapi.UpdateDetails, asyncAllowed bool) (brokerapi.IsAsync, error) {
	return brokerapi.IsAsync(false), brokerapi.ErrPlanChangeNotSupported
}

func main() {
	logger := lager.NewLogger("example")
	logger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.INFO)) //默认日志级别

	brokerAPI := brokerapi.New(&myServiceBroker{}, logger,
		brokerapi.BrokerCredentials{
			Username: "user",
			Password: "pass",
		})
	http.Handle("/", brokerAPI)
	addr := fmt.Sprintf(":%v", *flag.Int("port", 33333, "server port"))
	fmt.Println("Server started: http://localhost" + addr)
	fmt.Println(http.ListenAndServe(addr, nil))
}
