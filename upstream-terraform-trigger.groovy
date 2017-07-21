node {
    build job: 'tectonic-installer/master',
    parameters: [
        string(
            name: 'builder_image',
            value: 'quay.io/coreos/tectonic-builder:v1.36-upstream-terraform'
        )
    ]
}
