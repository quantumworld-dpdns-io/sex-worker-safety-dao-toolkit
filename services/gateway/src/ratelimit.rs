use governor::clock::DefaultClock;
use governor::middleware::NoOpMiddleware;
use governor::state::direct::NotKeyed;
use governor::state::in_memory::InMemoryState;
use governor::{Quota, RateLimiter as GovRateLimiter};
use std::collections::HashMap;
use std::num::NonZeroU32;
use std::sync::Mutex;

type PerIpLimiter = GovRateLimiter<NotKeyed, InMemoryState, DefaultClock, NoOpMiddleware>;

pub struct RateLimiter {
    limiters: Mutex<HashMap<String, PerIpLimiter>>,
    quota: Quota,
}

impl RateLimiter {
    pub fn new(requests: u64, window_secs: u64) -> Self {
        let window = std::cmp::max(1, window_secs);
        let per_sec = std::cmp::max(1, (requests + window - 1) / window) as u32;
        let burst = NonZeroU32::new(requests as u32).unwrap();
        let rate = NonZeroU32::new(per_sec).unwrap();
        let quota = Quota::per_second(rate).allow_burst(burst);

        Self {
            limiters: Mutex::new(HashMap::new()),
            quota,
        }
    }

    pub fn check_rate_limit(&self, ip: &str) -> bool {
        let mut limiters = self.limiters.lock().unwrap();
        let limiter = limiters.entry(ip.to_string()).or_insert_with(|| {
            GovRateLimiter::direct(self.quota)
        });
        limiter.check().is_ok()
    }
}
