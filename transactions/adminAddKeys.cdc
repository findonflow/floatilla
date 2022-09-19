
transaction(number:Int, key: String) {
	prepare(signer: AuthAccount) {

		let publicKey=key.decodeHex()
		let pk = PublicKey( publicKey: publicKey, signatureAlgorithm: SignatureAlgorithm.ECDSA_P256)

	  var i=0
		while i < number {

        signer.keys.add(
            publicKey: pk,
            hashAlgorithm: HashAlgorithm.SHA3_256,
            weight: 1000.0
        )
				i=i+1
			}
	}
}

