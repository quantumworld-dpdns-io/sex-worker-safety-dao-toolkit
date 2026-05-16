use sqlx::postgres::PgPool;
use sqlx::Row;
use uuid::Uuid;

use crate::models::{Attestation, NewAttestation};

#[derive(Debug)]
pub enum DbError {
    Sqlx(sqlx::Error),
    NotFound,
}

impl std::fmt::Display for DbError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            DbError::Sqlx(e) => write!(f, "database error: {e}"),
            DbError::NotFound => write!(f, "not found"),
        }
    }
}

impl std::error::Error for DbError {}

impl From<sqlx::Error> for DbError {
    fn from(e: sqlx::Error) -> Self {
        DbError::Sqlx(e)
    }
}

pub async fn init_pool(database_url: &str) -> Result<PgPool, DbError> {
    let pool = PgPool::connect(database_url).await?;
    Ok(pool)
}

pub async fn create_attestation(
    pool: &PgPool,
    att: NewAttestation,
) -> Result<Attestation, DbError> {
    let id = Uuid::new_v4();
    let now = chrono::Utc::now();

    let row = sqlx::query(
        r#"
        INSERT INTO attestations (id, circuit_type, proof_data, proof_hash, user_id_hash, payload, is_verified, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING id, circuit_type, proof_data, proof_hash, user_id_hash, payload, is_verified, created_at, updated_at
        "#,
    )
    .bind(id)
    .bind(&att.circuit_type)
    .bind("")
    .bind("")
    .bind(&att.user_id_hash)
    .bind(&att.payload)
    .bind(false)
    .bind(now)
    .bind(now)
    .fetch_one(pool)
    .await?;

    Ok(Attestation {
        id: row.get("id"),
        circuit_type: row.get("circuit_type"),
        proof_data: row.get("proof_data"),
        proof_hash: row.get("proof_hash"),
        user_id_hash: row.get("user_id_hash"),
        payload: row.get("payload"),
        is_verified: row.get("is_verified"),
        created_at: row.get("created_at"),
        updated_at: row.get("updated_at"),
    })
}

pub async fn get_attestation(
    pool: &PgPool,
    id: Uuid,
) -> Result<Option<Attestation>, DbError> {
    let row = sqlx::query(
        r#"
        SELECT id, circuit_type, proof_data, proof_hash, user_id_hash, payload, is_verified, created_at, updated_at
        FROM attestations
        WHERE id = $1
        "#,
    )
    .bind(id)
    .fetch_optional(pool)
    .await?;

    Ok(row.map(|r| Attestation {
        id: r.get("id"),
        circuit_type: r.get("circuit_type"),
        proof_data: r.get("proof_data"),
        proof_hash: r.get("proof_hash"),
        user_id_hash: r.get("user_id_hash"),
        payload: r.get("payload"),
        is_verified: r.get("is_verified"),
        created_at: r.get("created_at"),
        updated_at: r.get("updated_at"),
    }))
}

pub async fn verify_attestation(pool: &PgPool, id: Uuid) -> Result<(), DbError> {
    let result = sqlx::query(
        r#"
        UPDATE attestations
        SET is_verified = true, updated_at = $2
        WHERE id = $1
        "#,
    )
    .bind(id)
    .bind(chrono::Utc::now())
    .execute(pool)
    .await?;

    if result.rows_affected() == 0 {
        return Err(DbError::NotFound);
    }

    Ok(())
}

pub async fn list_circuits(pool: &PgPool) -> Result<Vec<String>, DbError> {
    let rows = sqlx::query(
        r#"
        SELECT DISTINCT circuit_type
        FROM attestations
        ORDER BY circuit_type
        "#,
    )
    .fetch_all(pool)
    .await?;

    Ok(rows.into_iter().map(|r| r.get("circuit_type")).collect())
}
