package systemd

import (
	"strings"
	"testing"
)

var data string = `
{
	"__CURSOR" : "s=74ca9654289b49dc904ebc649cc80943;i=11c6;b=d87bb9b0bf4e418badc5b2a5762d823a;m=593e126d;t=5b3750cd2d9ed;x=15a8b3547325a0",
	"__REALTIME_TIMESTAMP" : "1604690191243757",
	"__MONOTONIC_TIMESTAMP" : "1497240173",
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
	"CODE_LINE" : "3532",
	"CODE_FUNC" : "service_dispatch_timer",
	"MESSAGE" : "bootkube.service: Service RestartSec=5s expired, scheduling restart.",
	"INVOCATION_ID" : "d026e6248ad24604b8450541ff5f1eca",
	"_SOURCE_REALTIME_TIMESTAMP" : "1604690191243740"
}
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
{
	"__CURSOR" : "s=74ca9654289b49dc904ebc649cc80943;i=11c8;b=d87bb9b0bf4e418badc5b2a5762d823a;m=593e1555;t=5b3750cd2dcd4;x=11d7ab1a97962d88",
	"__REALTIME_TIMESTAMP" : "1604690191244500",
	"__MONOTONIC_TIMESTAMP" : "1497240917",
	"_BOOT_ID" : "d87bb9b0bf4e418badc5b2a5762d823a",
	"_MACHINE_ID" : "ec207e03c25f54405b7e8b2f0280b4e1",
	"PRIORITY" : "6",
	"SYSLOG_FACILITY" : "3",
	"_UID" : "0",
	"_GID" : "0",
	"_CAP_EFFECTIVE" : "3fffffffff",
	"CODE_FILE" : "../src/core/job.c",
	"CODE_LINE" : "827",
	"CODE_FUNC" : "job_log_status_message",
	"SYSLOG_IDENTIFIER" : "systemd",
	"JOB_RESULT" : "done",
	"_TRANSPORT" : "journal",
	"_PID" : "1",
	"_COMM" : "systemd",
	"_EXE" : "/usr/lib/systemd/systemd",
	"_SYSTEMD_CGROUP" : "/init.scope",
	"_SYSTEMD_UNIT" : "init.scope",
	"_SYSTEMD_SLICE" : "-.slice",
	"_HOSTNAME" : "ip-10-0-4-236",
	"MESSAGE_ID" : "9d1aaa27d60140bd96365438aad20286",
	"_SELINUX_CONTEXT" : "system_u:system_r:init_t:s0",
	"_CMDLINE" : "/usr/lib/systemd/systemd --switched-root --system --deserialize 16",
	"UNIT" : "bootkube.service",
	"MESSAGE" : "Stopped Bootstrap a Kubernetes cluster.",
	"JOB_TYPE" : "restart",
	"INVOCATION_ID" : "d026e6248ad24604b8450541ff5f1eca",
	"_SOURCE_REALTIME_TIMESTAMP" : "1604690191244271"
}
{
	"__CURSOR" : "s=74ca9654289b49dc904ebc649cc80943;i=11c9;b=d87bb9b0bf4e418badc5b2a5762d823a;m=593e1e39;t=5b3750cd2e5b9;x=85e60042f2e8731c",
	"__REALTIME_TIMESTAMP" : "1604690191246777",
	"__MONOTONIC_TIMESTAMP" : "1497243193",
	"_BOOT_ID" : "d87bb9b0bf4e418badc5b2a5762d823a",
	"_MACHINE_ID" : "ec207e03c25f54405b7e8b2f0280b4e1",
	"PRIORITY" : "6",
	"SYSLOG_FACILITY" : "3",
	"_UID" : "0",
	"_GID" : "0",
	"_CAP_EFFECTIVE" : "3fffffffff",
	"CODE_FILE" : "../src/core/job.c",
	"CODE_LINE" : "827",
	"CODE_FUNC" : "job_log_status_message",
	"SYSLOG_IDENTIFIER" : "systemd",
	"JOB_TYPE" : "start",
	"JOB_RESULT" : "done",
	"MESSAGE_ID" : "39f53479d3a045ac8e11786248231fbf",
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
	"MESSAGE" : "Started Bootstrap a Kubernetes cluster.",
	"UNIT" : "bootkube.service",
	"INVOCATION_ID" : "76e11720db5940948eb008a8eab2d5c0",
	"_SOURCE_REALTIME_TIMESTAMP" : "1604690191246053"
}
{
	"__CURSOR" : "s=74ca9654289b49dc904ebc649cc80943;i=11ca;b=d87bb9b0bf4e418badc5b2a5762d823a;m=593f710c;t=5b3750cd4388b;x=d0ab83b4af8fb625",
	"__REALTIME_TIMESTAMP" : "1604690191333515",
	"__MONOTONIC_TIMESTAMP" : "1497329932",
	"_BOOT_ID" : "d87bb9b0bf4e418badc5b2a5762d823a",
	"_MACHINE_ID" : "ec207e03c25f54405b7e8b2f0280b4e1",
	"PRIORITY" : "6",
	"SYSLOG_FACILITY" : "3",
	"_UID" : "0",
	"_GID" : "0",
	"_SYSTEMD_SLICE" : "system.slice",
	"_TRANSPORT" : "stdout",
	"_EXE" : "/usr/bin/bash",
	"_CAP_EFFECTIVE" : "3fffffffff",
	"_HOSTNAME" : "ip-10-0-4-236",
	"_SELINUX_CONTEXT" : "system_u:system_r:init_t:s0",
	"_COMM" : "bash",
	"SYSLOG_IDENTIFIER" : "bootkube.sh",
	"MESSAGE" : "/usr/local/bin/bootkube.sh: line 6: i-am-a-command-that-does-not-exist: command not found",
	"_CMDLINE" : "bash /usr/local/bin/bootkube.sh",
	"_SYSTEMD_UNIT" : "bootkube.service",
	"_SYSTEMD_CGROUP" : "/system.slice/bootkube.service",
	"_STREAM_ID" : "d266c2bb271448bb92d956fa4e036391",
	"_PID" : "15907",
	"_SYSTEMD_INVOCATION_ID" : "76e11720db5940948eb008a8eab2d5c0"
}
{
	"__CURSOR" : "s=74ca9654289b49dc904ebc649cc80943;i=11cb;b=d87bb9b0bf4e418badc5b2a5762d823a;m=593f7377;t=5b3750cd43af7;x=72a6951ba697a82d",
	"__REALTIME_TIMESTAMP" : "1604690191334135",
	"__MONOTONIC_TIMESTAMP" : "1497330551",
	"_BOOT_ID" : "d87bb9b0bf4e418badc5b2a5762d823a",
	"PRIORITY" : "5",
	"_MACHINE_ID" : "ec207e03c25f54405b7e8b2f0280b4e1",
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
	"CODE_LINE" : "3216",
	"CODE_FUNC" : "service_sigchld_event",
	"MESSAGE" : "bootkube.service: Main process exited, code=exited, status=127/n/a",
	"EXIT_CODE" : "exited",
	"EXIT_STATUS" : "127",
	"INVOCATION_ID" : "76e11720db5940948eb008a8eab2d5c0",
	"_SOURCE_REALTIME_TIMESTAMP" : "1604690191334126"
}
{
	"__CURSOR" : "s=74ca9654289b49dc904ebc649cc80943;i=11cc;b=d87bb9b0bf4e418badc5b2a5762d823a;m=593f7432;t=5b3750cd43bb2;x=e418aea7132d34e0",
	"__REALTIME_TIMESTAMP" : "1604690191334322",
	"__MONOTONIC_TIMESTAMP" : "1497330738",
	"_BOOT_ID" : "d87bb9b0bf4e418badc5b2a5762d823a",
	"_MACHINE_ID" : "ec207e03c25f54405b7e8b2f0280b4e1",
	"PRIORITY" : "4",
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
	"CODE_LINE" : "1683",
	"_HOSTNAME" : "ip-10-0-4-236",
	"_SELINUX_CONTEXT" : "system_u:system_r:init_t:s0",
	"_CMDLINE" : "/usr/lib/systemd/systemd --switched-root --system --deserialize 16",
	"UNIT" : "bootkube.service",
	"CODE_FILE" : "../src/core/service.c",
	"CODE_FUNC" : "service_enter_dead",
	"MESSAGE" : "bootkube.service: Failed with result 'exit-code'.",
	"INVOCATION_ID" : "76e11720db5940948eb008a8eab2d5c0",
	"_SOURCE_REALTIME_TIMESTAMP" : "1604690191334314"
}
{
	"__CURSOR" : "s=74ca9654289b49dc904ebc649cc80943;i=11cd;b=d87bb9b0bf4e418badc5b2a5762d823a;m=598e2e3e;t=5b3750d22f5be;x=9492525567447af8",
	"__REALTIME_TIMESTAMP" : "1604690196493758",
	"__MONOTONIC_TIMESTAMP" : "1502490174",
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
	"CODE_LINE" : "3532",
	"CODE_FUNC" : "service_dispatch_timer",
	"MESSAGE" : "bootkube.service: Service RestartSec=5s expired, scheduling restart.",
	"INVOCATION_ID" : "76e11720db5940948eb008a8eab2d5c0",
	"_SOURCE_REALTIME_TIMESTAMP" : "1604690196493742"
}
{
	"__CURSOR" : "s=74ca9654289b49dc904ebc649cc80943;i=11ce;b=d87bb9b0bf4e418badc5b2a5762d823a;m=598e2ff5;t=5b3750d22f775;x=513a85f2f6b632aa",
	"__REALTIME_TIMESTAMP" : "1604690196494197",
	"__MONOTONIC_TIMESTAMP" : "1502490613",
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
	"INVOCATION_ID" : "76e11720db5940948eb008a8eab2d5c0",
	"MESSAGE" : "bootkube.service: Scheduled restart job, restart counter is at 278.",
	"N_RESTARTS" : "278",
	"_SOURCE_REALTIME_TIMESTAMP" : "1604690196494129"
}
{
	"__CURSOR" : "s=74ca9654289b49dc904ebc649cc80943;i=11cf;b=d87bb9b0bf4e418badc5b2a5762d823a;m=598e3156;t=5b3750d22f8d5;x=ddc8f54f297e172b",
	"__REALTIME_TIMESTAMP" : "1604690196494549",
	"__MONOTONIC_TIMESTAMP" : "1502490966",
	"_BOOT_ID" : "d87bb9b0bf4e418badc5b2a5762d823a",
	"_MACHINE_ID" : "ec207e03c25f54405b7e8b2f0280b4e1",
	"PRIORITY" : "6",
	"SYSLOG_FACILITY" : "3",
	"_UID" : "0",
	"_GID" : "0",
	"_CAP_EFFECTIVE" : "3fffffffff",
	"CODE_FILE" : "../src/core/job.c",
	"CODE_LINE" : "827",
	"CODE_FUNC" : "job_log_status_message",
	"SYSLOG_IDENTIFIER" : "systemd",
	"JOB_RESULT" : "done",
	"_TRANSPORT" : "journal",
	"_PID" : "1",
	"_COMM" : "systemd",
	"_EXE" : "/usr/lib/systemd/systemd",
	"_SYSTEMD_CGROUP" : "/init.scope",
	"_SYSTEMD_UNIT" : "init.scope",
	"_SYSTEMD_SLICE" : "-.slice",
	"_HOSTNAME" : "ip-10-0-4-236",
	"MESSAGE_ID" : "9d1aaa27d60140bd96365438aad20286",
	"_SELINUX_CONTEXT" : "system_u:system_r:init_t:s0",
	"_CMDLINE" : "/usr/lib/systemd/systemd --switched-root --system --deserialize 16",
	"UNIT" : "bootkube.service",
	"MESSAGE" : "Stopped Bootstrap a Kubernetes cluster.",
	"JOB_TYPE" : "restart",
	"INVOCATION_ID" : "76e11720db5940948eb008a8eab2d5c0",
	"_SOURCE_REALTIME_TIMESTAMP" : "1604690196494269"
}
{
	"__CURSOR" : "s=74ca9654289b49dc904ebc649cc80943;i=11d0;b=d87bb9b0bf4e418badc5b2a5762d823a;m=598e39dd;t=5b3750d23015d;x=c985b3968c937b83",
	"__REALTIME_TIMESTAMP" : "1604690196496733",
	"__MONOTONIC_TIMESTAMP" : "1502493149",
	"_BOOT_ID" : "d87bb9b0bf4e418badc5b2a5762d823a",
	"_MACHINE_ID" : "ec207e03c25f54405b7e8b2f0280b4e1",
	"PRIORITY" : "6",
	"SYSLOG_FACILITY" : "3",
	"_UID" : "0",
	"_GID" : "0",
	"_CAP_EFFECTIVE" : "3fffffffff",
	"CODE_FILE" : "../src/core/job.c",
	"CODE_LINE" : "827",
	"CODE_FUNC" : "job_log_status_message",
	"SYSLOG_IDENTIFIER" : "systemd",
	"JOB_TYPE" : "start",
	"JOB_RESULT" : "done",
	"MESSAGE_ID" : "39f53479d3a045ac8e11786248231fbf",
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
	"MESSAGE" : "Started Bootstrap a Kubernetes cluster.",
	"UNIT" : "bootkube.service",
	"INVOCATION_ID" : "65231f59bdf84a7497c30ac4a09e6f2c",
	"_SOURCE_REALTIME_TIMESTAMP" : "1604690196496016"
}
{
	"__CURSOR" : "s=74ca9654289b49dc904ebc649cc80943;i=11d1;b=d87bb9b0bf4e418badc5b2a5762d823a;m=598f84de;t=5b3750d244c5d;x=2709ebb1da03d4f2",
	"__REALTIME_TIMESTAMP" : "1604690196581469",
	"__MONOTONIC_TIMESTAMP" : "1502577886",
	"_BOOT_ID" : "d87bb9b0bf4e418badc5b2a5762d823a",
	"_MACHINE_ID" : "ec207e03c25f54405b7e8b2f0280b4e1",
	"PRIORITY" : "6",
	"SYSLOG_FACILITY" : "3",
	"_UID" : "0",
	"_GID" : "0",
	"_SYSTEMD_SLICE" : "system.slice",
	"_TRANSPORT" : "stdout",
	"_EXE" : "/usr/bin/bash",
	"_CAP_EFFECTIVE" : "3fffffffff",
	"_HOSTNAME" : "ip-10-0-4-236",
	"_SELINUX_CONTEXT" : "system_u:system_r:init_t:s0",
	"_COMM" : "bash",
	"SYSLOG_IDENTIFIER" : "bootkube.sh",
	"MESSAGE" : "/usr/local/bin/bootkube.sh: line 6: i-am-a-command-that-does-not-exist: command not found",
	"_CMDLINE" : "bash /usr/local/bin/bootkube.sh",
	"_SYSTEMD_UNIT" : "bootkube.service",
	"_SYSTEMD_CGROUP" : "/system.slice/bootkube.service",
	"_STREAM_ID" : "01685a93fb9a41ee99699f6f450ef3bd",
	"_PID" : "15950",
	"_SYSTEMD_INVOCATION_ID" : "65231f59bdf84a7497c30ac4a09e6f2c"
}
{
	"__CURSOR" : "s=74ca9654289b49dc904ebc649cc80943;i=11d2;b=d87bb9b0bf4e418badc5b2a5762d823a;m=598f8757;t=5b3750d244ed7;x=25b5b1f488146de1",
	"__REALTIME_TIMESTAMP" : "1604690196582103",
	"__MONOTONIC_TIMESTAMP" : "1502578519",
	"_BOOT_ID" : "d87bb9b0bf4e418badc5b2a5762d823a",
	"PRIORITY" : "5",
	"_MACHINE_ID" : "ec207e03c25f54405b7e8b2f0280b4e1",
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
	"CODE_LINE" : "3216",
	"CODE_FUNC" : "service_sigchld_event",
	"MESSAGE" : "bootkube.service: Main process exited, code=exited, status=127/n/a",
	"EXIT_CODE" : "exited",
	"EXIT_STATUS" : "127",
	"INVOCATION_ID" : "65231f59bdf84a7497c30ac4a09e6f2c",
	"_SOURCE_REALTIME_TIMESTAMP" : "1604690196582095"
}
{
	"__CURSOR" : "s=74ca9654289b49dc904ebc649cc80943;i=11d3;b=d87bb9b0bf4e418badc5b2a5762d823a;m=598f8835;t=5b3750d244fb5;x=2a021d7682aad7f",
	"__REALTIME_TIMESTAMP" : "1604690196582325",
	"__MONOTONIC_TIMESTAMP" : "1502578741",
	"_BOOT_ID" : "d87bb9b0bf4e418badc5b2a5762d823a",
	"_MACHINE_ID" : "ec207e03c25f54405b7e8b2f0280b4e1",
	"PRIORITY" : "4",
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
	"CODE_LINE" : "1683",
	"_HOSTNAME" : "ip-10-0-4-236",
	"_SELINUX_CONTEXT" : "system_u:system_r:init_t:s0",
	"_CMDLINE" : "/usr/lib/systemd/systemd --switched-root --system --deserialize 16",
	"UNIT" : "bootkube.service",
	"CODE_FILE" : "../src/core/service.c",
	"CODE_FUNC" : "service_enter_dead",
	"MESSAGE" : "bootkube.service: Failed with result 'exit-code'.",
	"INVOCATION_ID" : "65231f59bdf84a7497c30ac4a09e6f2c",
	"_SOURCE_REALTIME_TIMESTAMP" : "1604690196582317"
}
`

func TestLog(t *testing.T) {
	log, err := NewLog(strings.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Restarts", func(t *testing.T) {
		expected := 278
		actual := log.Restarts("bootkube.service")
		if actual != expected {
			t.Fatalf("got %d , expected %d", actual, expected)
		}
	})

	t.Run("Format", func(t *testing.T) {
		t.Run("all invocations", func(t *testing.T) {
			expected := []string{
				"2020-11-06T19:16:31Z ip-10-0-4-236 systemd[1]: bootkube.service: Service RestartSec=5s expired, scheduling restart.",
				"2020-11-06T19:16:31Z ip-10-0-4-236 systemd[1]: bootkube.service: Scheduled restart job, restart counter is at 277.",
				"2020-11-06T19:16:31Z ip-10-0-4-236 systemd[1]: Stopped Bootstrap a Kubernetes cluster.",
				"2020-11-06T19:16:31Z ip-10-0-4-236 systemd[1]: Started Bootstrap a Kubernetes cluster.",
				"2020-11-06T19:16:31Z ip-10-0-4-236 bootkube.sh[15907]: /usr/local/bin/bootkube.sh: line 6: i-am-a-command-that-does-not-exist: command not found",
				"2020-11-06T19:16:31Z ip-10-0-4-236 systemd[1]: bootkube.service: Main process exited, code=exited, status=127/n/a",
				"2020-11-06T19:16:31Z ip-10-0-4-236 systemd[1]: bootkube.service: Failed with result 'exit-code'.",
				"2020-11-06T19:16:36Z ip-10-0-4-236 systemd[1]: bootkube.service: Service RestartSec=5s expired, scheduling restart.",
				"2020-11-06T19:16:36Z ip-10-0-4-236 systemd[1]: bootkube.service: Scheduled restart job, restart counter is at 278.",
				"2020-11-06T19:16:36Z ip-10-0-4-236 systemd[1]: Stopped Bootstrap a Kubernetes cluster.",
				"2020-11-06T19:16:36Z ip-10-0-4-236 systemd[1]: Started Bootstrap a Kubernetes cluster.",
				"2020-11-06T19:16:36Z ip-10-0-4-236 bootkube.sh[15950]: /usr/local/bin/bootkube.sh: line 6: i-am-a-command-that-does-not-exist: command not found",
				"2020-11-06T19:16:36Z ip-10-0-4-236 systemd[1]: bootkube.service: Main process exited, code=exited, status=127/n/a",
				"2020-11-06T19:16:36Z ip-10-0-4-236 systemd[1]: bootkube.service: Failed with result 'exit-code'.",
			}
			actual := log.Format("bootkube.service", 0)
			assertStringArraysEqual(t, actual, expected)
		})

		t.Run("positive invocations", func(t *testing.T) {
			expected := []string{
				"2020-11-06T19:16:31Z ip-10-0-4-236 systemd[1]: bootkube.service: Service RestartSec=5s expired, scheduling restart.",
				"2020-11-06T19:16:31Z ip-10-0-4-236 systemd[1]: bootkube.service: Scheduled restart job, restart counter is at 277.",
				"2020-11-06T19:16:31Z ip-10-0-4-236 systemd[1]: Stopped Bootstrap a Kubernetes cluster.",
			}
			actual := log.Format("bootkube.service", 1)
			assertStringArraysEqual(t, actual, expected)
		})

		t.Run("negative invocations", func(t *testing.T) {
			expected := []string{
				"2020-11-06T19:16:36Z ip-10-0-4-236 systemd[1]: Started Bootstrap a Kubernetes cluster.",
				"2020-11-06T19:16:36Z ip-10-0-4-236 bootkube.sh[15950]: /usr/local/bin/bootkube.sh: line 6: i-am-a-command-that-does-not-exist: command not found",
				"2020-11-06T19:16:36Z ip-10-0-4-236 systemd[1]: bootkube.service: Main process exited, code=exited, status=127/n/a",
				"2020-11-06T19:16:36Z ip-10-0-4-236 systemd[1]: bootkube.service: Failed with result 'exit-code'.",
			}
			actual := log.Format("bootkube.service", -1)
			assertStringArraysEqual(t, actual, expected)
		})
	})
}

func assertStringArraysEqual(t *testing.T, actual, expected []string) {
	for i, e := range expected {
		if i >= len(actual) {
			t.Fatalf("no actual line %d, expected %q", i, e)
		} else if actual[i] != e {
			t.Fatalf("line %d: got %q , expected %q", i, actual[i], e)
		}
	}
	if len(actual) > len(expected) {
		t.Fatalf("unexpected line %d: %q", len(expected), actual[len(expected)])
	}
}
