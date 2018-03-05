# frozen_string_literal: true

require 'certificate_authority'
require 'fileutils'

# Generates necessary TLS certificates
def generate_tls(path, cluster_name, domain, expiration_date = 365)
  root_ca, kube_ca, aggregator_ca, etcd_ca = generate_ca_certs(expiration_date)
  server_ingress = generate_ingress_certs(root_ca, cluster_name, domain, expiration_date)

  [
    ['ca/root-ca.crt', root_ca.to_pem],
    ['ca/kube-ca.crt', kube_ca.to_pem],
    ['ca/kube-ca.key', kube_ca.key_material.private_key],
    ['ca/aggregator-ca.crt', aggregator_ca.to_pem],
    ['ca/aggregator-ca.key', aggregator_ca.key_material.private_key],
    ['ca/etcd-ca.crt', etcd_ca.to_pem],
    ['ca/etcd-ca.key', etcd_ca.key_material.private_key],

    ['ingress/ca.crt', root_ca.to_pem],
    ['ingress/ingress.key', server_ingress.key_material.private_key],
    ['ingress/ingress.crt', server_ingress.to_pem]
  ].each do |cert_name, contents|
    save_tls_to_file(path, cert_name, contents)
  end
end

def generate_ca_certs(expiration_date = 365)
  root_ca = certificate_authority(1, 'root-ca', 'tectonic', expiration_date)
  kube_ca = intermediate_certificate_authority(root_ca, 2, 'kube-ca', 'bootkube', expiration_date)
  aggregator_ca = intermediate_certificate_authority(root_ca, 3, 'aggregator', 'bootkube', expiration_date)
  etcd_ca = intermediate_certificate_authority(root_ca, 4, 'etcd-ca', 'etcd', expiration_date)

  [root_ca, kube_ca, aggregator_ca, etcd_ca]
end

def generate_ingress_certs(root_ca, cluster_name, domain, expiration_date = 365)
  dns_name = if domain.include?('example.com')
               "tectonic.#{domain}"
             else
               "#{cluster_name}.#{domain}"
             end

  signing_profile_server = {
    'extensions' => {
      'keyUsage' => { 'usage' => %w[critical keyEncipherment digitalSignature] },
      'extendedKeyUsage' => { 'usage' => %w[serverAuth clientAuth] },
      'subjectAltName' => { 'dns_names' => [dns_name] }
    }
  }
  server_ingress = server_certificate(root_ca, dns_name, '', signing_profile_server, expiration_date)
  server_ingress
end

def certificate_authority(serial_number, common_name, organization, expiration_date = 365)
  mem_key = CertificateAuthority::MemoryKeyMaterial.new
  mem_key.generate_key

  root = CertificateAuthority::Certificate.new
  root.subject.common_name = common_name
  root.subject.organization = organization
  root.serial_number.number = serial_number
  root.signing_entity = true
  root.not_after = (Time.now + expiration_date * 86_400).utc # default is 365 days
  root.key_material = mem_key

  ca_profile = {
    'extensions' => {
      'keyUsage' => { 'usage' => %w[critical keyCertSign digitalSignature keyEncipherment] }
    }
  }
  root.sign!(ca_profile)
  root
end

def intermediate_certificate_authority(root, serial_number, common_name, organization, expiration_date = 365)
  mem_key = CertificateAuthority::MemoryKeyMaterial.new
  mem_key.generate_key

  inter = CertificateAuthority::Certificate.new
  inter.subject.common_name = common_name
  inter.subject.organization = organization
  inter.serial_number.number = serial_number
  inter.signing_entity = true
  inter.parent = root
  inter.not_after = (Time.now + expiration_date * 86_400).utc # default is 365 days
  inter.key_material = mem_key

  ca_profile = {
    'extensions' => {
      'keyUsage' => { 'usage' => %w[critical keyCertSign digitalSignature keyEncipherment] }
    }
  }
  inter.sign!(ca_profile)
  inter
end

def server_certificate(root, common_name, organization, signing_profile, expiration_date = 365)
  server = CertificateAuthority::Certificate.new
  server.subject.common_name = common_name
  server.subject.organization = organization
  server.serial_number.number = rand(3..100_000)
  server.parent = root
  server.not_after = (Time.now + expiration_date * 86_400).utc # + 1 day
  server.key_material.generate_key

  server.sign!(signing_profile)
  server
end

def save_tls_to_file(path, cert_name, contents)
  path_to_save = File.join(path, cert_name)

  FileUtils.mkdir_p(File.dirname(path_to_save))
  File.write(path_to_save, contents)
  File.chmod(0o600, path_to_save)
end
