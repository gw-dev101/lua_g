--check that shadowing works correctly
local a = 10
if a > 5 then
  local a = 5
  print(a) -- should print 5
end
print(a) -- should print 10