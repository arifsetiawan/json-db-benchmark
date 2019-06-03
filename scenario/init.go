package scenario

import (
	"math/rand"
	"time"
)

// Definition ...
type Definition struct {
	ID       string `json:"id"`
	TenantID string `json:"tenant_id"`
}

var (
	counter          uint64
	referencesApple  []string
	referencesGoogle []string
	version          []string
	tenantID         []string

	instanceStatus []string
)

// RandomInt ...
func RandomInt(min, max int) int {
	return rand.Intn(max-min) + min
}

func init() {
	rand.Seed(time.Now().UnixNano())

	tenantID = []string{
		"apple",
		"google",
	}

	referencesApple = []string{
		"files:bpmn:aa3e4140-ffa3-4fd7-bf03-a3d7c494736e",
		"files:bpmn:f257ce15-e2c7-4bd4-be07-931b036bb662",
		"files:bpmn:bd903081-fcf8-4828-924a-499a27100588",
		"files:bpmn:98143ee7-4b85-4760-bc6a-80bf3b412786",
	}

	referencesGoogle = []string{
		"files:bpmn:3a953454-510c-4ee0-a50a-c87facd0fdc2",
		"files:bpmn:ecbab1bc-0323-4c81-a82f-13f73be96a42",
		"files:bpmn:90e88a6b-3e3c-4e29-be78-897189a8e92e",
		"files:bpmn:eb5383f0-b25a-4d1c-ab3f-69b9fb421daf",
	}

	version = []string{
		"0.1.0",
		"0.2.0",
		"0.3.0",
	}

	instanceStatus = []string{
		"completed",
		"failed",
		"running",
	}

}
