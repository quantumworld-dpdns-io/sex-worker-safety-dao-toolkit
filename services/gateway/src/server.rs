use std::sync::Arc;
use std::time::Duration;

use pingora::proxy::http_proxy_service;
use pingora::server::configuration::Opt;
use pingora::server::Server;
use pingora::services::background::background_service;

use crate::config::Config;
use crate::proxy::{GatewayProxy, ProxyConfig};
use crate::ratelimit::RateLimiter;

pub fn create_server(config: Config) -> Server {
    let mut server = Server::new(Some(Opt::default())).unwrap();
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
        if let Err(e) = setup_tls(domain, email) {
            tracing::error!(error = %e, "TLS setup failed, continuing without TLS");
        }
    }

    server
}

#[allow(unused)]
fn setup_tls(domain: &str, email: &str) -> Result<(), Box<dyn std::error::Error>> {
    tracing::info!(domain = %domain, email = %email, "TLS via ACME requested");

    use acme_lib::persist::FilePersist;
    use acme_lib::{Directory, DirectoryUrl};

    let persist = FilePersist::new("/tmp/acme-certs");
    let dir = Directory::from_url(persist, DirectoryUrl::LetsEncrypt)?;
    let acc = dir.account(email)?;

    match acc.certificate(domain)? {
        Some(cert) => {
            tracing::info!(
                domain = %domain,
                days_left = cert.valid_days_left(),
                "found existing certificate"
            );
            let _cert_der = cert.certificate_der();
            let _key_der = cert.private_key_der();
            tracing::info!(domain = %domain, "TLS certificate loaded from storage");
            Ok(())
        }
        None => {
            tracing::warn!(
                domain = %domain,
                "no existing certificate found, automatic issuance requires serving HTTP-01 \
                 challenges on port 80. Configure separately."
            );
            Ok(())
        }
    }
}

struct HealthCheckTask {
    upstreams: Vec<String>,
}

#[async_trait::async_trait]
impl pingora::services::background::BackgroundService for HealthCheckTask {
    async fn start(&self, _shutdown: tokio::sync::watch::Receiver<bool>) {
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
