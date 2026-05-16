use aes_gcm::aead::{Aead, AeadCore, KeyInit, OsRng};
use aes_gcm::{Aes256Gcm, Key, Nonce};

/// Encrypt plaintext using AES-256-GCM.
/// Returns the nonce (first 12 bytes) concatenated with the ciphertext.
pub fn encrypt_payload(plaintext: &[u8], key: &[u8]) -> Vec<u8> {
    let key = Key::<Aes256Gcm>::from_slice(key);
    let cipher = Aes256Gcm::new(key);
    let nonce = Aes256Gcm::generate_nonce(&mut OsRng);
    let ciphertext = cipher
        .encrypt(&nonce, plaintext)
        .expect("encryption failure");
    [nonce.as_slice(), ciphertext.as_slice()].concat()
}

/// Decrypt ciphertext that was encrypted with `encrypt_payload`.
/// Expects the nonce as the first 12 bytes followed by the ciphertext.
pub fn decrypt_payload(ciphertext: &[u8], key: &[u8]) -> Vec<u8> {
    let key = Key::<Aes256Gcm>::from_slice(key);
    let cipher = Aes256Gcm::new(key);
    let nonce = Nonce::from_slice(&ciphertext[..12]);
    cipher
        .decrypt(nonce, &ciphertext[12..])
        .expect("decryption failure")
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_roundtrip() {
        let key = b"01234567890123456789012345678901"; // 32 bytes
        let plaintext = b"sensitive data here";
        let encrypted = encrypt_payload(plaintext, key);
        let decrypted = decrypt_payload(&encrypted, key);
        assert_eq!(decrypted, plaintext);
    }
}
