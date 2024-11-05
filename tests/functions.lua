hello = function()
  print("hello world")
end
function(a, b)
  if a < b then
    return a - b
  else
    print("It's not right!")
  end
  return a + b
end
f = function()
  print("cool")
end
g, b = function()
end, function()
end
function()
  if true then
    return
  end
end
t = {my_method = function(self)
  print(self)
end, function()
  a = true
  return "hello"
end, function(self, foo, bar)
  print("world", foo, bar)
end, function()
  "world"
end}
