from bs4 import BeautifulSoup
import urllib.request
# httpsL//redditlist.com
print("Getting page...")
page = ''
with open("temp/redditlist.html") as f:
    page = f.read()
print("Finding subreddits...")
soup = BeautifulSoup(page, 'html5lib')
subs = []

table = soup.findAll('div', attrs={'class':'listing'})

for row in table:
    print(row)
