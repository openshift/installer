'use strict';

const https = require('https');
const fs = require('fs');

const PORT = process.env.PORT || 6444;

// const keysFilenames = process.env.KEY || 'kube-apiserver-service-network-server.key,kube-apiserver-localhost-server.key,kube-apiserver-lb-server.key,kube-apiserver-internal-lb-server.key';
// const certsFilenames = process.env.CERT || 'kube-apiserver-service-network-server.crt,kube-apiserver-localhost-server.crt,kube-apiserver-lb-server.crt,kube-apiserver-internal-lb-server.crt';
const keysFilenames = process.env.KEY || 'kube-apiserver-lb-server.key';
const certsFilenames = process.env.CERT || 'kube-apiserver-lb-server.crt';
const caBundleFilename = process.env.CA_BUNDLE || 'aggregator-ca-bundle.crt';
const caBundle = fs.readFileSync(caBundleFilename);

function readFiles(filenames) {
    filenames.split(",").map(f => fs.readFileSync(f.trim()));
}

const options = {
    key: readFiles(keysFilenames),
    cert: readFiles(certsFilenames),
};

const server = https.createServer(options, (req, res) => {
    console.log(`Processing request from ${req.socket.remoteAddress}:${req.socket.remotePort} for ${req.url}`);
    if (req.url.startsWith("/bootstrap")) {
        onRequestForBootstrapGather(req, res);
    } else if  (req.url.startsWith("/readyz")) {
        onRequestForReadyz(req, res);
    } else {
        onRequestForAPIServer(req, res);
    }
});

server.listen(PORT, () => {
    console.log("Listening on " + PORT);
});

function onRequestForAPIServer(req, res) {
    var connector = https.request({
        host: "localhost",
        port: 6443,
        path: req.url,
        method: req.method,
        headers: req.headers,
        ca: caBundle,
    }, (resp) => {
        console.log("Forwarding request for " + req.url);
        res.statusCode = resp.statusCode;
        new Promise((resolve, reject) => {
            const chunks = [];
            resp.on('data', chunk => chunks.push(chunk));
            resp.on('error', reject);
            resp.on('end', () => {
                resolve(Buffer.concat(chunks).toString('utf8'));
            });
        }).then( (resp) => {
            res.end(resp);
        }).catch( (error) => {
            console.error(error);
            res.statusCode = 500;
            res.end(error);
        })
    });

    connector.on('error', (e) => {
        console.error(`error connecting to API server: ${e.message}`);
        res.statusCode = 400;
        res.end(e.toString());
    });

    req.pipe(connector);
}

function onRequestForReadyz(req, res) {
    res.end("ready");
}

function onRequestForBootstrapGather(req, res) {
    res.end("Connected for bootstrap gather logs");
}