# CoreOS bootimages data

## Scenario: use embedded metadata

Given an existing target cluster,
When building a new ISO using the node-joiner command
Then the node-joiner command must use the embedded stream data to fetch the base RHCOS ISO (when required)