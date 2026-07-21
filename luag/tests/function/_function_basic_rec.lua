local function factorial(n)
 	if n == 0 then
 		return 1
 	else
 		return n * factorial(n - 1)
 	end
end

local result = factorial(5)
print(result) -- should print 120
local function fib(n,a,b)
 	if n == 0 then
 		return a
 	else
 		return fib(n - 1, b, a + b)
 	end
end

local result = fib(10, 0, 1)
print(result) -- should print 55
