# tectonic-stats-extender

[![Container Repository on Quay](https://quay.io/repository/coreos/tectonic-stats-extender/status?token=17a1091d-1fe9-4da2-af92-416b8a0e5f54 "Container Repository on Quay")](https://quay.io/repository/coreos/tectonic-stats-extender)

tectonic-stats-extender is a basic implementation of a sidecar for [Spartakus](https://github.com/kubernetes-incubator/spartakus) that is customized to report information on a Tectonic license. The extender is in charge of writing extra data you may want to report to a file that can be read by Spartakus.

## Running

The extender binary can be run like:

```sh
./bin/extender --license=/path/to/tectonic/license --output=/path/to/output/file --extensions=extra:data --extensions=to:report
```

The default file generation interval is 1 hour, though this can be customized with the `--period` flag.

## Testing

To test this binary with the example license provided in the tectonic-licensing repo, try the following. First build the tectonic-extender binary:

```sh
./build
```

Then run the following command:

```sh
./bin/extender --output=extensions --license=vendor/github.com/coreos-inc/tectonic-licensing/license/test-license.txt --public-key=vendor/github.com/coreos-inc/tectonic-licensing/license/test-signing-key.pub
```

If this worked, you will see something like the following logged:

```sh
INFO[0000] started stats-extender
INFO[0000] successfully generated extensions
```

You should now see a file called `extensions` with the following contents:

```json
{"accountID":"ACC-FA720BE4-6C55-476A-812C-C4CA6862"}
```

To try setting some custom extensions using the `--extension` flag in addition to the license, try the following:

```sh
./bin/extender --output=extensions --license=vendor/github.com/coreos-inc/tectonic-licensing/license/test-license.txt --public-key=vendor/github.com/coreos-inc/tectonic-licensing/license/test-signing-key.pub --extension=newKey:newValue
```

If this worked, you will see something like the following logged:

```sh
INFO[0000] started stats-extender
INFO[0000] successfully generated extensions
```

You should see the following content in the `extensions` file:

```json
{"accountID":"ACC-FA720BE4-6C55-476A-812C-C4CA6862","newKey":"newValue"}
```
