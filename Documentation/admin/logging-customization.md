# Customizing Fluentd

The [fluentd-configmap.yaml][fluentd-config] provided has been designed to be easily customizable. Generally you'll want to avoid modifying anything other than the `fluentd.conf` and `output.conf` sections of the configmap.

## Customizing log destination

The [customizing log destination][customizing-log-destination] document explains how to configure where logs are sent.

## Add custom parsing/filtering logic

To add additional filters or parsers, add them to the `extra.conf` section in the [fluentd-configmap.yaml][fluentd-config]. The `extra.conf` already has a very brief example of how to add an extra field to log entries, and a more detailed example is shown below.

For details on Fluentd post-processing, check out the Fluentd [fliters][fluentd-docs-filter] and [parsers][fluentd-docs-parser] documents.

### Targeting a specific application's logs

Fluentd routes event based on tags. Events flowing through Fluentd can be routed based on the value of the `tag` using `<match>` and `<filter>` directives. The configuration tags events using the following conventions:

- systemd logs have a tag `systemd.<unit-name-here>`
- Kubernetes container logs have a tag `kube.<namespace-name>.<container-name>`
- Kubernetes API audit logs have a tag `kube-apiserver-audit`

The existing configuration already does additional post processing based on some of these tags.

For example, the host's `kubelet.service` log's are parsed by matching on the tag `systemd.kubelet`, and we do the same for parsing the Docker engine's logs using the tag `systemd.docker`. These filters set their `key_name` parameter to `MESSAGE` which is the actual field name for the log message when it originates from journald.

Similarly, we parse the logs of the kube-apiserver, kube-scheduler, and other controller components by performing a wildcard match on the tag: `kube.kube-system.**`. This filter set its `key_name` parameter to `log`, which is the field for log messages originating from Docker.

#### Example: Parsing the guestbook apache logs

The following configuration will parse the `frontend` component's logs from the guestbook example app deployed in the ["Deploy your second app"][second-app] tutorial. To use it, copy and paste the snippet below [fluentd-configmap.yaml][fluentd-config]'s `extra.conf` section (make sure you indent to the correct level).

```
<filter kube.default.php-redis>
  @type parser
  # Fluentd provides a few built-in formats for popular and common formats such as "apache" and "json".
  format apache2
  key_name log
  # Retain the original "log" field after parsing out the data.
  reserve_data true

  # The access logs and error logs are interleaved with each other and have
  # different formats, so ignore parse errors, as they're expected
  suppress_parse_error_log true
</filter>

<filter kube.default.php-redis>
  @type parser
  format apache_error
  key_name log
  reserve_data true

  # The access logs and error logs are interleaved with each other and have
  # different formats, so ignore parse errors, as they're expected
  suppress_parse_error_log true
</filter>
```

Once you've updated your config, you will need to delete and recreate your `fluentd` pods in order for the configuration to take effect:

```
$ kubectl delete pods --namespace logging -l app=fluentd
```

Once your pods have restarted, any new logs being parsed should be using the new configuration.

[fluentd-config]: ../files/logging/fluentd-configmap.yaml
[quay-fluentd-kubernetes]: https://quay.io/repository/coreos/fluentd-kubernetes?tab=tags
[fluentd-match]: http://docs.fluentd.org/v0.12/articles/config-file#2-ldquomatchrdquo-tell-fluentd-what-to-do
[fluentd-docs-filter]: http://docs.fluentd.org/v0.12/articles/filter-plugin-overview
[fluentd-docs-parser]: http://docs.fluentd.org/v0.12/articles/parser-plugin-overview
[second-app]: ../tutorials/second-app.md
[customizing-log-destination]: logging-destination.md
