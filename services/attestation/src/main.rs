use axum::routing::{get, post};
use axum::Router;
use sqlx::migrate::MigrateDatabase;
use tower_http::cors::CorsLayer;
use tower_http::trace::TraceLayer;
use tracing_subscriber::EnvFilter;

mod config;
mod crypto;
mod db;
mod handlers;
mod lib;
mod models;
mod proof;

use lib::AppState;

#[tokio::main]
async fn main() {
    tracing_subscriber::fmt()
        .json()
        .with_env_filter(EnvFilter::from_default_env())
        .init();

    let cfg = config::Config::from_env();

    if !sqlx::Postgres::database_exists(&cfg.database_url)
        .await
        .unwrap_or(false)
    {
        sqlx::Postgres::create_database(&cfg.database_url)
            .await
            .expect("failed to create database");
    }

    let pool = db::init_pool(&cfg.database_url)
        .await
        .expect("failed to connect to database");

    sqlx::migrate!()
        .run(&pool)
        .await
        .expect("failed to run migrations");

    let state = AppState {
        db: pool,
        config: cfg.clone(),
    };

    let app = Router::new()
        .route("/health", get(handlers::health))
        .route(
            "/api/attestations",
            post(handlers::create_attestation),
        )
        .route(
            "/api/attestations/{id}",
            get(handlers::get_attestation),
        )
        .route(
            "/api/attestations/{id}/verify",
            post(handlers::verify_attestation),
        )
        .route("/api/circuits", get(handlers::list_circuits))
        .layer(TraceLayer::new_for_http())
        .layer(CorsLayer::permissive())
        .with_state(state);

    let addr = format!("0.0.0.0:{}", cfg.port);
    tracing::info!("starting attestation service on {addr}");

    let listener = tokio::net::TcpListener::bind(&addr)
        .await
        .expect("failed to bind address");

    axum::serve(listener, app)
        .await
        .expect("server failed");
}
