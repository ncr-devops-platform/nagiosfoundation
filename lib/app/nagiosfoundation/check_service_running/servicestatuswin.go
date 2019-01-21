// +build windows

package nagiosfoundation

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/StackExchange/wmi"
	"github.com/golang/glog"
)

type Win32_Service struct {
	Name string
}

func CheckServiceRunning() {
	var serviceName = flag.String("service_name", "", "the name of the service to check")
	flag.Parse()

	var dst []Win32_Service
	whereClause := fmt.Sprintf("where Name ='%v' AND State LIKE 'Running'", *serviceName)
	glog.Infof("Where clause: %v", whereClause)
	query := wmi.CreateQuery(&dst, whereClause)
	glog.Infof("Service Query: %s", query)
	err := wmi.Query(query, &dst)
	if err != nil {
		log.Fatal(err)
	}
	var msg string
	var retcode int
	if len(dst) <= 0 {
		msg = fmt.Sprintf("CheckServiceRunning CRITICAL - %s not in a running state", *serviceName)
		retcode = 2
	} else {
		msg = fmt.Sprintf("CheckServiceRunning OK- %s in a running state", *serviceName)
		retcode = 0
	}
	fmt.Println(msg)
	os.Exit(retcode)
}
