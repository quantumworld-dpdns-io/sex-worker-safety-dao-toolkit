use async_trait::async_trait;
use http::{HeaderValue, Uri};
use pingora::http::{RequestHeader, ResponseHeader};
use pingora::proxy::{ProxyHttp, Session};
use pingora::upstreams::peer::HttpPeer;
use pingora::Result;
use std::str::FromStr;
use std::sync::Arc;
use uuid::Uuid;

pub struct ProxyConfig {
    pub attestation_upstream: String,
    pub checkin_upstream: String,
    pub analytics_upstream: String,
}

impl ProxyConfig {
    pub fn route(&self, path: &str) -> (&str, Option<&str>) {
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
            return Ok(Box::new(HttpPeer::new("127.0.0.1:8080", false, "")));
        }

        let peer = HttpPeer::new(upstream, false, "");
        Ok(Box::new(peer))
    }

    async fn request_filter(&self, session: &mut Session, _ctx: &mut ()) -> Result<()> {
        let path = session.req_header().uri.path().to_string();

        if path == "/health" {
            let resp = ResponseHeader::build(200, None)?;
            session.respond_header(resp).await?;
            session
                .respond_body(axum::body::Bytes::from("OK"), true)
                .await?;
            tracing::info!(path = %path, "health check");
            return Ok(());
        }

        let client_ip = session
            .req_header()
            .headers
            .get("X-Forwarded-For")
            .and_then(|v| v.to_str().ok())
            .unwrap_or("unknown")
            .to_string();

        if !self.rate_limiter.check_rate_limit(&client_ip) {
            let resp = ResponseHeader::build(429, None)?;
            session.respond_header(resp).await?;
            session
                .respond_body(axum::body::Bytes::from("rate limit exceeded"), true)
                .await?;
            tracing::warn!(ip = %client_ip, "rate limit exceeded");
            return Ok(());
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
                            let resp = ResponseHeader::build(401, None)?;
                            session.respond_header(resp).await?;
                            session
                                .respond_body(axum::body::Bytes::from("unauthorized"), true)
                                .await?;
                            return Ok(());
                        }
                    }
                } else {
                    let resp = ResponseHeader::build(401, None)?;
                    session.respond_header(resp).await?;
                    session
                        .respond_body(axum::body::Bytes::from("invalid authorization scheme"), true)
                        .await?;
                    return Ok(());
                }
            } else {
                let resp = ResponseHeader::build(401, None)?;
                session.respond_header(resp).await?;
                session
                    .respond_body(axum::body::Bytes::from("missing authorization header"), true)
                    .await?;
                return Ok(());
            }
        }

        Ok(())
    }

    async fn upstream_request_filter(
        &self,
        session: &mut Session,
        _ctx: &mut (),
    ) -> Result<()> {
        let headers = &mut session.req_header_mut().headers;

        headers.insert(
            "X-Request-ID",
            HeaderValue::from_str(&Uuid::new_v4().to_string()).unwrap(),
        );

        let path = session.req_header().uri.path().to_string();
        if let Some(stripped) = path.strip_prefix("/api/v1") {
            let new_uri = Uri::from_str(stripped)
                .map_err(|e| pingora::Error::new(pingora::ErrorType::InternalError))?;
            *session.req_header_mut().uri_mut() = new_uri;
        }

        Ok(())
    }

    async fn error_response(
        &self,
        session: &mut Session,
        _ctx: &mut (),
    ) -> Result<()> {
        let resp = ResponseHeader::build(502, None)?;
        session.respond_header(resp).await?;
        session
            .respond_body(axum::body::Bytes::from("upstream unavailable"), true)
            .await?;
        Ok(())
    }
}
