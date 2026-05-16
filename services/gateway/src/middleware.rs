use std::future::Future;
use std::pin::Pin;
use std::sync::Arc;
use std::task::{Context, Poll};

use http::{Request, Response, StatusCode};
use tower::{Layer, Service};

use crate::auth;

/// Extracts and validates JWT from the Authorization header (Bearer scheme).
/// Injects `X-User-Id` and `X-User-Role` headers into the request for
/// downstream services. Non-authenticated routes (`/health`, `/api/v1/auth/*`)
/// are passed through unchanged.
#[derive(Clone)]
pub struct AuthMiddleware<S> {
    inner: S,
    jwt_secret: Arc<Vec<u8>>,
}

impl<S> AuthMiddleware<S> {
    pub fn new(inner: S, jwt_secret: Vec<u8>) -> Self {
        Self {
            inner,
            jwt_secret: Arc::new(jwt_secret),
        }
    }
}

impl<S, ReqBody> Service<Request<ReqBody>> for AuthMiddleware<S>
where
    S: Service<Request<ReqBody>, Response = Response<axum::body::Body>> + Clone + Send + 'static,
    S::Error: Into<Box<dyn std::error::Error + Send + Sync>>,
    S::Future: Send + 'static,
    ReqBody: Send + 'static,
{
    type Response = Response<axum::body::Body>;
    type Error = Box<dyn std::error::Error + Send + Sync>;
    type Future = Pin<Box<dyn Future<Output = Result<Self::Response, Self::Error>> + Send>>;

    fn poll_ready(&mut self, cx: &mut Context<'_>) -> Poll<Result<(), Self::Error>> {
        self.inner.poll_ready(cx).map_err(Into::into)
    }

    fn call(&mut self, req: Request<ReqBody>) -> Self::Future {
        let path = req.uri().path().to_string();

        if path == "/health" || path.starts_with("/api/v1/auth") {
            let future = self.inner.call(req);
            return Box::pin(async move { future.await.map_err(Into::into) });
        }

        let secret = self.jwt_secret.clone();
        let mut inner = self.inner.clone();

        Box::pin(async move {
            let auth_header = req
                .headers()
                .get(http::header::AUTHORIZATION)
                .and_then(|v| v.to_str().ok());

            let token = match auth_header.and_then(|h| h.strip_prefix("Bearer ")) {
                Some(t) => t,
                None => {
                    return Ok(Response::builder()
                        .status(StatusCode::UNAUTHORIZED)
                        .body(axum::body::Body::from("missing or invalid authorization header"))
                        .unwrap());
                }
            };

            match auth::validate_token(token, &secret) {
                Ok(claims) => {
                    let (mut parts, body) = req.into_parts();
                    parts.headers.insert(
                        "X-User-Id",
                        claims.sub.parse().unwrap(),
                    );
                    parts.headers.insert(
                        "X-User-Role",
                        claims.role.parse().unwrap(),
                    );
                    let req = Request::from_parts(parts, body);
                    inner.call(req).await.map_err(Into::into)
                }
                Err(e) => {
                    tracing::warn!(error = %e, token = %token, "JWT validation failed in middleware");
                    Ok(Response::builder()
                        .status(StatusCode::UNAUTHORIZED)
                        .body(axum::body::Body::from("unauthorized"))
                        .unwrap())
                }
            }
        })
    }
}

#[derive(Clone)]
pub struct AuthLayer {
    jwt_secret: Vec<u8>,
}

impl AuthLayer {
    pub fn new(jwt_secret: Vec<u8>) -> Self {
        Self { jwt_secret }
    }
}

impl<S> Layer<S> for AuthLayer {
    type Service = AuthMiddleware<S>;

    fn layer(&self, inner: S) -> Self::Service {
        AuthMiddleware::new(inner, self.jwt_secret.clone())
    }
}

/// Checks rate limit per client IP using the shared `RateLimiter`.
/// Returns 429 Too Many Requests when the limit is exceeded.
#[derive(Clone)]
pub struct RateLimitMiddleware<S> {
    inner: S,
    rate_limiter: Arc<crate::ratelimit::RateLimiter>,
}

impl<S> RateLimitMiddleware<S> {
    pub fn new(inner: S, rate_limiter: Arc<crate::ratelimit::RateLimiter>) -> Self {
        Self {
            inner,
            rate_limiter,
        }
    }
}

impl<S, ReqBody> Service<Request<ReqBody>> for RateLimitMiddleware<S>
where
    S: Service<Request<ReqBody>, Response = Response<axum::body::Body>> + Clone + Send + 'static,
    S::Error: Into<Box<dyn std::error::Error + Send + Sync>>,
    S::Future: Send + 'static,
    ReqBody: Send + 'static,
{
    type Response = Response<axum::body::Body>;
    type Error = Box<dyn std::error::Error + Send + Sync>;
    type Future = Pin<Box<dyn Future<Output = Result<Self::Response, Self::Error>> + Send>>;

    fn poll_ready(&mut self, cx: &mut Context<'_>) -> Poll<Result<(), Self::Error>> {
        self.inner.poll_ready(cx).map_err(Into::into)
    }

    fn call(&mut self, req: Request<ReqBody>) -> Self::Future {
        let client_ip = req
            .headers()
            .get("X-Forwarded-For")
            .and_then(|v| v.to_str().ok())
            .unwrap_or("unknown")
            .to_string();

        let limiter = self.rate_limiter.clone();
        let mut inner = self.inner.clone();

        Box::pin(async move {
            if !limiter.check_rate_limit(&client_ip) {
                tracing::warn!(ip = %client_ip, "rate limit exceeded in middleware");
                return Ok(Response::builder()
                    .status(StatusCode::TOO_MANY_REQUESTS)
                    .body(axum::body::Body::from("rate limit exceeded"))
                    .unwrap());
            }
            inner.call(req).await.map_err(Into::into)
        })
    }
}

#[derive(Clone)]
pub struct RateLimitLayer {
    rate_limiter: Arc<crate::ratelimit::RateLimiter>,
}

impl RateLimitLayer {
    pub fn new(rate_limiter: Arc<crate::ratelimit::RateLimiter>) -> Self {
        Self { rate_limiter }
    }
}

impl<S> Layer<S> for RateLimitLayer {
    type Service = RateLimitMiddleware<S>;

    fn layer(&self, inner: S) -> Self::Service {
        RateLimitMiddleware::new(inner, self.rate_limiter.clone())
    }
}

/// Logs method, path, status code, and duration for every request.
#[derive(Clone)]
pub struct RequestLogMiddleware<S> {
    inner: S,
}

impl<S> RequestLogMiddleware<S> {
    pub fn new(inner: S) -> Self {
        Self { inner }
    }
}

impl<S, ReqBody, ResBody> Service<Request<ReqBody>> for RequestLogMiddleware<S>
where
    S: Service<Request<ReqBody>, Response = Response<ResBody>> + Clone + Send + 'static,
    S::Error: std::error::Error + Send + Sync + 'static,
    S::Future: Send + 'static,
    ReqBody: Send + 'static,
    ResBody: Send + 'static,
{
    type Response = Response<ResBody>;
    type Error = S::Error;
    type Future = Pin<Box<dyn Future<Output = Result<Self::Response, Self::Error>> + Send>>;

    fn poll_ready(&mut self, cx: &mut Context<'_>) -> Poll<Result<(), Self::Error>> {
        self.inner.poll_ready(cx)
    }

    fn call(&mut self, req: Request<ReqBody>) -> Self::Future {
        let start = std::time::Instant::now();
        let method = req.method().clone();
        let path = req.uri().path().to_string();
        let mut inner = self.inner.clone();

        Box::pin(async move {
            let response = inner.call(req).await?;
            let status = response.status();
            let duration = start.elapsed();
            tracing::info!(
                method = %method,
                path = %path,
                status = status.as_u16(),
                duration_ms = duration.as_millis() as u64,
                "request completed"
            );
            Ok(response)
        })
    }
}

#[derive(Clone)]
pub struct RequestLogLayer;

impl<S> Layer<S> for RequestLogLayer {
    type Service = RequestLogMiddleware<S>;

    fn layer(&self, inner: S) -> Self::Service {
        RequestLogMiddleware::new(inner)
    }
}
