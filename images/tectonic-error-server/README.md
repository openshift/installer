# tectonic-error-server

[![Container Repository on Quay](https://quay.io/repository/coreos/tectonic-error-server/status?token=545314b1-7617-4b5a-bd96-c476cf311c4a "Container Repository on Quay")](https://quay.io/repository/coreos/tectonic-error-server)

In case of an error in a request to the Nginx Ingress controller, the body of the response is obtained from the default backend image. The server image in this repo is used as the default backend image for the tectonic Ingress controller. It provides custom error pages depending on the error code returned.

## Description

The `X-Code` value in the header indicates the HTTP error code encountered by the Nginx Ingress controller. The custom error pages to be returned are located in the directory `/web`. If a custom error page is unavailable for an error code, a generic `index.html` page is displayed.

## Usage

Run the `build-docker` script followed by the `push` script to build and push a docker image to the [coreos/tectonic-error-server](https://quay.io/repository/coreos/tectonic-error-server) repo on quay.io.

```
./build-docker
sudo ./push <TAG>
```
