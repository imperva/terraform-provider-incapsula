resource "incapsula_abp_credential" "my_credential" {
  account_id = var.account_id

  # RSA public key (PEM encoded) used to encrypt the returned secret.
  # Decrypt with, e.g.:
  #   terraform output -raw encrypted_secret | base64 -d \
  #     | openssl pkeyutl -decrypt -inkey <your-private-key-file> \
  #         -pkeyopt rsa_padding_mode:oaep -pkeyopt rsa_oaep_md:sha256 \
  #         -pkeyopt rsa_oaep_label:$(echo -n 'abp_credential' | xxd -p)
  rsa_key = file("${path.module}/public_key.pem")
}

output "encrypted_secret" {
  value = incapsula_abp_credential.my_credential.encrypted_secret
}
