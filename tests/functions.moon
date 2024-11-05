


hello = () ->
  print "hello world"

(a, b) ->
  if a < b
    return a - b
  else
    print "It's not right!"

  return a + b


f = -> print "cool"
g, b = ->, ->
