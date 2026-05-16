use axum::routing::{get, post};
use axum::Router;
use sqlx::migrate::MigrateDatabase;
use tower_http::cors::CorsLayer;
use tower_http::trace::TraceLayer;
use tracing_subscriber::EnvFilter;

use attestation_service::config;
use attestation_service::db;
use attestation_service::handlers;
use attestation_service::AppState;

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

    db::run_migrations(&pool)
        .await
        .expect("failed to run migrations");

    let port = cfg.port;

    let state = AppState {
        db: pool,
        config: cfg,
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

    let addr = format!("0.0.0.0:{port}");
    tracing::info!("starting attestation service on {addr}");

    let listener = tokio::net::TcpListener::bind(&addr)
        .await
        .expect("failed to bind address");

    axum::serve(listener, app)
        .await
        .expect("server failed");
}
