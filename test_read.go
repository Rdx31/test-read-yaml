package read

import "testing"

func TestReadCfg(t *testing.T) {
	cfg := ServiceCfg{}
	LoadCfg("./test_data.yaml", &cfg)

	want := "localhost"
	got := cfg.GRPC.Server.Host
	if want != got {
		t.Errorf("want: %v , got: %v ", want, got)
	}
}
