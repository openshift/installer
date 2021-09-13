#!/bin/bash

# Add badges:
cat > ./README.md << 'EOF'
[![Build Status](https://travis-ci.org/PromonLogicalis/asn1.svg?branch=master)](https://travis-ci.org/PromonLogicalis/asn1) [![Go Report Card](https://goreportcard.com/badge/github.com/PromonLogicalis/asn1)](https://goreportcard.com/report/github.com/PromonLogicalis/asn1) [![GoDoc](https://godoc.org/github.com/PromonLogicalis/asn1?status.svg)](https://godoc.org/github.com/PromonLogicalis/asn1)
EOF

# Add Go doc
godocdown ./ >> ./README.md
