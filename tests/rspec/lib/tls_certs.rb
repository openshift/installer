# frozen_string_literal: true

require 'certificate_authority'
require 'fileutils'

# Generates necessary TLS certificates
def generate_tls(path, cluster_name, domain, etcd_server_count, expiration_date = 365)
  root_kube, server_kube, client_kube = generate_api_certs(cluster_name, domain, expiration_date)
  server_identity, client_identity = generate_identity_certs(root_kube, expiration_date)
  root_etcd, server_etcd, server_peer_etcd, client_etcd = generate_etcd_certs(cluster_name, domain,
                                                                              etcd_server_count, expiration_date)
  server_ingress = generate_ingress_certs(root_kube, cluster_name, domain, expiration_date)

  [
    ['kube/ca.crt', root_kube.to_pem],
    ['kube/kubelet.key', client_kube.key_material.private_key],
    ['kube/kubelet.crt', client_kube.to_pem],
    ['kube/apiserver.key', server_kube.key_material.private_key],
    ['kube/apiserver.crt', server_kube.to_pem],

    ['identity/identity-client.key', client_identity.key_material.private_key],
    ['identity/identity-client.crt', client_identity.to_pem],
    ['identity/identity-server.key', server_identity.key_material.private_key],
    ['identity/identity-server.crt', server_identity.to_pem],

    ['etcd/etcd-ca.crt', root_etcd.to_pem],
    ['etcd/etcd-client.key', client_etcd.key_material.private_key],
    ['etcd/etcd-client.crt', client_etcd.to_pem],
    ['etcd/etcd-server.key', server_etcd.key_material.private_key],
    ['etcd/etcd-server.crt', server_etcd.to_pem],
    ['etcd/etcd-peer.key', server_peer_etcd.key_material.private_key],
    ['etcd/etcd-peer.crt', server_peer_etcd.to_pem],

    ['ingress/ca.crt', root_kube.to_pem],
    ['ingress/ingress.key', server_ingress.key_material.private_key],
    ['ingress/ingress.crt', server_ingress.to_pem]
  ].each do |cert_name, contents|
    save_tls_to_file(path, cert_name, contents)
  end
end

def generate_api_certs(cluster_name, domain, expiration_date = 365)
  root_kube = certificate_authority('kube-ca', 'bootkube', expiration_date)

  # when in bare-metal tests we have a pre defined dns name (example.com)
  # and we need to set the correct names for the certs we are about to generate
  # same for lines: 139, 215, 230 which is a bit different from the other platforms
  dns_main_node = if domain.include?('example.com')
                    'cluster.example.com'
                  else
                    "#{cluster_name}-api.#{domain}"
                  end

  signing_profile_server = {
    'extensions' => {
      'keyUsage' => { 'usage' => %w[critical keyCertSign digitalSignature keyEncipherment] },
      'extendedKeyUsage' => { 'usage' => %w[serverAuth clientAuth] },
      'subjectAltName' => { 'dns_names' => [dns_main_node, 'kubernetes', 'kubernetes.default',
                                            'kubernetes.default.svc', 'kubernetes.default.svc.cluster.local'],
                            'ips' => ['10.3.0.1'] }
    }
  }
  server_kube = server_certificate(root_kube, 'kube-apiserver', 'kube-master', signing_profile_server, expiration_date)

  signing_profile_client = {
    'extensions' => {
      'keyUsage' => { 'usage' => %w[critical digitalSignature keyEncipherment] },
      'extendedKeyUsage' => {
        'usage' => %w[clientAuth serverAuth]
      }
    }
  }
  client_kube = client_certificate(root_kube, 'kubelet', 'system:masters', signing_profile_client, expiration_date)

  [root_kube, server_kube, client_kube]
end

def generate_identity_certs(root_kube, expiration_date = 365)
  signing_profile_server = {
    'extensions' => {
      'keyUsage' => { 'usage' => [] },
      'extendedKeyUsage' => { 'usage' => ['serverAuth'] }
    }
  }
  server_identity = server_certificate(root_kube, 'tectonic-identity-api.tectonic-system.svc.cluster.local',
                                       '', signing_profile_server, expiration_date)

  signing_profile_client = {
    'extensions' => {
      'extendedKeyUsage' => {
        'usage' => ['clientAuth']
      }
    }
  }
  client_identity = client_certificate(root_kube, 'tectonic-identity-api.tectonic-system.svc.cluster.local',
                                       '', signing_profile_client, expiration_date)

  [server_identity, client_identity]
end

def generate_etcd_certs(cluster_name, domain, etcd_server_count, expiration_date = 365)
  dns_etcd_server, dns_etcd_peer, dns_main_node = generate_etcd_dns(cluster_name, domain, etcd_server_count)

  root_etcd = certificate_authority(dns_main_node, '', expiration_date)

  signing_profile_server_etcd = {
    'extensions' => {
      'keyUsage' => { 'usage' => %w[critical keyEncipherment] },
      'extendedKeyUsage' => { 'usage' => ['serverAuth'] },
      'subjectAltName' => { 'dns_names' => dns_etcd_server,
                            'ips' => ['127.0.0.1', '10.3.0.15', '10.3.0.20'] }
    }
  }
  server_etcd = server_certificate(root_etcd, dns_main_node, '', signing_profile_server_etcd, expiration_date)

  signing_profile_peer_etcd = {
    'extensions' => {
      'keyUsage' => { 'usage' => %w[critical keyEncipherment] },
      'extendedKeyUsage' => { 'usage' => %w[serverAuth clientAuth] },
      'subjectAltName' => { 'dns_names' => dns_etcd_peer,
                            'ips' => ['10.3.0.15', '10.3.0.20'] }
    }
  }
  server_peer_etcd = server_certificate(root_etcd, dns_main_node, '', signing_profile_peer_etcd, expiration_date)

  signing_profile_client = {
    'extensions' => {
      'keyUsage' => { 'usage' => %w[critical digitalSignature keyEncipherment] },
      'extendedKeyUsage' => {
        'usage' => %w[clientAuth serverAuth]
      }
    }
  }
  client_etcd = client_certificate(root_etcd, dns_main_node, '', signing_profile_client, expiration_date)

  [root_etcd, server_etcd, server_peer_etcd, client_etcd]
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

def certificate_authority(common_name, organization, expiration_date = 365)
  mem_key = CertificateAuthority::MemoryKeyMaterial.new
  mem_key.generate_key

  root = CertificateAuthority::Certificate.new
  root.subject.common_name = common_name
  root.subject.organization = organization
  root.serial_number.number = 1
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

def client_certificate(root_ca, common_name, organization, signing_profile, expiration_date = 365)
  client = CertificateAuthority::Certificate.new
  client.subject.common_name = common_name
  client.subject.organization = organization
  client.serial_number.number = 2
  client.parent = root_ca
  client.not_after = (Time.now + expiration_date * 86_400).utc # + 1 day

  client.key_material.generate_key

  client.sign!(signing_profile)
  client
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

def generate_etcd_dns(cluster_name, domain, etcd_server_count)
  dns_etcd_server = generate_etcd_dns_names(cluster_name, domain, etcd_server_count, true)
  dns_etcd_peer = generate_etcd_dns_names(cluster_name, domain, etcd_server_count, false)

  dns_main_node = if domain.include?('azure')
                    "#{cluster_name}-etcd"
                  elsif domain.include?('example.com')
                    "node1.#{domain}"
                  else
                    "#{cluster_name}-etcd.#{domain}"
                  end

  [dns_etcd_server, dns_etcd_peer, dns_main_node]
end

def generate_etcd_dns_names(cluster_name, domain, etcd_server_count, include_localhost)
  dns_etcd = []
  count = etcd_server_count.zero? ? 1 : etcd_server_count
  (0..count - 1).each do |etcd_node|
    dns_etcd_node = if domain.include?('azure')
                      "#{cluster_name}-etcd-#{etcd_node}"
                    elsif domain.include?('example.com')
                      "node1.#{domain}"
                    else
                      "#{cluster_name}-etcd-#{etcd_node}.#{domain}"
                    end
    dns_etcd.push(dns_etcd_node)
  end
  dns_etcd.push('localhost') if include_localhost
  dns_etcd.push('*.kube-etcd.kube-system.svc.cluster.local')
  dns_etcd.push('kube-etcd-client.kube-system.svc.cluster.local')

  dns_etcd
end
