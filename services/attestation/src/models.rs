use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use uuid::Uuid;

#[derive(Debug, Clone, Serialize, Deserialize, sqlx::FromRow)]
pub struct Attestation {
    pub id: Uuid,
    pub circuit_type: String,
    pub proof_data: String,
    pub proof_hash: String,
    pub user_id_hash: String,
    pub payload: Vec<u8>,
    pub is_verified: bool,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct NewAttestation {
    pub circuit_type: String,
    pub user_id_hash: String,
    pub payload: Vec<u8>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Circuit {
    pub name: String,
    pub description: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ProofRequest {
    pub circuit_type: String,
    pub user_id: String,
    pub payload: serde_json::Value,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ProofResponse {
    pub attestation_id: Uuid,
    pub proof_hash: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AttestationResponse {
    pub id: Uuid,
    pub circuit_type: String,
    pub proof_data: String,
    pub proof_hash: String,
    pub user_id_hash: String,
    pub is_verified: bool,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}
