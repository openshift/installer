output "tnc_cert_pem" {
  value = "${tls_locally_signed_cert.tnc.cert_pem}"
}

output "tnc_key_pem" {
  value = "${tls_private_key.tnc.private_key_pem}"
}

output "id" {
  value = "${sha1("
  ${join(" ",
    list(local_file.tnc_cert.id,
    local_file.tnc_key.id)
    )}
  ")}"
}
