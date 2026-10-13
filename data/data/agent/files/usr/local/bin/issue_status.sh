#!/bin/bash

issue_file() {
    printf "/etc/issue.d/%s.issue" "$1"
}

set_issue() {
    local outfile
    outfile="$(issue_file "$1")"
    local tmp
    tmp="$(mktemp -p /etc/issue.d/)"
    {
        printf '\n'
        cat -
        printf '\n'
    } >"${tmp}"
    if ! diff "${tmp}" "${outfile}" >/dev/null 2>&1; then
        mv "${tmp}" "${outfile}"
        agetty --reload
    else
        rm "${tmp}"
    fi
}

clear_issue() {
    local outfile
    outfile="$(issue_file "$1")"
    if [ -f "${outfile}" ]; then
        rm "${outfile}"
        agetty --reload
    fi
}
