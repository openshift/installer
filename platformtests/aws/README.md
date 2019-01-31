# AWS Tests

This directory contains test suites checking AWS-specific assumptions.
Run with:

```console
$ AWS_PROFILE=your-profile go test .
```

or similar (it needs access to [your AWS credentials][credentials]).

[credentials]: https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html
