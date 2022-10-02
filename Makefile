enc:
	cat plain.txt | go run cmd/koso2/main.go ddddddO > encrypted_by_koso2

dec:
	openssl pkeyutl -decrypt -in encrypted_by_koso2 -inkey ~/.ssh/id_rsa \
    -pkeyopt rsa_padding_mode:oaep -pkeyopt rsa_oaep_md:sha256