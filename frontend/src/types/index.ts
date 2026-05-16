export interface User {
  id: string
  walletAddress: string
  role: "admin" | "moderator" | "member"
  displayName?: string
  isActive: boolean
  createdAt: string
}

export interface Attestation {
  id: string
  userIdHash: string
  proofType: "noir" | "risc0"
  proofHash: string
  circuitType: string
  isVerified: boolean
  createdAt: string
}

export interface CheckIn {
  id: string
  userId: string
  scheduledAt: string
  windowMinutes: number
  status: "pending" | "completed" | "missed" | "emergency"
  completedAt?: string
  createdAt: string
}

export interface EmergencyAlert {
  id: string
  userId: string
  triggerType: string
  locationData?: object
  resolvedAt?: string
  createdAt: string
}

export interface BadClientReport {
  id: string
  encryptedDetails: string
  locationRegion?: string
  reportCategory: string
  status: "pending" | "verified" | "dismissed"
  confidenceScore: number
  createdAt: string
}

export interface DAOProposal {
  id: string
  title: string
  description: string
  proposalType: string
  status: string
  votingType: string
  startsAt: string
  endsAt: string
  createdBy: string
}

export interface DAOVote {
  id: string
  proposalId: string
  voterId: string
  vote: "yes" | "no" | "abstain"
  votingWeight: number
  createdAt: string
}
