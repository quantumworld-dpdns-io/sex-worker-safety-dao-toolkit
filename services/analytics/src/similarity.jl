function compute_embedding(text::String)::Vector{Float32}
    words = split(lowercase(text))
    embedding = zeros(Float32, 128)

    for word in words
        h = hash(word)
        idx1 = mod(h, 128) + 1
        idx2 = mod(h >> 32, 128) + 1
        embedding[idx1] += 1.0f0
        embedding[idx2] += 0.5f0
    end

    norm = sqrt(sum(embedding .^ 2))
    if norm > 0.0f0
        embedding ./= norm
    end

    return embedding
end

function cosine_similarity(a::Vector{Float32}, b::Vector{Float32})::Float32
    dot = sum(a .* b)
    norm_a = sqrt(sum(a .^ 2))
    norm_b = sqrt(sum(b .^ 2))
    if norm_a == 0.0f0 || norm_b == 0.0f0
        return 0.0f0
    end
    return dot / (norm_a * norm_b)
end

function find_duplicates(client::QdrantClient, text::String, threshold::Float64=0.85)::Vector{Dict}
    embedding = compute_embedding(text)
    results = search_points(client, "reports", embedding; limit=20)

    return filter(r -> r["score"] >= threshold, results)
end

function deduplicate_report(client::QdrantClient, report_id::String, report_text::String)::Dict
    embedding = compute_embedding(report_text)

    search_results = search_points(client, "reports", embedding; limit=10)

    high_similarity = filter(r -> r["score"] >= 0.85, search_results)

    if !isempty(high_similarity)
        best = high_similarity[1]
        return Dict(
            "is_duplicate" => true,
            "report_id" => report_id,
            "matches" => high_similarity,
            "best_match_id" => best["id"],
            "best_score" => best["score"],
            "stored" => false,
        )
    end

    try
        upsert_point(client, "reports", report_id, embedding, Dict(
            "report_id" => report_id,
            "timestamp" => string(now()),
        ))
    catch e
        @error "Failed to store embedding" report_id=report_id exception=(e, catch_backtrace())
        return Dict(
            "is_duplicate" => false,
            "report_id" => report_id,
            "matches" => [],
            "stored" => false,
            "error" => string(e),
        )
    end

    return Dict(
        "is_duplicate" => false,
        "report_id" => report_id,
        "matches" => [],
        "stored" => true,
    )
end
