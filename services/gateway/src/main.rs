mod auth;
mod config;
mod crypto;
mod middleware;
mod proxy;
mod ratelimit;
mod server;

#[tokio::main]
async fn main() {
    tracing_subscriber::fmt()
        .json()
        .with_env_filter(tracing_subscriber::EnvFilter::from_default_env())
        .init();

    let cfg = config::Config::from_env();

    tracing::info!(
        port = cfg.gateway_port,
        "starting API gateway"
    );

    let srv = server::create_server(cfg);
    srv.run_forever();
}
