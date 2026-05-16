use jsonwebtoken::{decode, encode, DecodingKey, EncodingKey, Header, Validation};
use serde::{Deserialize, Serialize};
use std::time::{SystemTime, UNIX_EPOCH};

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct Claims {
    pub sub: String,
    pub exp: usize,
    pub role: String,
    pub iat: usize,
}

pub fn create_token(
    user_id: &str,
    role: &str,
    secret: &[u8],
    expiry_hours: u64,
) -> String {
    let now = SystemTime::now()
        .duration_since(UNIX_EPOCH)
        .unwrap()
        .as_secs() as usize;

    let claims = Claims {
        sub: user_id.to_string(),
        exp: now + (expiry_hours as usize * 3600),
        role: role.to_string(),
        iat: now,
    };

    encode(
        &Header::new(jsonwebtoken::Algorithm::HS256),
        &claims,
        &EncodingKey::from_secret(secret),
    )
    .expect("JWT encoding failed")
}

pub fn validate_token(
    token: &str,
    secret: &[u8],
) -> Result<Claims, jsonwebtoken::errors::Error> {
    let validation = Validation::new(jsonwebtoken::Algorithm::HS256);
    let token_data = decode::<Claims>(
        token,
        &DecodingKey::from_secret(secret),
        &validation,
    )?;
    Ok(token_data.claims)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_create_and_validate_token() {
        let secret = b"mysecretkey12345678901234567890";
        let token = create_token("user123", "admin", secret, 24);
        let claims = validate_token(&token, secret).unwrap();
        assert_eq!(claims.sub, "user123");
        assert_eq!(claims.role, "admin");
        assert!(claims.exp > claims.iat);
    }
}
