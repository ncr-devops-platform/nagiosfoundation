package nagiosformatters

import (
	"fmt"
)

// GreaterFormatNagiosCheck compares a value against thresholds and returns nagios output
func GreaterFormatNagiosCheck(name string, value float64, warning float64, critical float64, metricName string) (string, int) {
	if value > critical {
		return fmt.Sprintf("%s CRITICAL - value = %f | %s=%f", name, value, metricName, value), 2
	}
	if value > warning {
		return fmt.Sprintf("%s WARNING - value = %f | %s=%f", name, value, metricName, value), 1
	}
	return fmt.Sprintf("%s OK - value = %f | %s=%f", name, value, metricName, value), 0
}

// LesserFormatNagiosCheck compares a value against thresholds and returns nagios output
func LesserFormatNagiosCheck(name string, value float64, warning float64, critical float64, metricName string) (string, int) {
	if value < critical {
		return fmt.Sprintf("%s CRITICAL - value = %f | %s=%f", name, value, metricName, value), 2
	}
	if value < warning {
		return fmt.Sprintf("%s WARNING - value = %f | %s=%f", name, value, metricName, value), 1
	}
	return fmt.Sprintf("%s OK - value = %f | %s=%f", name, value, metricName, value), 0
}
