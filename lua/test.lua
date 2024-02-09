file = io.open("settings.json", "rb")
contents = file.read(file)
io.close(file)

print(contents)

a = 3
b = 4

function myTest(value)
  print("comparing value " .. value .. " with 4")
  assert(value < 4)
end

myTest(a)
myTest(b)
