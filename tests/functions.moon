


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


-> if true then return

t = {
  my_method: => print @
  ->
    a = true
    return "hello"
  (foo, bar) =>
    print "world", foo, bar
  -> "world"
}
