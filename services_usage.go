package deploystack

import (
	"context"
	"fmt"
	"strings"
	"time"

	"google.golang.org/api/serviceusage/v1"
)

var (
	enabledServices     = make(map[string]bool)
	serviceUsageService *serviceusage.Service
)

func getServiceUsageService(project string) (*serviceusage.Service, error) {
	if serviceUsageService != nil {
		return serviceUsageService, nil
	}

	ctx := context.Background()
	svc, err := serviceusage.NewService(ctx, opts)
	if err != nil {
		return nil, err
	}

	serviceUsageService = svc

	return svc, nil
}

// ServiceEnable enable a service in the selected project so that query calls
// to various lists will work.
func ServiceEnable(project, service string) error {
	if _, ok := enabledServices[service]; ok {
		return nil
	}

	svc, err := getServiceUsageService(project)
	if err != nil {
		return err
	}

	enabled, err := ServiceIsEnabled(project, service)
	if err != nil {
		return err
	}

	if enabled {
		fmt.Printf("Service %s already enabled in project %s: ..\n", service, project)
		enabledServices[service] = true
		return nil
	}

	s := fmt.Sprintf("projects/%s/services/%s", project, service)
	fmt.Printf("Enabling service %s in project %s.\n", service, project)
	op, err := svc.Services.Enable(s, &serviceusage.EnableServiceRequest{}).Do()
	if err != nil {
		return fmt.Errorf("could not enable service: %s", err)
	}

	if !strings.Contains(string(op.Response), "ENABLED") {

		fmt.Printf("Waiting for service to be enabled...")

		for i := 0; i < 60; i++ {
			enabled, err = ServiceIsEnabled(project, service)
			if err != nil {
				return err
			}
			if enabled {
				fmt.Printf("complete.\n")
				enabledServices[service] = true
				return nil
			}
			fmt.Printf(".")
			time.Sleep(1 * time.Second)
		}

	}

	enabledServices[service] = true
	return nil
}

func ServiceIsEnabled(project, service string) (bool, error) {
	svc, err := getServiceUsageService(project)

	s := fmt.Sprintf("projects/%s/services/%s", project, service)
	current, err := svc.Services.Get(s).Do()
	if err != nil {
		return false, err
	}

	if current.State == "ENABLED" {
		return true, nil
	}

	return false, nil
}

func ServiceDisable(project, service string) error {
	svc, err := getServiceUsageService(project)
	if err != nil {
		return err
	}
	s := fmt.Sprintf("projects/%s/services/%s", project, service)
	if _, err := svc.Services.Disable(s, &serviceusage.DisableServiceRequest{}).Do(); err != nil {
		return fmt.Errorf("could not disable service: %s", err)
	}

	return nil
}
