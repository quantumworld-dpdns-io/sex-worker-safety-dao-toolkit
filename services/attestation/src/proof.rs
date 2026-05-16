use sha2::{Digest, Sha256};

#[derive(Debug)]
pub enum ProofError {
    Generation(String),
    Verification(String),
}

impl std::fmt::Display for ProofError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            ProofError::Generation(e) => write!(f, "proof generation error: {e}"),
            ProofError::Verification(e) => write!(f, "proof verification error: {e}"),
        }
    }
}

impl std::error::Error for ProofError {}

#[derive(Debug, Clone)]
pub struct ProofOutput {
    pub proof_data: String,
    pub proof_hash: String,
    pub public_outputs: Vec<u8>,
}

pub trait ProofEngine: Send + Sync {
    fn generate_proof(
        &self,
        circuit_type: &str,
        private_inputs: &[u8],
        public_inputs: &[u8],
    ) -> Result<ProofOutput, ProofError>;

    fn verify_proof(
        &self,
        proof_data: &str,
        public_inputs: &[u8],
    ) -> Result<bool, ProofError>;
}

pub struct NoirProver;

impl NoirProver {
    pub fn new() -> Self {
        Self
    }
}

impl ProofEngine for NoirProver {
    fn generate_proof(
        &self,
        circuit_type: &str,
        private_inputs: &[u8],
        public_inputs: &[u8],
    ) -> Result<ProofOutput, ProofError> {
        let mut hasher = Sha256::new();
        hasher.update(circuit_type.as_bytes());
        hasher.update(private_inputs);
        hasher.update(public_inputs);
        let hash = hasher.finalize();
        let proof_hash = hex::encode(hash);

        let proof_data = format!(
            "mock_proof:{}:{}",
            proof_hash,
            hex::encode(private_inputs)
        );

        Ok(ProofOutput {
            proof_data,
            proof_hash,
            public_outputs: public_inputs.to_vec(),
        })
    }

    fn verify_proof(
        &self,
        proof_data: &str,
        public_inputs: &[u8],
    ) -> Result<bool, ProofError> {
        let parts: Vec<&str> = proof_data.splitn(2, ':').collect();
        if parts.len() != 2 {
            return Ok(false);
        }

        let expected_prefix = "mock_proof";
        if parts[0] != expected_prefix {
            return Ok(false);
        }

        let hash_str = parts[1];
        let hash_bytes = hex::decode(hash_str)
            .map_err(|_| ProofError::Verification("invalid hex in proof data".into()))?;

        let mut hasher = Sha256::new();
        hasher.update(public_inputs);
        let expected = hasher.finalize();

        Ok(hash_bytes == expected.as_slice())
    }
}

pub fn generate_proof_hash(proof_data: &[u8]) -> String {
    let mut hasher = Sha256::new();
    hasher.update(proof_data);
    hex::encode(hasher.finalize())
}
