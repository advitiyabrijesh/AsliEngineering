local drivers = redis.call('georadius', KEYS[1], ARGV[1], ARGV[2], 5, 'km', 'WITHCOORD', 'WITHDIST')
local matchedDrivers = {}

for i, driver in ipairs(drivers) do
    local driverKey = 'location_info:' .. driver[1]
    local driverParams = redis.call('HGETALL', driverKey)

    for j, field in ipairs(driverParams) do
        if field == 'Rating' and tonumber(driverParams[j + 1]) >= tonumber(ARGV[3]) then
            table.insert(matchedDrivers, driver)
        end
    end
end

return matchedDrivers
