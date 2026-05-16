use std::env;

#[derive(Clone)]
pub struct Config {
    pub port: u16,
    pub database_url: String,
    pub encryption_key: Vec<u8>,
    pub jwt_secret: String,
}

impl Config {
    pub fn from_env() -> Self {
        let port = env::var("ATTESTATION_PORT")
            .unwrap_or_else(|_| "3001".to_string())
            .parse()
            .expect("ATTESTATION_PORT must be a valid u16");

        let database_url = env::var("DATABASE_URL").expect("DATABASE_URL must be set");
        let encryption_key_hex =
            env::var("ATTESTATION_ENCRYPTION_KEY").expect("ATTESTATION_ENCRYPTION_KEY must be set");
        let encryption_key = hex::decode(&encryption_key_hex)
            .expect("ATTESTATION_ENCRYPTION_KEY must be valid hex");
        assert_eq!(
            encryption_key.len(),
            32,
            "ATTESTATION_ENCRYPTION_KEY must be exactly 32 bytes (64 hex chars)"
        );

        let jwt_secret = env::var("JWT_SECRET").expect("JWT_SECRET must be set");

        Self {
            port,
            database_url,
            encryption_key,
            jwt_secret,
        }
    }
}
