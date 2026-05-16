struct Config
    port::Int
    qdrant_host::String
    qdrant_port::Int
    database_url::String
end

function load_config()::Config
    Config(
        parse(Int, get(ENV, "ANALYTICS_PORT", "3003")),
        get(ENV, "QDRANT_HOST", "localhost"),
        parse(Int, get(ENV, "QDRANT_PORT", "6333")),
        get(ENV, "DATABASE_URL", ""),
    )
end
