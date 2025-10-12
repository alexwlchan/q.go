package q

import "testing"

func TestChooseDisplayName(t *testing.T) {
	for _, tt := range []struct {
		label        string
		functionName string
		displayName  string
	}{
		{
			label:        "test-function",
			functionName: "tailscale.com/wgengine/magicsock.TestNetworkDownSendErrors",
			displayName:  "TestNetworkDownSendErrors",
		},
		{
			label:        "test-function-with-parameters",
			functionName: "tailscale.com/tstest/integration.TestOneNodeUpAuth.func1",
			displayName:  "TestOneNodeUpAuth",
		},
	} {
		t.Run(tt.label, func(t *testing.T) {
			if got := chooseDisplayName(tt.functionName); got != tt.displayName {
				t.Fatalf("wanted: %s, got: %s", tt.displayName, got)
			}
		})

	}
}
