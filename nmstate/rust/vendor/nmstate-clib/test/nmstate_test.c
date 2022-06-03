#include <stddef.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>

#include <nmstate.h>

int main(void) {
	int rc = EXIT_SUCCESS;
	char *state = NULL;
	char *err_kind = NULL;
	char *err_msg = NULL;
	char *log = NULL;
	uint32_t flag = NMSTATE_FLAG_KERNEL_ONLY;

	if (nmstate_net_state_retrieve(flag, &state, &log, &err_kind, &err_msg)
	    == NMSTATE_PASS) {
		printf("%s\n", state);
	} else {
		printf("%s: %s\n", err_kind, err_msg);
		rc = EXIT_FAILURE;
	}

	nmstate_cstring_free(state);
	nmstate_cstring_free(err_kind);
	nmstate_cstring_free(err_msg);
	nmstate_cstring_free(log);
	exit(rc);
}
