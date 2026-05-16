const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1"

class ApiError extends Error {
  status: number
  constructor(message: string, status: number) {
    super(message)
    this.name = "ApiError"
    this.status = status
  }
}

async function request<T>(
  endpoint: string,
  options?: RequestInit,
): Promise<T> {
  const url = `${API_URL}${endpoint}`
  const token =
    typeof window !== "undefined" ? localStorage.getItem("sws_auth_token") : null

  const headers: HeadersInit = {
    "Content-Type": "application/json",
    ...options?.headers,
  }

  if (token) {
    ;(headers as Record<string, string>)["Authorization"] = `Bearer ${token}`
  }

  const res = await fetch(url, { ...options, headers })

  if (!res.ok) {
    const body = await res.text().catch(() => "")
    throw new ApiError(
      body || `Request failed with status ${res.status}`,
      res.status,
    )
  }

  if (res.status === 204) return undefined as T
  return res.json()
}

function buildQueryString(params?: Record<string, string | number | boolean | undefined>): string {
  if (!params) return ""
  const search = new URLSearchParams()
  for (const [key, value] of Object.entries(params)) {
    if (value !== undefined) search.set(key, String(value))
  }
  const qs = search.toString()
  return qs ? `?${qs}` : ""
}

export const api = {
  createAttestation: (data: Record<string, unknown>) =>
    request<{ id: string }>("/attestations", {
      method: "POST",
      body: JSON.stringify(data),
    }),

  getAttestation: (id: string) =>
    request(`/attestations/${id}`),

  verifyAttestation: (id: string) =>
    request<{ verified: boolean }>(`/attestations/${id}/verify`, {
      method: "POST",
    }),

  createCheckIn: (data: Record<string, unknown>) =>
    request<{ id: string }>("/checkins", {
      method: "POST",
      body: JSON.stringify(data),
    }),

  completeCheckIn: (id: string) =>
    request(`/checkins/${id}/complete`, { method: "POST" }),

  listCheckIns: () =>
    request<unknown[]>("/checkins"),

  triggerAlert: (data: Record<string, unknown>) =>
    request<{ id: string }>("/emergency", {
      method: "POST",
      body: JSON.stringify(data),
    }),

  resolveAlert: (id: string) =>
    request(`/emergency/${id}/resolve`, { method: "POST" }),

  submitReport: (data: Record<string, unknown>) =>
    request<{ id: string }>("/registry/reports", {
      method: "POST",
      body: JSON.stringify(data),
    }),

  searchReports: (params?: Record<string, string | number | boolean | undefined>) =>
    request<unknown[]>(`/registry/reports${buildQueryString(params)}`),

  createProposal: (data: Record<string, unknown>) =>
    request<{ id: string }>("/dao/proposals", {
      method: "POST",
      body: JSON.stringify(data),
    }),

  listProposals: () =>
    request<unknown[]>("/dao/proposals"),

  castVote: (proposalId: string, data: Record<string, unknown>) =>
    request(`/dao/proposals/${proposalId}/votes`, {
      method: "POST",
      body: JSON.stringify(data),
    }),
}
