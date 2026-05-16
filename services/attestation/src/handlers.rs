use axum::extract::State;
use axum::Json;
use uuid::Uuid;

use crate::crypto;
use crate::db;
use crate::AppState;
use crate::models::{AttestationResponse, Circuit, ProofRequest, ProofResponse};
use crate::proof::{NoirProver, ProofEngine};

#[derive(Debug)]
pub enum HandlerError {
    Db(db::DbError),
    Crypto(crypto::CryptoError),
    Proof(crate::proof::ProofError),
    Internal(String),
}

impl std::fmt::Display for HandlerError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            HandlerError::Db(e) => write!(f, "{e}"),
            HandlerError::Crypto(e) => write!(f, "{e}"),
            HandlerError::Proof(e) => write!(f, "{e}"),
            HandlerError::Internal(e) => write!(f, "internal error: {e}"),
        }
    }
}

impl std::error::Error for HandlerError {}

impl axum::response::IntoResponse for HandlerError {
    fn into_response(self) -> axum::response::Response {
        let status = match &self {
            HandlerError::Db(db::DbError::NotFound) => axum::http::StatusCode::NOT_FOUND,
            _ => axum::http::StatusCode::INTERNAL_SERVER_ERROR,
        };
        let body = serde_json::json!({ "error": self.to_string() });
        (status, axum::Json(body)).into_response()
    }
}

impl From<db::DbError> for HandlerError {
    fn from(e: db::DbError) -> Self {
        HandlerError::Db(e)
    }
}

impl From<crypto::CryptoError> for HandlerError {
    fn from(e: crypto::CryptoError) -> Self {
        HandlerError::Crypto(e)
    }
}

impl From<crate::proof::ProofError> for HandlerError {
    fn from(e: crate::proof::ProofError) -> Self {
        HandlerError::Proof(e)
    }
}

pub async fn create_attestation(
    State(state): State<AppState>,
    Json(body): Json<ProofRequest>,
) -> Result<Json<ProofResponse>, HandlerError> {
    let payload_bytes =
        serde_json::to_vec(&body.payload).map_err(|e| HandlerError::Internal(e.to_string()))?;

    let encrypted_payload = crypto::encrypt(&payload_bytes, &state.config.encryption_key)?;

    let salt = rand::random::<[u8; 16]>();
    let user_id_hash = crypto::hash_user_id(&body.user_id, &salt);

    let prover = NoirProver::new();
    let proof_output = prover.generate_proof(&body.circuit_type, &payload_bytes, &[])?;

    let new_att = crate::models::NewAttestation {
        circuit_type: body.circuit_type,
        user_id_hash,
        payload: encrypted_payload,
    };

    let mut attestation = db::create_attestation(&state.db, new_att).await?;
    attestation.proof_data = proof_output.proof_data;
    attestation.proof_hash = proof_output.proof_hash.clone();

    Ok(Json(ProofResponse {
        attestation_id: attestation.id,
        proof_hash: proof_output.proof_hash,
    }))
}

pub async fn get_attestation(
    State(state): State<AppState>,
    axum::extract::Path(id): axum::extract::Path<Uuid>,
) -> Result<Json<AttestationResponse>, HandlerError> {
    let attestation = db::get_attestation(&state.db, id)
        .await?
        .ok_or(db::DbError::NotFound)?;

    Ok(Json(AttestationResponse {
        id: attestation.id,
        circuit_type: attestation.circuit_type,
        proof_data: attestation.proof_data,
        proof_hash: attestation.proof_hash,
        user_id_hash: attestation.user_id_hash,
        is_verified: attestation.is_verified,
        created_at: attestation.created_at,
        updated_at: attestation.updated_at,
    }))
}

pub async fn verify_attestation(
    State(state): State<AppState>,
    axum::extract::Path(id): axum::extract::Path<Uuid>,
) -> Result<Json<serde_json::Value>, HandlerError> {
    let attestation = db::get_attestation(&state.db, id)
        .await?
        .ok_or(db::DbError::NotFound)?;

    let prover = NoirProver::new();
    let public_inputs = &[];

    let is_valid = prover.verify_proof(&attestation.proof_data, public_inputs)?;

    if is_valid {
        db::verify_attestation(&state.db, id).await?;
    }

    Ok(Json(serde_json::json!({
        "id": id,
        "is_valid": is_valid
    })))
}

pub async fn list_circuits(
    State(state): State<AppState>,
) -> Result<Json<Vec<Circuit>>, HandlerError> {
    let circuit_names = db::list_circuits(&state.db).await?;

    let circuits: Vec<Circuit> = circuit_names
        .into_iter()
        .map(|name| Circuit {
            description: format!("ZK circuit for {name}"),
            name,
        })
        .collect();

    Ok(Json(circuits))
}

pub async fn health() -> &'static str {
    "ok"
}
