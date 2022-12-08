# Azure Authentication using Client certificates

Azure accepts client certificates as a means of authentication in the service principal and terraform accepts it too. 
As of 4.12, the Installer accepts certificate-based service principals in addition to secret-based service principals.

### Pitfalls
Although the installer can now use the certs to authenticate, CCO does not support this and hence the installer
should create the cluster in manual credentials mode only.

### Prerequisites
- A certificate that is suitable for Azure is created. 
- Register the certificate with Azure as part of App Registrations in the Azure AD.
More information on how to do these steps is below.

## Steps
The installer takes the service principal for authentication and the current requirement is that these fields are
populated.

After [1], we no longer need to pass the client secret and can pass the clientCertificate and clientCertificatePassword fields.

1. Populate the service principal with the following fields. The service principal is by default in ~/.azure/osServicePrincipal.json
-- subscriptionId
-- tenantId
-- clientId
-- clientCertificate (this must be the path to the pfx file generated)
-- clientCertificatePassword (optional)
2. Run openshift-installer

The installer will automatically pick up the values in the sevice principal and switch to certificate based authentication.

## Extras
### Creating a certificate
Azure expects a PEM file certificate for App registrations that are used for authentication. It then expects the PEM certificate and the
key to be combined into a PFX file for authentication requests any application makes. We can use openssl to create these certificates.
There are multiple ways to generate certificate but the key points to remember is that terraform does not accept the latest algorithms 
for generating pfx files and we need to convert them using old algorithms [2].

To generate a certificate, enter the following command using openssl

`openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -sha256 -days 365`

This command will ask a few questions that need to be answered and it generates a PEM certificate and a PEM key file that are valid for one year.

Once this command is done, we need to register this certificate with Azure AD. Navigate to the AD and click on App Registrations.
We can either reuse an existing registration or create a new one in which case, enter the name and make sure the Supported Account Types is set to
the appropriate permission. Redirect URI is optional and need not be entered.

Once the registration is created, click on it, navigate to the Certificates and Secrets section and click on the Upload Certificate to upload the PEM
certificate that we created. This marks the end of the certificate generation section.

We need to now create the pfx file that we would need to authenticate with azure. This can be done with the following command.

`openssl pkcs12 -certpbe PBE-SHA1-3DES -keypbe PBE-SHA1-3DES -export -macalg sha1 -inkey key.pem -in cert.pem -export -out cert.pfx`

This ensures the pfx file is generated with an algorithm that terraform understands. The pfx file is now ready to use for auth and can be set to the
"clientCertificate" key in the osServicePrincipal.json file mentioned above.

### References
[1] - PR for enabling Azure certs auth : https://github.com/openshift/installer/pull/6250
[2] - https://github.com/hashicorp/terraform-provider-azurerm/issues/16228