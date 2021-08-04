package twsaws

import (
	"fmt"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
	"strings"
	"time"
)

// validationPhase provides a set of utilities to accumulate possible Terraform schema problems into diagnostics
type validationPhase struct {
	problems diag.Diagnostics
	data *schema.ResourceData
}

//hasProblems determines if any performed validations resulted in problems.
func (v *validationPhase) hasProblems() bool  {
	return len(v.problems) > 0
}

//journalError records a validation error for later reporting.  This method is intended to be indirectly called via
//the series of validation functions
func (v *validationPhase) journalError(element string, dataType string, problem error) {
	summary := fmt.Sprintf("%s is not a valid %s", element, dataType)
	details := fmt.Sprintf("%s resulted in %s",element, problem.Error())
	path := cty.IndexStringPath(element)

	v.problems = append(v.problems, diag.Diagnostic{
		Severity:      diag.Error,
		Summary:       summary,
		Detail:        details,
		AttributePath: path,
	})
}

//validateStringOneOf ensures the given schema element consists of a string with one of the given options
func (v *validationPhase) validateStringOneOf(element string, choices []string) string {
	value := v.data.Get(element).(string)

	for _, k := range choices {
		if k == value {
			return k
		}
	}

	summary := fmt.Sprintf("%s has invalid choice", element)
	details := fmt.Sprintf("%s must have one of %#v",element, choices)
	path := cty.IndexStringPath(element)

	v.problems = append(v.problems, diag.Diagnostic{
		Severity:      diag.Error,
		Summary:       summary,
		Detail:        details,
		AttributePath: path,
	})
	return ""
}

//parseExtendedDurationDay parses the given string into a duration, adding the capability to parse days.  A day is
//assumed to be 24 hours regardless of any leap seconds, days, minutes, etc.
func parseExtendedDurationDay(fromValue string, suffix string) (time.Duration, error) {
	totalLength := len(fromValue)
	suffixLength := len(suffix)
	numberPartEnd := totalLength - suffixLength
	numberPart := fromValue[0:numberPartEnd]
	count, err :=  strconv.Atoi(numberPart)
	if err != nil { return 0, err }
	return time.Hour * 24 * time.Duration(count), nil
}

//validateExtendedDuration validates the given schema attribute contains a string which can be parsed into a duration
//format with an extension for days.
func (v *validationPhase) validateExtendedDuration(element string) time.Duration {
	value := v.data.Get(element).(string)

	if strings.HasSuffix(value, "day") {
		amount, err := parseExtendedDurationDay(value, "day")
		if err == nil {
			return amount
		} else {
			v.journalError(element,"extended duration", err)
		}
	}
	if strings.HasSuffix(value, "days") {
		amount, err := parseExtendedDurationDay(value, "days")
		if err == nil {
			return amount
		} else {
			v.journalError(element,"extended duration", err)
		}
	}

	amount, err := time.ParseDuration(value)
	if err != nil {
		v.journalError(element,"extended duration", err)
	}
	return amount
}
