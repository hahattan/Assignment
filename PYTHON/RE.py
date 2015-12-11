import re

sum = 0
hand = open('actual.txt')
for line in hand:
    line = line.rstrip()
    temp = re.findall('[0-9]+', line)
    for num in temp:
        sum += int(num)

print sum
