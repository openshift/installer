#!/usr/bin/env bash


count=1
interval=15m

# quick and dirty opts parsing (be careful)
while [[ "$#" -gt 0 ]]; do
  case $1 in
    --repo)
        gh_args="$gh_args $1 $2"
        shift
        ;;
    --repo=*)
        gh_args="$gh_args --repo ${1#*=}"
        ;;
    --interval)
        interval="$2"
        shift
        ;;
    --count)
        count="$2"
        shift
        ;;
    *)
      pr="$1";;
  esac
  shift
done

interval=${interval/h/*60*60+};interval=${interval/m/*60+};interval=${interval/s/+};interval="${interval}0"
interval=$(($interval))

gh pr view --json number,title --template '{{printf "#%v %v\n" .number .title}}'

for (( c=0 ; c<${count} ; c++ )) ; do
  data="$(gh pr checks ${gh_args} ${pr} | grep -v tide)"

  read failed pending passed <<< $(awk 'BEGIN {failed=0; pending=0; passed=0} ; /fail/ {failed++} ; /pending/ {pending++} ;/pass/ {passed++}; END {print failed, pending, passed}' <(printf "%s" "$data"))

  printf "$(date):  "

  if [ $(($failed+$pending)) -eq 0 ] ; then
    printf  "All checks passed.\n"
    break
  fi

  [ $pending -gt 0 ] && printf "%d checks pending. " $pending
  [ $failed -gt 0 ] && printf "%d checks failing. " $failed
  printf "\n"
  awk '/fail|pending/ {printf "%s %s %s\n", $2, $1, $4}' <(printf "%s" "$data") | sed -e 's/fail/Failed/' -e 's/pass/Passed️/' -e 's/pending/Pending️/' | column -t

  if [ $failed -gt 0 ] ; then
    printf "$(date):  /retest $(gh pr comment ${gh_args} ${pr} --body "/retest")\n"
  fi

  sleep $interval

done
