pub mod config;
pub mod crypto;
pub mod db;
pub mod handlers;
pub mod models;
pub mod proof;

use sqlx::PgPool;

#[derive(Clone)]
pub struct AppState {
    pub db: PgPool,
    pub config: config::Config,
}
