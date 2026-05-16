using Genie.Renderer.Json

route("/health") do
    Dict("status" => "ok", "service" => "analytics")
end

route("/api/v1/analytics/dedup", method=POST) do
    body = rawpayload()
    if isempty(body)
        return Respond(400, JSON3.write(Dict("error" => "Empty request body")))
    end

    local data
    try
        data = JSON3.parse(body)
    catch e
        @warn "JSON parse error in /dedup" exception=(e, catch_backtrace())
        return Respond(400, JSON3.write(Dict("error" => "Invalid JSON body")))
    end

    report_text = get(data, "report_text", "")
    report_id = get(data, "report_id", "")

    if isempty(report_text)
        return Respond(400, JSON3.write(Dict("error" => "Missing 'report_text' field")))
    end
    if isempty(report_id)
        return Respond(400, JSON3.write(Dict("error" => "Missing 'report_id' field")))
    end

    try
        cfg = load_config()
        client = QdrantClient(cfg.qdrant_host, cfg.qdrant_port)
        result = deduplicate_report(client, report_id, report_text)
        return Respond(200, JSON3.write(result))
    catch e
        @error "Dedup endpoint failed" exception=(e, catch_backtrace())
        return Respond(500, JSON3.write(Dict(
            "error" => "Internal server error",
            "detail" => string(e),
        )))
    end
end

route("/api/v1/analytics/trends", method=GET) do
    try
        sample_data = DataFrame(
            category = ["safety", "harassment", "safety", "medical", "harassment", "safety", "medical", "legal"],
            region = ["north", "south", "east", "west", "north", "south", "east", "west"],
            date = [Date("2025-01-15"), Date("2025-02-20"), Date("2025-03-10"),
                    Date("2025-01-25"), Date("2025-02-05"), Date("2025-03-15"),
                    Date("2025-01-30"), Date("2025-02-28")],
        )
        result = report_trends(sample_data)
        return Respond(200, JSON3.write(result))
    catch e
        @error "Trends endpoint failed" exception=(e, catch_backtrace())
        return Respond(500, JSON3.write(Dict("error" => "Internal server error")))
    end
end

route("/api/v1/analytics/dao", method=GET) do
    try
        sample_votes = DataFrame(
            proposal_id = ["prop-1", "prop-1", "prop-1", "prop-2", "prop-2", "prop-3", "prop-3", "prop-3", "prop-3", "prop-4"],
            vote = ["for", "against", "for", "for", "abstain", "for", "for", "against", "for", "abstain"],
            voter_id = ["voter-a", "voter-b", "voter-c", "voter-a", "voter-d", "voter-a", "voter-b", "voter-e", "voter-c", "voter-f"],
            turnout = [0.75, 0.75, 0.75, 0.60, 0.60, 0.80, 0.80, 0.80, 0.80, 0.45],
        )
        result = dao_participation(sample_votes)
        return Respond(200, JSON3.write(result))
    catch e
        @error "DAO endpoint failed" exception=(e, catch_backtrace())
        return Respond(500, JSON3.write(Dict("error" => "Internal server error")))
    end
end

route("/api/v1/analytics/embed", method=POST) do
    body = rawpayload()
    if isempty(body)
        return Respond(400, JSON3.write(Dict("error" => "Empty request body")))
    end

    local data
    try
        data = JSON3.parse(body)
    catch e
        @warn "JSON parse error in /embed" exception=(e, catch_backtrace())
        return Respond(400, JSON3.write(Dict("error" => "Invalid JSON body")))
    end

    text = get(data, "text", "")
    if isempty(text)
        return Respond(400, JSON3.write(Dict("error" => "Missing 'text' field")))
    end

    try
        embedding = compute_embedding(text)
        return Respond(200, JSON3.write(Dict("embedding" => embedding)))
    catch e
        @error "Embedding computation failed" exception=(e, catch_backtrace())
        return Respond(500, JSON3.write(Dict("error" => "Internal server error")))
    end
end
