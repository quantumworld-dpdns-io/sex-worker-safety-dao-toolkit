use aes_gcm::aead::{Aead, KeyInit, OsRng};
use aes_gcm::{Aes256Gcm, Nonce};
use rand::RngCore;

#[derive(Debug)]
pub enum CryptoError {
    Aead(aes_gcm::Error),
    KeyLength,
}

impl std::fmt::Display for CryptoError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            CryptoError::Aead(e) => write!(f, "encryption error: {e}"),
            CryptoError::KeyLength => write!(f, "invalid key length"),
        }
    }
}

impl std::error::Error for CryptoError {}

impl From<aes_gcm::Error> for CryptoError {
    fn from(e: aes_gcm::Error) -> Self {
        CryptoError::Aead(e)
    }
}

pub fn encrypt(plaintext: &[u8], key: &[u8]) -> Result<Vec<u8>, CryptoError> {
    let cipher =
        Aes256Gcm::new_from_slice(key).map_err(|_| CryptoError::KeyLength)?;

    let mut nonce_bytes = [0u8; 12];
    OsRng.fill_bytes(&mut nonce_bytes);
    let nonce = Nonce::from_slice(&nonce_bytes);

    let ciphertext = cipher.encrypt(nonce, plaintext)?;

    let mut result = Vec::with_capacity(12 + ciphertext.len());
    result.extend_from_slice(&nonce_bytes);
    result.extend_from_slice(&ciphertext);

    Ok(result)
}

pub fn decrypt(ciphertext: &[u8], key: &[u8]) -> Result<Vec<u8>, CryptoError> {
    if ciphertext.len() < 12 {
        return Err(CryptoError::Aead(aes_gcm::Error));
    }

    let cipher =
        Aes256Gcm::new_from_slice(key).map_err(|_| CryptoError::KeyLength)?;

    let (nonce_bytes, ct) = ciphertext.split_at(12);
    let nonce = Nonce::from_slice(nonce_bytes);

    let plaintext = cipher.decrypt(nonce, ct)?;
    Ok(plaintext)
}

pub fn hash_user_id(user_id: &str, salt: &[u8]) -> String {
    use sha2::{Digest, Sha256};

    let mut hasher = Sha256::new();
    hasher.update(salt);
    hasher.update(user_id.as_bytes());
    hex::encode(hasher.finalize())
}
