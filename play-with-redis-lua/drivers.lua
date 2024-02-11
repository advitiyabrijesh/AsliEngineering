-- Sample: redis-cli eval "$(cat drivers.lua)" 2 driver_locations location_info 17.484710933268207 78.36023090530539 6

local drivers = redis.call('georadius', KEYS[1], ARGV[1], ARGV[2], 5, 'km', 'WITHCOORD', 'WITHDIST')
local matchedDrivers = {}

local function compare(a, b)
    return tonumber(a[2]) < tonumber(b[2])
end

for i, driver in ipairs(drivers) do
    local driverKey = 'location_info:' .. driver[1]
    local driverParams = redis.call('HGETALL', driverKey)

    for j, field in ipairs(driverParams) do
        if field == 'Rating' and tonumber(driverParams[j + 1]) >= tonumber(ARGV[3]) then
            table.insert(matchedDrivers, driver)
        end
    end
end

table.sort(matchedDrivers, compare)

return matchedDrivers
