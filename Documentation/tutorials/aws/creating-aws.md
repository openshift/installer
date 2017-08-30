# Creating an AWS account

You’re well on your way to getting up and running with CoreOS Tectonic. Let’s get started.

In this tutorial you will:

* Create an Amazon Web Services (AWS) account
* Secure your AWS access keys
* Configure automated DNS with Amazon Route 53

## Creating and configuring your AWS account

To create your AWS account:

1. Go to https://aws.amazon.com and click *Create a Free Account*.
2. Enter your payment information using a valid credit card.
3. Complete the identity verification process by answering Amazon’s phone call and entering the PIN show in your browser.
4. Select the support plan that will best serve your business needs.

After successfully creating the account, sign in to the AWS console.

## Configuring AWS Route 53

Route 53 is an Amazon service that allows you to perform DNS management, traffic management, availability monitoring and domain registration. DNS management is the only feature of Route 53 required to install Tectonic.

### Create an AWS Route 53 Hosted Zone

1. From Services, select *Networking & Content Delivery > Route 53*.
2. Select *Hosted zones* from the left pane, and click *Create Hosted Zone*.
3. Enter an existing, registered Domain Name, a Comment (if desired), and select a Type (Public or Private Hosted Zone).

When creating an AWS Route 53 Hosted Zone, enter a domain that you own and can manage.

The Tectonic installation requires a hosted zone domain in which it will create two subdomains; one for the Tectonic console, and one for the Kubernetes API server. This allows Tectonic to access and use the listed domain.

Enter the domain and click *Create*.

AWS provides 4 DNS nameservers for the new zone. The domain or sub-domain must be [configured to use these nameservers][aws-r53-doc]. Visit the domain registrar to add the Route53 NS records.

1. Go to your domain registrar’s website.
2. Go to the DNS settings page and enter the four (4) nameservers Amazon provides you.
3. Save your updated domain settings.

Please note that it may take anywhere from minutes to hours for the changes to take effect.

To verify which nameservers are associated with your domain, use a tool like Dig or nslookup. If no nameservers are returned when you look up your domain, changes may still be pending. Here's an example command:

```bash
$ dig -t ns [example.com]
```

You will know the nameservers were set up correctly when the lookup yields the four provided by AWS.

Next, create an EC2 SSH key pair.

## Creating an EC2 SSH key pair

Both an Access key and an EC2 SSH key pair are required for Tectonic installation.

Before creating the key pair, make sure you are in the correct EC2 region. You can confirm your current region by clicking on the location next to your profile.

Next, configure an SSH key pair from the AWS console.
1. Go to *Services > Compute > EC2*.
2. Confirm that you are in the correct EC2 region by checking the location next to your username in the menu bar.
3. Under *Network & Security*, click *Key Pairs*.
4. Click *Create Key Pair*.
5. Name your key pair, and click *Create*.

Your private key will download automatically, and your key pair name and fingerprint will be listed in the page. Make a note of the Key pair name.

Confirm that RW permissions are correctly set on your .prm file. Running `ls -l` for the file should return `rw-------`. Permissions should be read and write (rw) for owner, and none (---) for all other users.

Finally, create an AWS access key.

## Creating AWS access keys

To use Tectonic with AWS, you must supply Tectonic with a set of security credentials that allow it to authenticate against your AWS account. Use the AWS console to create these credentials.

Create a new Access key ID and Secret access key pair from the AWS console.
1. Select *Services > Security, Identity & Compliance > IAM*.
2. From the left hand pane, click *Users*, then click your username in the list provided.
3. From the Summary page, click the *Security credentials* tab,  and click *Create access key*.
4. Click *Download .csv file*, and save your Secret access key for later.

Both the Access Key ID and its corresponding Secret Access Key are used during the Tectonic installation process.

[**NEXT:** Downloading and installing Tectonic on AWS][installing-tectonic]


[installing-tectonic]: installing-tectonic.md
[aws-r53-doc]: https://docs.aws.amazon.com/Route53/latest/DeveloperGuide/creating-migrating.html
