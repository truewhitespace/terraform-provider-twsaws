package twsaws

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"testing"
	"time"
)

func mockData(t *testing.T, value string) *validationPhase  {
	data := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"test_attribute": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
	}, map[string]interface{}{
		"test_attribute": value,
	})
	validator := &validationPhase{
		data:     data,
	}
	return validator
}

func (v *validationPhase) assertNoErrors(t *testing.T) {
	if v.hasProblems() {
		t.Errorf("Expected no errors, got %#v", v.problems)
	}
}

func TestParseSecondsDuration(t *testing.T)  {
	validator := mockData(t, "1s")
	result := validator.validateExtendedDuration("test_attribute")
	validator.assertNoErrors(t)
	if result != (1 * time.Second) {
		t.Errorf("Expected to parse as 1 second, got %d", result)
	}
}

func TestParseHoursDuration(t *testing.T)  {
	validator := mockData(t, "2h")
	result := validator.validateExtendedDuration("test_attribute")

	validator.assertNoErrors(t)
	if result != (2 * time.Hour) {
		t.Errorf("Expected to parse as 2 hours, got %d", result)
	}
}

func TestParseDayDuration(t *testing.T)  {
	validator := mockData(t, "1day")
	result := validator.validateExtendedDuration("test_attribute")

	validator.assertNoErrors(t)
	if result != (24 * time.Hour) {
		t.Errorf("Expected to parse as 1 day, got %d", result)
	}
}

func TestParseDaysDuration(t *testing.T)  {
	validator := mockData(t, "2days")
	result := validator.validateExtendedDuration("test_attribute")

	validator.assertNoErrors(t)
	if result != (2 * 24 * time.Hour) {
		t.Errorf("Expected to parse as 2 days, got %d", result)
	}
}
