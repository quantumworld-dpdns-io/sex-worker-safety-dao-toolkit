function report_trends(reports::DataFrame)::Dict
    total = nrow(reports)
    if total == 0
        return Dict(
            "total_reports" => 0,
            "by_category" => Dict(),
            "by_region" => Dict(),
            "by_month" => Dict(),
            "trending_categories" => [],
        )
    end

    by_category = Dict{String,Int}()
    if hasproperty(reports, :category)
        for cat in unique(reports.category)
            by_category[cat] = count(==(cat), reports.category)
        end
    end

    by_region = Dict{String,Int}()
    if hasproperty(reports, :region)
        for r in unique(reports.region)
            by_region[r] = count(==(r), reports.region)
        end
    end

    by_month = Dict{String,Int}()
    if hasproperty(reports, :date)
        for d in reports.date
            month_key = Dates.format(Dates.Date(d), "Y-m")
            by_month[month_key] = get(by_month, month_key, 0) + 1
        end
    end

    sorted_cats = sort(collect(by_category), by=x -> x[2], rev=true)
    trending = [c[1] for c in sorted_cats[1:min(5, end)]]

    return Dict(
        "total_reports" => total,
        "by_category" => by_category,
        "by_region" => by_region,
        "by_month" => by_month,
        "trending_categories" => trending,
    )
end

function checkin_stats(checkins::DataFrame)::Dict
    total = nrow(checkins)
    if total == 0
        return Dict(
            "total_checkins" => 0,
            "completion_rate" => 0.0,
            "avg_duration_minutes" => 0.0,
            "by_hour" => Dict(),
            "by_day" => Dict(),
        )
    end

    completed = 0
    if hasproperty(checkins, :status)
        completed = count(==("completed"), checkins.status)
    end

    completion_rate = total > 0 ? completed / total * 100.0 : 0.0

    avg_duration = 0.0
    if hasproperty(checkins, :duration_minutes)
        avg_duration = mean(skipmissing(checkins.duration_minutes))
    end

    by_hour = Dict{String,Int}()
    by_day = Dict{String,Int}()
    if hasproperty(checkins, :timestamp)
        for ts in checkins.timestamp
            dt = Dates.DateTime(ts)
            hour_key = string(Dates.hour(dt))
            by_hour[hour_key] = get(by_hour, hour_key, 0) + 1
            day_key = Dates.dayname(dt)
            by_day[day_key] = get(by_day, day_key, 0) + 1
        end
    end

    return Dict(
        "total_checkins" => total,
        "completion_rate" => round(completion_rate, digits=1),
        "avg_duration_minutes" => round(avg_duration, digits=1),
        "by_hour" => by_hour,
        "by_day" => by_day,
    )
end

function dao_participation(votes::DataFrame)::Dict
    total = nrow(votes)
    if total == 0
        return Dict(
            "total_proposals" => 0,
            "total_votes" => 0,
            "avg_turnout" => 0.0,
            "vote_distribution" => Dict("for" => 0, "against" => 0, "abstain" => 0),
            "voter_retention" => 0.0,
        )
    end

    total_proposals = 0
    if hasproperty(votes, :proposal_id)
        total_proposals = length(unique(votes.proposal_id))
    end

    vote_dist = Dict("for" => 0, "against" => 0, "abstain" => 0)
    if hasproperty(votes, :vote)
        for v in votes.vote
            key = lowercase(string(v))
            if haskey(vote_dist, key)
                vote_dist[key] += 1
            end
        end
    end

    avg_turnout = 0.0
    if hasproperty(votes, :turnout)
        avg_turnout = mean(skipmissing(votes.turnout)) * 100.0
    end

    voter_retention = 0.0
    if hasproperty(votes, :voter_id) && total_proposals > 1
        voter_counts = combine(groupby(votes, :voter_id), nrow => :count)
        returned_voters = count(r -> r.count > 1, eachrow(voter_counts))
        voter_retention = total > 0 ? returned_voters / nrow(voter_counts) * 100.0 : 0.0
    end

    return Dict(
        "total_proposals" => total_proposals,
        "total_votes" => total,
        "avg_turnout" => round(avg_turnout, digits=1),
        "vote_distribution" => vote_dist,
        "voter_retention" => round(voter_retention, digits=1),
    )
end

function compute_confidence(report::Dict, similar_reports::Vector{Dict})::Float32
    if isempty(similar_reports)
        return 0.0f0
    end

    scores = [Float32(get(r, "score", 0.0)) for r in similar_reports]
    avg_score = mean(scores)
    count_factor = min(length(scores) / 5.0, 1.0)

    confidence = avg_score * 0.7f0 + Float32(count_factor) * 0.3f0
    return min(confidence, 1.0f0)
end
