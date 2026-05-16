using HTTP

struct QdrantClient
    host::String
    port::Int
    headers::Dict{String,String}
end

function QdrantClient(host::String, port::Int)
    QdrantClient(host, port, Dict("Content-Type" => "application/json"))
end

function _url(client::QdrantClient, path::String)
    "http://$(client.host):$(client.port)$path"
end

function create_collection(client::QdrantClient, name::String, vector_size::Int)
    url = _url(client, "/collections/$(name)")
    body = JSON3.write(Dict(
        "vectors" => Dict(
            "size" => vector_size,
            "distance" => "Cosine"
        )
    ))
    try
        resp = HTTP.put(url, client.headers, body)
        data = JSON3.parse(String(resp.body))
        if get(data, "result", false) === true
            @info "Created Qdrant collection" collection=name vector_size=vector_size
        else
            @warn "Qdrant create_collection response" status=resp.status result=data
        end
        return data
    catch e
        @error "Failed to create Qdrant collection" collection=name exception=(e, catch_backtrace())
        rethrow(e)
    end
end

function upsert_point(client::QdrantClient, collection::String, point_id::String, vector::Vector{Float32}, payload::Dict)
    url = _url(client, "/collections/$(collection)/points")
    body = JSON3.write(Dict(
        "points" => [
            Dict(
                "id" => point_id,
                "vector" => vector,
                "payload" => payload
            )
        ]
    ))
    try
        resp = HTTP.put(url, client.headers, body)
        data = JSON3.parse(String(resp.body))
        @info "Upserted point to Qdrant" collection=collection point_id=point_id
        return data
    catch e
        @error "Failed to upsert point" collection=collection point_id=point_id exception=(e, catch_backtrace())
        rethrow(e)
    end
end

function search_points(client::QdrantClient, collection::String, vector::Vector{Float32}; limit::Int=10)
    url = _url(client, "/collections/$(collection)/points/search")
    body = JSON3.write(Dict(
        "vector" => vector,
        "limit" => limit,
        "with_payload" => true,
        "with_vector" => false
    ))
    try
        resp = HTTP.post(url, client.headers, body)
        data = JSON3.parse(String(resp.body))
        results = get(data, "result", [])
        return [
            Dict(
                "id" => get(r, "id", ""),
                "score" => get(r, "score", 0.0),
                "payload" => get(r, "payload", Dict())
            )
            for r in results
        ]
    catch e
        @error "Failed to search Qdrant" collection=collection exception=(e, catch_backtrace())
        rethrow(e)
    end
end

function delete_collection(client::QdrantClient, name::String)
    url = _url(client, "/collections/$(name)")
    try
        resp = HTTP.delete(url, client.headers)
        data = JSON3.parse(String(resp.body))
        @info "Deleted Qdrant collection" collection=name
        return data
    catch e
        @error "Failed to delete Qdrant collection" collection=name exception=(e, catch_backtrace())
        rethrow(e)
    end
end
