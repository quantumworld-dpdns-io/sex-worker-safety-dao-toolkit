use async_trait::async_trait;
use bytes::Bytes;
use http::HeaderValue;
use pingora::http::{RequestHeader, ResponseHeader};
use pingora::proxy::{ProxyHttp, Session};
use pingora::upstreams::peer::HttpPeer;
use pingora::Result;
use std::sync::Arc;
use uuid::Uuid;

pub struct ProxyConfig {
    pub attestation_upstream: String,
    pub checkin_upstream: String,
    pub analytics_upstream: String,
}

impl ProxyConfig {
    pub fn route<'a>(&'a self, path: &'a str) -> (&'a str, Option<&'a str>) {
        if path == "/health" {
            return ("", None);
        }

        let stripped = path.strip_prefix("/api/v1");

        if let Some(rest) = stripped {
            if rest.starts_with("/attestations") {
                return (&self.attestation_upstream, Some(rest));
            }
            if rest.starts_with("/checkins")
                || rest.starts_with("/emergency")
                || rest.starts_with("/dao")
                || rest.starts_with("/registry")
            {
                return (&self.checkin_upstream, Some(rest));
            }
            if rest.starts_with("/analytics") {
                return (&self.analytics_upstream, Some(rest));
            }
        }

        (&self.checkin_upstream, None)
    }
}

fn extract_addr(upstream: &str) -> &str {
    upstream
        .strip_prefix("http://")
        .or_else(|| upstream.strip_prefix("https://"))
        .unwrap_or(upstream)
}

#[derive(Clone)]
pub struct GatewayProxy {
    pub config: Arc<ProxyConfig>,
    pub jwt_secret: Arc<Vec<u8>>,
    pub rate_limiter: Arc<crate::ratelimit::RateLimiter>,
}

#[async_trait]
impl ProxyHttp for GatewayProxy {
    type CTX = ();

    fn new_ctx(&self) -> Self::CTX {}

    async fn upstream_peer(&self, session: &mut Session, _ctx: &mut ()) -> Result<Box<HttpPeer>> {
        let path = session.req_header().uri.path();
        let (upstream, _) = self.config.route(path);

        if upstream.is_empty() {
            return Ok(Box::new(HttpPeer::new("127.0.0.1:8080", false, "".to_string())));
        }

        let addr = extract_addr(upstream);
        let peer = Box::new(HttpPeer::new(addr, false, "".to_string()));
        Ok(peer)
    }

    async fn request_filter(&self, session: &mut Session, _ctx: &mut ()) -> Result<bool> {
        let path = session.req_header().uri.path().to_string();

        if path == "/health" {
            let resp = ResponseHeader::build(200, None)?;
            session.write_response_header(Box::new(resp), false).await?;
            session.write_response_body(Some(Bytes::from_static(b"OK")), true).await?;
            tracing::info!(path = %path, "health check");
            return Ok(true);
        }

        let client_ip = session
            .req_header()
            .headers
            .get("X-Forwarded-For")
            .and_then(|v| v.to_str().ok())
            .unwrap_or("unknown")
            .to_string();

        if !self.rate_limiter.check_rate_limit(&client_ip) {
            session
                .respond_error_with_body(429, Bytes::from_static(b"rate limit exceeded"))
                .await?;
            tracing::warn!(ip = %client_ip, "rate limit exceeded");
            return Ok(true);
        }

        if !path.starts_with("/api/v1/auth") && !path.starts_with("/health") {
            let auth_header = session
                .req_header()
                .headers
                .get(http::header::AUTHORIZATION)
                .and_then(|v| v.to_str().ok());

            if let Some(header_value) = auth_header {
                if let Some(token) = header_value.strip_prefix("Bearer ") {
                    match crate::auth::validate_token(token, &self.jwt_secret) {
                        Ok(claims) => {
                            session
                                .req_header_mut()
                                .headers
                                .insert("X-User-Id", HeaderValue::from_str(&claims.sub).unwrap());
                            session
                                .req_header_mut()
                                .headers
                                .insert("X-User-Role", HeaderValue::from_str(&claims.role).unwrap());
                        }
                        Err(e) => {
                            tracing::warn!(error = %e, "JWT validation failed");
                            session
                                .respond_error_with_body(401, Bytes::from_static(b"unauthorized"))
                                .await?;
                            return Ok(true);
                        }
                    }
                } else {
                    session
                        .respond_error_with_body(
                            401,
                            Bytes::from_static(b"invalid authorization scheme"),
                        )
                        .await?;
                    return Ok(true);
                }
            } else {
                session
                    .respond_error_with_body(
                        401,
                        Bytes::from_static(b"missing authorization header"),
                    )
                    .await?;
                return Ok(true);
            }
        }

        Ok(false)
    }

    async fn upstream_request_filter(
        &self,
        _session: &mut Session,
        upstream_request: &mut RequestHeader,
        _ctx: &mut (),
    ) -> Result<()> {
        upstream_request
            .insert_header(
                "X-Request-ID",
                HeaderValue::from_str(&Uuid::new_v4().to_string()).unwrap(),
            )
            .unwrap();

        let path = upstream_request.uri.path().to_string();
        if let Some(stripped) = path.strip_prefix("/api/v1") {
            let new_uri: http::Uri = stripped.parse().map_err(|_| {
                pingora::Error::new(pingora::ErrorType::InternalError)
            })?;
            upstream_request.set_uri(new_uri);
        }

        Ok(())
    }
}
