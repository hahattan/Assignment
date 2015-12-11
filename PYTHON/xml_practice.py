import urllib
import xml.etree.ElementTree as ET

url = 'http://python-data.dr-chuck.net/comments_194438.xml'
sum = 0

uh = urllib.urlopen(url)
data = uh.read()
tree = ET.fromstring(data)

print data

comments = tree.find('comments').findall('comment')
print "Count:", len(comments)
#print comments
for comment in comments:
	count = comment.find('count').text
	sum += int(count)

print sum
