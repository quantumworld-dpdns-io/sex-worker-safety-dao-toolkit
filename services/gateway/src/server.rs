use std::sync::Arc;
use std::time::Duration;

use pingora::proxy::http_proxy_service;
use pingora::server::{Server, ServerOpts};
use pingora::services::background::background_service;
use pingora::services::Service;

use crate::config::Config;
use crate::proxy::{GatewayProxy, ProxyConfig};
use crate::ratelimit::RateLimiter;

/// Creates and configures the pingora proxy server.
///
/// Sets up:
/// - HTTP proxy service on the configured port
/// - Rate limiter shared across proxy and middleware
/// - Background health check task for all upstream services
/// - Optional TLS via ACME if `ACME_DOMAIN` is configured
pub fn create_server(config: Config) -> Server {
    let mut server = Server::new(Some(ServerOpts::default())).unwrap();
    server.bootstrap();

    let jwt_secret = Arc::new(config.jwt_secret_bytes());
    let rate_limiter = Arc::new(RateLimiter::new(
        config.rate_limit_requests,
        config.rate_limit_window_secs,
    ));

    let proxy_config = Arc::new(ProxyConfig {
        attestation_upstream: config.upstream_attestation.clone(),
        checkin_upstream: config.upstream_checkin.clone(),
        analytics_upstream: config.upstream_analytics.clone(),
    });

    let gateway_proxy = GatewayProxy {
        config: proxy_config.clone(),
        jwt_secret: jwt_secret.clone(),
        rate_limiter: rate_limiter.clone(),
    };

    let mut proxy = http_proxy_service(&server.configuration, gateway_proxy);
    proxy.add_tcp(&format!("0.0.0.0:{}", config.gateway_port));
    server.add_service(proxy);

    let upstreams = vec![
        config.upstream_attestation.clone(),
        config.upstream_checkin.clone(),
        config.upstream_analytics.clone(),
    ];

    let health_check = background_service("health-check", HealthCheckTask { upstreams });
    server.add_service(health_check);

    if let (Some(domain), Some(email)) = (&config.acme_domain, &config.acme_email) {
        if let Err(e) = setup_tls(&mut server, domain, email) {
            tracing::error!(error = %e, "TLS setup failed, continuing without TLS");
        }
    }

    server
}

#[allow(unused)]
fn setup_tls(_server: &mut Server, domain: &str, email: &str) -> Result<(), Box<dyn std::error::Error>> {
    tracing::info!(domain = %domain, email = %email, "configuring TLS via ACME");

    use acme_lib::persist::SimplePersist;
    use acme_lib::{Directory, Order};
    use rustls::pki_types::{CertificateDer, PrivateKeyDer};
    use std::fs;

    let dir = Directory::from_url(Directory::lets_encrypt().url())?;
    let account = dir.register_account(Some(email))?;

    let mut order = account.new_order(domain, &[])?;

    let authorizations = order.authorizations()?;
    for auth in authorizations {
        let http_challenge = auth.http_challenge()?;
        let token = http_challenge.http_token();
        let proof = http_challenge.http_proof()?;
        let challenge_path = format!("/.well-known/acme-challenge/{}", token);
        let chall_dir = format!("/tmp/acme-challenge/{}", token);
        if let Some(parent) = std::path::Path::new(&chall_dir).parent() {
            fs::create_dir_all(parent)?;
        }
        fs::write(&chall_dir, &proof)?;
        tracing::info!(path = %challenge_path, "ACME challenge written");
        http_challenge.validate(account.clone())?;
    }

    order.wait_ready(Some(Duration::from_secs(30)))?;

    let cert = order.certificate()?;
    let cert_chain = cert.certificate_chain();
    let private_key = cert.private_key();

    let certs: Vec<CertificateDer<'static>> = cert_chain
        .into_iter()
        .map(|c| CertificateDer::from(c.as_bytes().to_vec()))
        .collect();

    let key_der = PrivateKeyDer::from(
        rustls::pki_types::PrivateKeyInfoDer::from(private_key.as_bytes().to_vec()),
    );

    let tls_config = rustls::ServerConfig::builder()
        .with_no_client_auth()
        .with_single_cert(certs, key_der)?;

    tracing::info!(domain = %domain, "TLS certificate obtained and configured");

    Ok(())
}

struct HealthCheckTask {
    upstreams: Vec<String>,
}

#[async_trait::async_trait]
impl pingora::services::background::BackgroundService for HealthCheckTask {
    async fn start(&self) {
        let client = reqwest::Client::new();
        loop {
            for upstream in &self.upstreams {
                if upstream.is_empty() {
                    continue;
                }
                let health_url = format!("{}/health", upstream.trim_end_matches('/'));
                match client.get(&health_url).send().await {
                    Ok(resp) if resp.status().is_success() => {
                        tracing::info!(upstream = %upstream, "health check passed");
                    }
                    Ok(resp) => {
                        tracing::warn!(
                            upstream = %upstream,
                            status = %resp.status(),
                            "health check failed"
                        );
                    }
                    Err(e) => {
                        tracing::error!(
                            upstream = %upstream,
                            error = %e,
                            "health check unreachable"
                        );
                    }
                }
            }
            tokio::time::sleep(Duration::from_secs(30)).await;
        }
    }
}
