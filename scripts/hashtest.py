#!/usr/bin/python3

import base64
import hashlib

private_key = b"ag1kZXZ-ZnV0YXJjaGlvcg0LEgRVc2VyIgNlcmIM#5577006791947779410"
h = hashlib.sha256()
h.update(private_key)
public_key = h.digest()
public_key_b64 = base64.urlsafe_b64encode(public_key)
print(public_key_b64)
