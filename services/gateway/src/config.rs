use std::env;

#[derive(Debug, Clone)]
pub struct Config {
    pub gateway_port: u16,
    pub jwt_secret: String,
    pub jwt_expiry_hours: u64,
    pub upstream_attestation: String,
    pub upstream_checkin: String,
    pub upstream_analytics: String,
    pub rate_limit_requests: u64,
    pub rate_limit_window_secs: u64,
    pub acme_domain: Option<String>,
    pub acme_email: Option<String>,
    pub encryption_key: Option<String>,
}

impl Config {
    pub fn from_env() -> Self {
        Self {
            gateway_port: env::var("GATEWAY_PORT")
                .ok()
                .and_then(|v| v.parse().ok())
                .unwrap_or(8080),
            jwt_secret: env::var("JWT_SECRET")
                .expect("JWT_SECRET must be set"),
            jwt_expiry_hours: env::var("JWT_EXPIRY_HOURS")
                .ok()
                .and_then(|v| v.parse().ok())
                .unwrap_or(24),
            upstream_attestation: env::var("UPSTREAM_ATTESTATION")
                .unwrap_or_else(|_| "http://localhost:3001".to_string()),
            upstream_checkin: env::var("UPSTREAM_CHECKIN")
                .unwrap_or_else(|_| "http://localhost:3002".to_string()),
            upstream_analytics: env::var("UPSTREAM_ANALYTICS")
                .unwrap_or_else(|_| "http://localhost:3003".to_string()),
            rate_limit_requests: env::var("RATE_LIMIT_REQUESTS")
                .ok()
                .and_then(|v| v.parse().ok())
                .unwrap_or(100),
            rate_limit_window_secs: env::var("RATE_LIMIT_WINDOW_SECS")
                .ok()
                .and_then(|v| v.parse().ok())
                .unwrap_or(60),
            acme_domain: env::var("ACME_DOMAIN").ok(),
            acme_email: env::var("ACME_EMAIL").ok(),
            encryption_key: env::var("ENCRYPTION_KEY").ok(),
        }
    }

    pub fn jwt_secret_bytes(&self) -> Vec<u8> {
        self.jwt_secret.as_bytes().to_vec()
    }

    pub fn encryption_key_bytes(&self) -> Option<Vec<u8>> {
        self.encryption_key.as_ref().map(|hex_str| {
            hex::decode(hex_str).expect("ENCRYPTION_KEY must be valid hex (64 hex chars = 32 bytes)")
        })
    }
}
