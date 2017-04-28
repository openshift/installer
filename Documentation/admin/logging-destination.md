# Customizing log destination

In order for Fluentd to send your logs to a different destination, you will need to use different Docker image with the correct Fluentd plugin for your destination. Once you have an image, you need to replace the contents of the `output.conf` section in your [fluentd-configmap.yaml][fluentd-config] with the appropriate [match directive][fluentd-match] for your output plugin.

## Prebuilt images

There are currently 4 prebuilt Debian based Docker images in the [quay.io/coreos/fluentd-kubernetes][quay-fluentd-Kubernetes] registry available for various logging destinations:

- [quay.io/coreos/fluentd-kubernetes:v0.12-debian-cloudwatch][quay-fluentd-kubernetes]
- [quay.io/coreos/fluentd-kubernetes:v0.12-debian-logentries][quay-fluentd-kubernetes]
- [quay.io/coreos/fluentd-kubernetes:v0.12-debian-loggly][quay-fluentd-kubernetes]
- [quay.io/coreos/fluentd-kubernetes:v0.12-debian-elasticsearch][quay-fluentd-kubernetes]

**Note**: there are Alpine based images which are automatically published along side the Debian images, but they cannot be used in conjunction with the systemd input plugin, because Alpine has no `libsystemd` package available.

To use one of these images, update the `image` field in your [fluentd-ds.yaml][fluentd-ds] manifest, and update your [fluentd-configmap.yaml][fluentd-config] `output.conf` with the correct match configuration for your configured output plugin.

If you deploy Elasticsearch into your cluster, ensure the hostname and port of the service match the value in the `output.conf` section of your [fluentd-configmap.yaml][fluentd-config].

### Using a different storage destination than Elasticsearch

To change where your logs are sent, change the image in [fluentd-ds.yaml][fluentd-ds] to an image providing the necessary output plugin. The `output.conf` stanza in [fluentd-configmap.yaml][fluentd-config] must also be updated to match the new output plugin.

#### Logentries

To change to the logentries image, replace the line containing `image: quay.io/coreos/fluentd-kubernetes:v0.12-debian-elasticsearch` with `image: quay.io/coreos/fluentd-kubernetes:v0.12-debian-logentries`.

Next, update your [fluentd-configmap.yaml][fluentd-config] `output.conf`:

```
<match **>
  # Plugin specific settings
  type logentries
  config_path /etc/logentries/logentries-tokens.conf

  # Buffer settings
  buffer_chunk_limit 2M
  buffer_queue_limit 32
  flush_interval 10s
  max_retry_wait 30
  disable_retry_limit
  num_threads 8
</match>
```

**Note**: You will need to also modify your `fluentd-ds.yaml` to add a secret/volumeMount for your `logentries-token.conf` referenced in the config above.

#### Cloudwatch

To change to the cloudwatch image, replace the line containing `image: quay.io/coreos/fluentd-kubernetes:v0.12-debian-elasticsearch` with `image: quay.io/coreos/fluentd-kubernetes:v0.12-debian-cloudwatch`. You will also need to create an IAM user and set the `AWS_ACCESS_KEY`, `AWS_SECRET_KEY` and `AWS_REGION` environment variables in your manifest [as documented in the plugin's README](https://github.com/ryotarai/fluent-plugin-cloudwatch-logs#preparation). We recommend you use secrets [as environment variables](https://kubernetes.io/docs/concepts/configuration/secret/#using-secrets-as-environment-variables) to accomplish setting the environment variables securely.

Next, update your [fluentd-configmap.yaml][fluentd-config] `output.conf`:

```
<match **>
  # Plugin specific settings
  type cloudwatch_logs
  log_group_name your-log-group
  log_stream_name your-log-stream
  auto_create_stream true

  # Buffer settings
  buffer_chunk_limit 2M
  buffer_queue_limit 32
  flush_interval 10s
  max_retry_wait 30
  disable_retry_limit
  num_threads 8
</match>
```

#### Loggly

To change to the cloudwatch image, replace the line containing `image: quay.io/coreos/fluentd-kubernetes:v0.12-debian-elasticsearch` with `image: quay.io/coreos/fluentd-kubernetes:v0.12-debian-loggly`.

Next, update your [fluentd-configmap.yaml][fluentd-config] `output.conf` (replace `xxx-xxxx-xxxx-xxxxx-xxxxxxxxxx` with your loggly customer token):

```
<match **>
  # Plugin specific settings
  type loggly_buffered
  loggly_url https://logs-01.loggly.com/bulk/xxx-xxxx-xxxx-xxxxx-xxxxxxxxxx

  # Buffer settings
  buffer_chunk_limit 2M
  buffer_queue_limit 32
  flush_interval 10s
  max_retry_wait 30
  disable_retry_limit
  num_threads 8
</match>
```

### Elasticsearch Authentication

If your Elasticsearch cluster is configured with x-pack or authentication methods, you will need to modify the `output.conf` section of your [fluentd-configmap.yaml][fluentd-config] to set credentials.

Installing the x-pack plugin on your Elasticsearch nodes enables authentication to Elasticsearch by default. The default user is `elastic` and the default password is `changeme`. Modify the configuration to include the `user` and `password` fields like so:

```
<match **>
  type elasticsearch
  log_level info
  include_tag_key true

  # Connection settings
  host elasticsearch.default.svc.cluster.local
  port 9200
  scheme https
  ssl_verify true
  user elastic
  password changeme

  logstash_format true
  template_file /fluentd/etc/elasticsearch-template-es5x.json
  template_name elasticsearch-template-es5x.json

  # Buffer settings
  buffer_chunk_limit 2M
  buffer_queue_limit 32
  flush_interval 5s
  max_retry_wait 30
  disable_retry_limit
  num_threads 8
</match>
```

### AWS Elasticsearch

AWS Elasticsearch isn't officially supported at this time. AWS Elasticsearch functions differently from standard Elasticsearch making the standard Fluentd Elasticsearch plugin not work with it by default. If you are interested in using or contributing support for AWS Elasticsearch, please file an issue on [Github](https://github.com/coreos-inc/tectonic-installer/issues). Alternatively, you could look at one of the many AWS Signing proxies, which may work for your purposes, and allow you to use the default Elasticsearch configuration, pointed at the signing proxy.


[fluentd-ds]: ../files/logging/fluentd-ds.yaml
[fluentd-config]: ../files/logging/fluentd-configmap.yaml
[quay-fluentd-kubernetes]: https://quay.io/repository/coreos/fluentd-kubernetes?tab=tags
[fluentd-docs-output]: http://docs.fluentd.org/v0.12/articles/output-plugin-overview
[fluentd-match]: http://docs.fluentd.org/v0.12/articles/config-file#2-ldquomatchrdquo-tell-fluentd-what-to-do
