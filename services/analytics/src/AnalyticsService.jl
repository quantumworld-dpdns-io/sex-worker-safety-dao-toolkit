module AnalyticsService

using Genie, Genie.Router, Genie.Requests, Genie.Responses
using JSON3
using DataFrames
using Dates
using Statistics
using Logging

include("config.jl")
include("qdrant.jl")
include("similarity.jl")
include("analytics.jl")
include("routes.jl")

function run()
    cfg = load_config()
    @info "Starting Analytics Service on port $(cfg.port)"
    Genie.config.run_as_server = true
    Genie.Server.up(cfg.port)
end

end
