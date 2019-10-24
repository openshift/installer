package systemd

import (
	"encoding/json"
	"testing"
)

var exampleEntry string = `
{
	"__CURSOR" : "s=74ca9654289b49dc904ebc649cc80943;i=11c7;b=d87bb9b0bf4e418badc5b2a5762d823a;m=593e13e9;t=5b3750cd2db69;x=83b2de9fcedfdd27",
	"__REALTIME_TIMESTAMP" : "1604690191244137",
	"__MONOTONIC_TIMESTAMP" : "1497240553",
	"_BOOT_ID" : "d87bb9b0bf4e418badc5b2a5762d823a",
	"_MACHINE_ID" : "ec207e03c25f54405b7e8b2f0280b4e1",
	"PRIORITY" : "6",
	"SYSLOG_FACILITY" : "3",
	"_UID" : "0",
	"_GID" : "0",
	"_CAP_EFFECTIVE" : "3fffffffff",
	"SYSLOG_IDENTIFIER" : "systemd",
	"_TRANSPORT" : "journal",
	"_PID" : "1",
	"_COMM" : "systemd",
	"_EXE" : "/usr/lib/systemd/systemd",
	"_SYSTEMD_CGROUP" : "/init.scope",
	"_SYSTEMD_UNIT" : "init.scope",
	"_SYSTEMD_SLICE" : "-.slice",
	"_HOSTNAME" : "ip-10-0-4-236",
	"_SELINUX_CONTEXT" : "system_u:system_r:init_t:s0",
	"_CMDLINE" : "/usr/lib/systemd/systemd --switched-root --system --deserialize 16",
	"UNIT" : "bootkube.service",
	"CODE_FILE" : "../src/core/service.c",
	"CODE_LINE" : "2147",
	"CODE_FUNC" : "service_enter_restart",
	"MESSAGE_ID" : "5eb03494b6584870a536b337290809b3",
	"INVOCATION_ID" : "d026e6248ad24604b8450541ff5f1eca",
	"MESSAGE" : "bootkube.service: Scheduled restart job, restart counter is at 277.",
	"N_RESTARTS" : "277",
	"_SOURCE_REALTIME_TIMESTAMP" : "1604690191244128"
}
`

func TestEntry(t *testing.T) {
	var entry Entry
	if err := json.Unmarshal([]byte(exampleEntry), &entry); err != nil {
		t.Fatal(err)
	}

	t.Run("InvocationID", func(t *testing.T) {
		expected := "d026e6248ad24604b8450541ff5f1eca"
		if entry.InvocationID != expected {
			t.Fatalf("got %q , expected %q", entry.InvocationID, expected)
		}
	})

	t.Run("Restarts", func(t *testing.T) {
		expected := 277
		if entry.Restarts != expected {
			t.Fatalf("got %d , expected %d", entry.Restarts, expected)
		}
	})

	t.Run("String", func(t *testing.T) {
		expected := "2020-11-06T19:16:31Z ip-10-0-4-236 systemd[1]: bootkube.service: Scheduled restart job, restart counter is at 277."
		actual := entry.String()
		if actual != expected {
			t.Fatalf("got %q , expected %q", actual, expected)
		}
	})

	t.Run("Unit", func(t *testing.T) {
		expected := "bootkube.service"
		if entry.Unit != expected {
			t.Fatalf("got %q , expected %q", entry.Unit, expected)
		}
	})

}
